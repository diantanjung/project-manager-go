package postgres

import (
	"strings"
	"testing"

	"project-manager-go/internal/domain"
)

func TestProjectWhereUsesExistsForAccessChecks(t *testing.T) {
	where := projectWhere(domain.AuthUser{
		ID:   42,
		Role: domain.RoleTeamMember,
	})

	clause := where.clause()
	if strings.Contains(clause, "primary_tm.user_id = $2 OR additional_tm.user_id = $3") {
		t.Fatalf("project access uses row-multiplying joins: %s", clause)
	}
	if got := strings.Count(clause, "EXISTS"); got != 2 {
		t.Fatalf("EXISTS count = %d, want 2 in clause: %s", got, clause)
	}
	if len(where.values) != 3 {
		t.Fatalf("values length = %d, want 3", len(where.values))
	}
	for i, value := range where.values {
		if value != 42 {
			t.Fatalf("values[%d] = %v, want 42", i, value)
		}
	}
}

func TestProjectWhereProductOwnerDoesNotFilterAccess(t *testing.T) {
	where := projectWhere(domain.AuthUser{
		ID:   42,
		Role: domain.RoleProductOwner,
	})

	if clause := where.clause(); clause != "" {
		t.Fatalf("clause = %q, want empty", clause)
	}
	if len(where.values) != 0 {
		t.Fatalf("values length = %d, want 0", len(where.values))
	}
}

func TestTaskWhereUsesExistsForAccessChecks(t *testing.T) {
	where := taskWhere(domain.AuthUser{
		ID:   42,
		Role: domain.RoleTeamMember,
	})

	clause := where.clause()
	if strings.Contains(clause, "LEFT JOIN") {
		t.Fatalf("task access should not require row-multiplying joins: %s", clause)
	}
	if got := strings.Count(clause, "EXISTS"); got != 4 {
		t.Fatalf("EXISTS count = %d, want 4 in clause: %s", got, clause)
	}
	if len(where.values) != 6 {
		t.Fatalf("values length = %d, want 6", len(where.values))
	}
	for i, value := range where.values {
		if value != 42 {
			t.Fatalf("values[%d] = %v, want 42", i, value)
		}
	}
}

func TestTaskWhereProductOwnerDoesNotFilterAccess(t *testing.T) {
	where := taskWhere(domain.AuthUser{
		ID:   42,
		Role: domain.RoleProductOwner,
	})

	if clause := where.clause(); clause != "" {
		t.Fatalf("clause = %q, want empty", clause)
	}
	if len(where.values) != 0 {
		t.Fatalf("values length = %d, want 0", len(where.values))
	}
}
