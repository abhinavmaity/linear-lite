package repositories

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

func TestStrPtrEqual_CoversNilAndEqualCases(t *testing.T) {
	t.Parallel()

	if !strPtrEqual(nil, nil) {
		t.Fatalf("expected nil and nil to be equal")
	}
	a := "x"
	if strPtrEqual(&a, nil) {
		t.Fatalf("expected non-nil and nil to be different")
	}
	b := "x"
	if !strPtrEqual(&a, &b) {
		t.Fatalf("expected equal pointed values to be equal")
	}
	c := "y"
	if strPtrEqual(&a, &c) {
		t.Fatalf("expected different pointed values to be different")
	}
}

func TestIsUniqueViolationConstraint_DetectsPgxAndPqErrors(t *testing.T) {
	t.Parallel()

	pgxErr := &pgconn.PgError{Code: "23505", ConstraintName: "uq_labels_lower_name"}
	if !isUniqueViolationConstraint(pgxErr, "uq_labels_lower_name") {
		t.Fatalf("expected pgx unique violation to match constraint")
	}

	pqErr := &pq.Error{Code: "23505", Constraint: "uq_labels_lower_name"}
	if !isUniqueViolationConstraint(pqErr, "uq_labels_lower_name") {
		t.Fatalf("expected pq unique violation to match constraint")
	}

	notUnique := errors.New("other error")
	if isUniqueViolationConstraint(notUnique, "uq_labels_lower_name") {
		t.Fatalf("did not expect non-unique error to match")
	}
}

func TestIsUniqueViolationConstraint_ConstraintNameMustMatchWhenProvided(t *testing.T) {
	t.Parallel()

	pgxErr := &pgconn.PgError{Code: "23505", ConstraintName: "uq_projects_key"}
	if isUniqueViolationConstraint(pgxErr, "uq_labels_lower_name") {
		t.Fatalf("did not expect mismatched constraint to match")
	}

	pqErr := &pq.Error{Code: "23505", Constraint: "uq_projects_key"}
	if isUniqueViolationConstraint(pqErr, "uq_labels_lower_name") {
		t.Fatalf("did not expect mismatched constraint to match")
	}
}
