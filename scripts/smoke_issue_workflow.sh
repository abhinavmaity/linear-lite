#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/docker-compose.yml"

if [[ -x "/usr/local/bin/docker" ]]; then
  DOCKER_BIN="${DOCKER_BIN:-/usr/local/bin/docker}"
elif [[ -x "/Applications/Docker.app/Contents/Resources/bin/docker" ]]; then
  DOCKER_BIN="${DOCKER_BIN:-/Applications/Docker.app/Contents/Resources/bin/docker}"
else
  DOCKER_BIN="${DOCKER_BIN:-docker}"
fi

compose() {
  "$DOCKER_BIN" compose -f "$COMPOSE_FILE" "$@"
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1"
    exit 1
  fi
}

assert_status() {
  local name="$1"
  local expected="$2"
  local actual="$3"
  if [[ "$expected" != "$actual" ]]; then
    echo "FAIL: $name expected $expected got $actual"
    exit 1
  fi
  echo "PASS: $name -> $actual"
}

json_extract() {
  local pattern="$1"
  local input="$2"
  printf '%s' "$input" | sed -n "s/$pattern/\\1/p"
}

http_json() {
  local method="$1"
  local url="$2"
  local out_file="$3"
  local auth="${4:-}"
  local data="${5:-}"

  local args=(-sS -o "$out_file" -w "%{http_code}" -X "$method" "$url")
  if [[ -n "$auth" ]]; then
    args+=(-H "Authorization: Bearer $auth")
  fi
  if [[ -n "$data" ]]; then
    args+=(-H "Content-Type: application/json" -d "$data")
  fi

  curl "${args[@]}"
}

seed_sql_returning_id() {
  local sql="$1"
  compose exec -T postgres psql -U postgres -d linear_lite -q -t -A -c "$sql" \
    | tr -d '\r' | head -n1 | xargs
}

require_cmd curl
require_cmd sed
require_cmd tr
require_cmd head
require_cmd xargs

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "==> Starting backend stack"
compose up -d --build

echo "==> Running migrations"
compose exec -T backend migrate

EMAIL="smoke.$(date +%s)@example.com"
REGISTER_PAYLOAD="{\"name\":\"Smoke Tester\",\"email\":\"$EMAIL\",\"password\":\"Password123\"}"
REGISTER_STATUS="$(http_json POST "http://localhost:8080/api/v1/auth/register" "$TMP_DIR/register.json" "" "$REGISTER_PAYLOAD")"
REGISTER_BODY="$(cat "$TMP_DIR/register.json")"
assert_status "auth register" "201" "$REGISTER_STATUS"

TOKEN="$(json_extract '.*"token":"\([^"]*\)".*' "$REGISTER_BODY")"
USER_ID="$(json_extract '.*"user":{"id":"\([^"]*\)".*' "$REGISTER_BODY")"
if [[ -z "$TOKEN" || -z "$USER_ID" ]]; then
  echo "FAIL: unable to parse token or user id from register response"
  echo "$REGISTER_BODY"
  exit 1
fi

PROJECT_KEY="SMK$(date +%H%M%S)"
PROJECT_ID="$(seed_sql_returning_id "INSERT INTO projects (name, description, key, created_by) VALUES ('Smoke Project', 'Smoke project', '$PROJECT_KEY', '$USER_ID') RETURNING id;")"
SPRINT_ID="$(seed_sql_returning_id "INSERT INTO sprints (name, description, project_id, start_date, end_date, status) VALUES ('Smoke Sprint', 'Sprint for smoke tests', '$PROJECT_ID', CURRENT_DATE, CURRENT_DATE + 14, 'planned') RETURNING id;")"
LABEL_ID="$(seed_sql_returning_id "INSERT INTO labels (name, color, description) VALUES ('smoke-label-$(date +%H%M%S)', '#2563EB', 'Smoke label') RETURNING id;")"

echo "==> Testing issue create"
CREATE_PAYLOAD="{\"title\":\"Smoke issue\",\"description\":\"Initial description\",\"project_id\":\"$PROJECT_ID\",\"sprint_id\":\"$SPRINT_ID\",\"assignee_id\":\"$USER_ID\",\"label_ids\":[\"$LABEL_ID\"],\"status\":\"backlog\",\"priority\":\"medium\"}"
CREATE_STATUS="$(http_json POST "http://localhost:8080/api/v1/issues" "$TMP_DIR/create.json" "$TOKEN" "$CREATE_PAYLOAD")"
CREATE_BODY="$(cat "$TMP_DIR/create.json")"
assert_status "issues create" "201" "$CREATE_STATUS"

ISSUE_ID="$(json_extract '.*"data":{"id":"\([^"]*\)".*' "$CREATE_BODY")"
if [[ -z "$ISSUE_ID" ]]; then
  echo "FAIL: unable to parse issue id from create response"
  echo "$CREATE_BODY"
  exit 1
fi

echo "==> Testing issue list/detail/update/archive/restore"
LIST_STATUS="$(http_json GET "http://localhost:8080/api/v1/issues?page=1&limit=20" "$TMP_DIR/list.json" "$TOKEN")"
assert_status "issues list" "200" "$LIST_STATUS"
if ! grep -q "$ISSUE_ID" "$TMP_DIR/list.json"; then
  echo "FAIL: issues list response missing created issue"
  cat "$TMP_DIR/list.json"
  exit 1
fi

DETAIL_STATUS="$(http_json GET "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/detail.json" "$TOKEN")"
assert_status "issue detail" "200" "$DETAIL_STATUS"

UPDATE_PAYLOAD='{"status":"in_progress","priority":"high","title":"Smoke issue updated"}'
UPDATE_STATUS="$(http_json PUT "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/update.json" "$TOKEN" "$UPDATE_PAYLOAD")"
assert_status "issue update" "200" "$UPDATE_STATUS"
if ! grep -q '"title":"Smoke issue updated"' "$TMP_DIR/update.json"; then
  echo "FAIL: update response missing updated title"
  cat "$TMP_DIR/update.json"
  exit 1
fi

PUT_ARCHIVED_TRUE_STATUS="$(http_json PUT "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/put_archived_true.json" "$TOKEN" '{"archived":true}')"
assert_status "put archived true rejected" "400" "$PUT_ARCHIVED_TRUE_STATUS"

ARCHIVE_STATUS="$(http_json DELETE "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/archive.txt" "$TOKEN")"
assert_status "issue archive" "204" "$ARCHIVE_STATUS"

DETAIL_AFTER_ARCHIVE_STATUS="$(http_json GET "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/detail_after_archive.json" "$TOKEN")"
assert_status "detail after archive default" "404" "$DETAIL_AFTER_ARCHIVE_STATUS"

DETAIL_WITH_ARCHIVE_STATUS="$(http_json GET "http://localhost:8080/api/v1/issues/$ISSUE_ID?include_archived=true" "$TMP_DIR/detail_include_archived.json" "$TOKEN")"
assert_status "detail include archived" "200" "$DETAIL_WITH_ARCHIVE_STATUS"

RESTORE_STATUS="$(http_json PUT "http://localhost:8080/api/v1/issues/$ISSUE_ID" "$TMP_DIR/restore.json" "$TOKEN" '{"archived":false}')"
assert_status "issue restore" "200" "$RESTORE_STATUS"

echo
echo "Smoke workflow passed."
echo "issue_id=$ISSUE_ID"
