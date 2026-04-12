#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/docker-compose.backend.yml"

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

assert_eq() {
  local name="$1"
  local expected="$2"
  local actual="$3"
  if [[ "$expected" != "$actual" ]]; then
    echo "FAIL: $name expected '$expected' got '$actual'"
    exit 1
  fi
  echo "PASS: $name -> $actual"
}

assert_contains() {
  local name="$1"
  local needle="$2"
  local file="$3"
  if ! grep -q "$needle" "$file"; then
    echo "FAIL: $name missing '$needle'"
    cat "$file"
    exit 1
  fi
  echo "PASS: $name contains '$needle'"
}

assert_not_contains() {
  local name="$1"
  local needle="$2"
  local file="$3"
  if grep -q "$needle" "$file"; then
    echo "FAIL: $name unexpectedly contains '$needle'"
    cat "$file"
    exit 1
  fi
  echo "PASS: $name does not contain '$needle'"
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

seed_sql_exec() {
  local sql="$1"
  compose exec -T postgres psql -U postgres -d linear_lite -q -c "$sql" >/dev/null
}

redis_flush() {
  compose exec -T redis redis-cli FLUSHALL >/dev/null
}

redis_count() {
  local pattern="$1"
  compose exec -T redis redis-cli --raw KEYS "$pattern" \
    | sed '/^$/d' \
    | wc -l \
    | tr -d ' '
}

redis_key_exists() {
  local key="$1"
  local count
  count="$(compose exec -T redis redis-cli EXISTS "$key" | tr -d '\r\n')"
  if [[ "$count" == "1" ]]; then
    echo "1"
  else
    echo "0"
  fi
}

require_cmd curl
require_cmd sed
require_cmd tr
require_cmd head
require_cmd xargs
require_cmd grep
require_cmd wc

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "==> Starting backend stack"
compose up -d --build

echo "==> Running migrations"
compose exec -T backend migrate

EMAIL="cache.$(date +%s)@example.com"
REGISTER_PAYLOAD="{\"name\":\"Cache Tester\",\"email\":\"$EMAIL\",\"password\":\"Password123\"}"
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

PROJECT_KEY="CCH$(date +%H%M%S)"
PROJECT_ID="$(seed_sql_returning_id "INSERT INTO projects (name, description, key, created_by) VALUES ('Cache Project', 'Cache project', '$PROJECT_KEY', '$USER_ID') RETURNING id;")"
SPRINT_ID="$(seed_sql_returning_id "INSERT INTO sprints (name, description, project_id, start_date, end_date, status) VALUES ('Cache Sprint', 'Sprint for cache tests', '$PROJECT_ID', CURRENT_DATE, CURRENT_DATE + 14, 'planned') RETURNING id;")"
LABEL_ID="$(seed_sql_returning_id "INSERT INTO labels (name, color, description) VALUES ('cache-label-$(date +%H%M%S)', '#2563EB', 'Cache label') RETURNING id;")"

echo "==> Users cache: miss/hit + register invalidation"
redis_flush

USERS_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/users?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/users_1.json" "$TOKEN")"
assert_status "users first fetch" "200" "$USERS_1_STATUS"
assert_eq "users cache key count after miss" "1" "$(redis_count 'users:list:*')"

DIRECT_USER_EMAIL="cache.direct.$(date +%s)@example.com"
seed_sql_exec "INSERT INTO users (email, password_hash, name) VALUES ('$DIRECT_USER_EMAIL', 'x', 'Cache Direct User');"

USERS_2_STATUS="$(http_json GET "http://localhost:8080/api/v1/users?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/users_2.json" "$TOKEN")"
assert_status "users second fetch (cache hit)" "200" "$USERS_2_STATUS"
assert_not_contains "users cached response" "$DIRECT_USER_EMAIL" "$TMP_DIR/users_2.json"

SECOND_REGISTER_EMAIL="cache.invalidate.$(date +%s)@example.com"
SECOND_REGISTER_PAYLOAD="{\"name\":\"Cache Invalidate\",\"email\":\"$SECOND_REGISTER_EMAIL\",\"password\":\"Password123\"}"
SECOND_REGISTER_STATUS="$(http_json POST "http://localhost:8080/api/v1/auth/register" "$TMP_DIR/register_2.json" "" "$SECOND_REGISTER_PAYLOAD")"
assert_status "auth register invalidates users cache" "201" "$SECOND_REGISTER_STATUS"
assert_eq "users cache key count after register invalidation" "0" "$(redis_count 'users:*')"

USERS_3_STATUS="$(http_json GET "http://localhost:8080/api/v1/users?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/users_3.json" "$TOKEN")"
assert_status "users fetch after invalidation" "200" "$USERS_3_STATUS"
assert_contains "users refreshed response includes direct insert" "$DIRECT_USER_EMAIL" "$TMP_DIR/users_3.json"

echo "==> Projects cache: miss/hit + mutation invalidation"
redis_flush

PROJECTS_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/projects?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/projects_1.json" "$TOKEN")"
assert_status "projects first fetch" "200" "$PROJECTS_1_STATUS"
assert_eq "projects list key count after miss" "1" "$(redis_count 'projects:list:*')"

PROJECT_DETAIL_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/projects/$PROJECT_ID" "$TMP_DIR/project_detail_1.json" "$TOKEN")"
assert_status "project detail first fetch" "200" "$PROJECT_DETAIL_1_STATUS"
PROJECT_DETAIL_KEY="projects:detail:$PROJECT_ID"
assert_eq "project detail key exists" "1" "$(redis_key_exists "$PROJECT_DETAIL_KEY")"

HIDDEN_PROJECT_KEY="HID$(date +%H%M%S)"
seed_sql_exec "INSERT INTO projects (name, description, key, created_by) VALUES ('Hidden Project', 'Should be hidden by cache hit', '$HIDDEN_PROJECT_KEY', '$USER_ID');"

PROJECTS_2_STATUS="$(http_json GET "http://localhost:8080/api/v1/projects?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/projects_2.json" "$TOKEN")"
assert_status "projects second fetch (cache hit)" "200" "$PROJECTS_2_STATUS"
assert_not_contains "projects cached response" "$HIDDEN_PROJECT_KEY" "$TMP_DIR/projects_2.json"

CREATE_PROJECT_PAYLOAD="{\"name\":\"Created Project\",\"description\":\"from api\",\"key\":\"API$(date +%H%M%S)\"}"
CREATE_PROJECT_STATUS="$(http_json POST "http://localhost:8080/api/v1/projects" "$TMP_DIR/project_create.json" "$TOKEN" "$CREATE_PROJECT_PAYLOAD")"
assert_status "project create" "201" "$CREATE_PROJECT_STATUS"
assert_eq "projects list keys invalidated after project create" "0" "$(redis_count 'projects:list:*')"
assert_eq "old project detail key invalidated after project create" "0" "$(redis_key_exists "$PROJECT_DETAIL_KEY")"

echo "==> Sprints cache: miss/hit + mutation invalidation"
redis_flush

SPRINTS_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/sprints?page=1&limit=50&sort_by=start_date&sort_order=desc" "$TMP_DIR/sprints_1.json" "$TOKEN")"
assert_status "sprints first fetch" "200" "$SPRINTS_1_STATUS"
assert_eq "sprints list key count after miss" "1" "$(redis_count 'sprints:list:*')"

SPRINT_DETAIL_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/sprints/$SPRINT_ID" "$TMP_DIR/sprint_detail_1.json" "$TOKEN")"
assert_status "sprint detail first fetch" "200" "$SPRINT_DETAIL_1_STATUS"
SPRINT_DETAIL_KEY="sprints:detail:$SPRINT_ID"
assert_eq "sprint detail key exists" "1" "$(redis_key_exists "$SPRINT_DETAIL_KEY")"

HIDDEN_SPRINT_NAME="Hidden Sprint $(date +%s)"
seed_sql_exec "INSERT INTO sprints (name, description, project_id, start_date, end_date, status) VALUES ('$HIDDEN_SPRINT_NAME', 'hidden by cache', '$PROJECT_ID', CURRENT_DATE + 20, CURRENT_DATE + 30, 'planned');"

SPRINTS_2_STATUS="$(http_json GET "http://localhost:8080/api/v1/sprints?page=1&limit=50&sort_by=start_date&sort_order=desc" "$TMP_DIR/sprints_2.json" "$TOKEN")"
assert_status "sprints second fetch (cache hit)" "200" "$SPRINTS_2_STATUS"
assert_not_contains "sprints cached response" "$HIDDEN_SPRINT_NAME" "$TMP_DIR/sprints_2.json"

DASHBOARD_WARM_STATUS="$(http_json GET "http://localhost:8080/api/v1/dashboard/stats" "$TMP_DIR/dashboard_warm_sprint.json" "$TOKEN")"
assert_status "dashboard warm before sprint update" "200" "$DASHBOARD_WARM_STATUS"
assert_eq "dashboard key exists before sprint update" "1" "$(redis_key_exists "dashboard:stats:$USER_ID")"

SPRINT_UPDATE_PAYLOAD="{\"name\":\"Cache Sprint Updated\"}"
SPRINT_UPDATE_STATUS="$(http_json PUT "http://localhost:8080/api/v1/sprints/$SPRINT_ID" "$TMP_DIR/sprint_update.json" "$TOKEN" "$SPRINT_UPDATE_PAYLOAD")"
assert_status "sprint update" "200" "$SPRINT_UPDATE_STATUS"
assert_eq "sprints list keys invalidated after sprint update" "0" "$(redis_count 'sprints:list:*')"
assert_eq "projects keys invalidated after sprint update" "0" "$(redis_count 'projects:*')"
assert_eq "dashboard keys invalidated after sprint update" "0" "$(redis_count 'dashboard:*')"

echo "==> Labels cache: miss/hit + mutation invalidation"
redis_flush

LABELS_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/labels?page=1&limit=100&sort_by=name&sort_order=asc" "$TMP_DIR/labels_1.json" "$TOKEN")"
assert_status "labels first fetch" "200" "$LABELS_1_STATUS"
assert_eq "labels list key count after miss" "1" "$(redis_count 'labels:list:*')"

HIDDEN_LABEL_NAME="hidden-label-$(date +%s)"
seed_sql_exec "INSERT INTO labels (name, color, description) VALUES ('$HIDDEN_LABEL_NAME', '#22C55E', 'hidden by cache');"

LABELS_2_STATUS="$(http_json GET "http://localhost:8080/api/v1/labels?page=1&limit=100&sort_by=name&sort_order=asc" "$TMP_DIR/labels_2.json" "$TOKEN")"
assert_status "labels second fetch (cache hit)" "200" "$LABELS_2_STATUS"
assert_not_contains "labels cached response" "$HIDDEN_LABEL_NAME" "$TMP_DIR/labels_2.json"

# Keep updated label names unique per run so reruns on reused DB state
# do not fail with expected 409 conflicts unrelated to cache behavior.
LABEL_UPDATE_NAME="cache-label-updated-$(date +%s)"
LABEL_UPDATE_PAYLOAD="{\"name\":\"$LABEL_UPDATE_NAME\",\"color\":\"#EF4444\"}"
LABEL_UPDATE_STATUS="$(http_json PUT "http://localhost:8080/api/v1/labels/$LABEL_ID" "$TMP_DIR/label_update.json" "$TOKEN" "$LABEL_UPDATE_PAYLOAD")"
assert_status "label update" "200" "$LABEL_UPDATE_STATUS"
assert_eq "labels keys invalidated after label update" "0" "$(redis_count 'labels:*')"

echo "==> Dashboard + issue write invalidation"
redis_flush

DASHBOARD_1_STATUS="$(http_json GET "http://localhost:8080/api/v1/dashboard/stats" "$TMP_DIR/dashboard_1.json" "$TOKEN")"
assert_status "dashboard first fetch" "200" "$DASHBOARD_1_STATUS"
assert_eq "dashboard key exists after first fetch" "1" "$(redis_key_exists "dashboard:stats:$USER_ID")"

TOTAL_ISSUES_1="$(json_extract '.*"total_issues":\([0-9][0-9]*\).*' "$(cat "$TMP_DIR/dashboard_1.json")")"
if [[ -z "$TOTAL_ISSUES_1" ]]; then
  echo "FAIL: unable to parse dashboard total_issues"
  cat "$TMP_DIR/dashboard_1.json"
  exit 1
fi

seed_sql_exec "INSERT INTO issues (identifier, title, status, priority, project_id, sprint_id, created_by) VALUES ('${PROJECT_KEY}-9999', 'Dashboard hidden by cache', 'backlog', 'medium', '$PROJECT_ID', '$SPRINT_ID', '$USER_ID');"
DASHBOARD_2_STATUS="$(http_json GET "http://localhost:8080/api/v1/dashboard/stats" "$TMP_DIR/dashboard_2.json" "$TOKEN")"
assert_status "dashboard second fetch (cache hit)" "200" "$DASHBOARD_2_STATUS"
TOTAL_ISSUES_2="$(json_extract '.*"total_issues":\([0-9][0-9]*\).*' "$(cat "$TMP_DIR/dashboard_2.json")")"
assert_eq "dashboard cached total_issues unchanged" "$TOTAL_ISSUES_1" "$TOTAL_ISSUES_2"

PROJECTS_WARM_STATUS="$(http_json GET "http://localhost:8080/api/v1/projects?page=1&limit=50&sort_by=name&sort_order=asc" "$TMP_DIR/projects_warm_issue.json" "$TOKEN")"
assert_status "projects warm before issue create" "200" "$PROJECTS_WARM_STATUS"
SPRINTS_WARM_STATUS="$(http_json GET "http://localhost:8080/api/v1/sprints?page=1&limit=50&sort_by=start_date&sort_order=desc" "$TMP_DIR/sprints_warm_issue.json" "$TOKEN")"
assert_status "sprints warm before issue create" "200" "$SPRINTS_WARM_STATUS"

ISSUE_CREATE_PAYLOAD="{\"title\":\"Cache invalidation issue\",\"description\":\"test\",\"project_id\":\"$PROJECT_ID\",\"sprint_id\":\"$SPRINT_ID\",\"assignee_id\":\"$USER_ID\",\"label_ids\":[\"$LABEL_ID\"],\"status\":\"todo\",\"priority\":\"medium\"}"
ISSUE_CREATE_STATUS="$(http_json POST "http://localhost:8080/api/v1/issues" "$TMP_DIR/issue_create_cache.json" "$TOKEN" "$ISSUE_CREATE_PAYLOAD")"
assert_status "issue create invalidates caches" "201" "$ISSUE_CREATE_STATUS"
assert_eq "dashboard keys invalidated after issue create" "0" "$(redis_count 'dashboard:*')"
assert_eq "projects keys invalidated after issue create" "0" "$(redis_count 'projects:*')"
assert_eq "sprints keys invalidated after issue create" "0" "$(redis_count 'sprints:*')"

echo
echo "Milestone 4 cache smoke passed."
