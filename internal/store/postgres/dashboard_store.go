package postgres

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"project-manager-go/internal/domain"
)

type statusCount struct {
	Status domain.TaskStatus `db:"status"`
	Count  int               `db:"count"`
}

type memberWorkload struct {
	MemberID int `db:"member_id"`
	Count    int `db:"count"`
}

func (s *Store) DashboardSummary(ctx context.Context, user domain.AuthUser) (domain.DashboardSummary, error) {
	summary := domain.DashboardSummary{
		TaskCountPerStatus: map[domain.TaskStatus]int{
			domain.TaskStatusBacklog:    0,
			domain.TaskStatusTodo:       0,
			domain.TaskStatusInProgress: 0,
			domain.TaskStatusReview:     0,
			domain.TaskStatusDone:       0,
		},
		WorkloadPerMember:    map[string]int{},
		RecentlyUpdatedTasks: []domain.Task{},
		RecentTasks:          []domain.Task{},
		UpcomingDeadlines:    []domain.Task{},
		HighPriorityTasks:    []domain.Task{},
		LatestUpdates:        []domain.ActivityLog{},
	}

	totalActiveProjects, err := s.dashboardActiveProjectCount(ctx, user)
	if err != nil {
		return domain.DashboardSummary{}, err
	}
	summary.TotalActiveProjects = totalActiveProjects

	counts, err := s.dashboardTaskCountPerStatus(ctx, user)
	if err != nil {
		return domain.DashboardSummary{}, err
	}
	for _, item := range counts {
		summary.TaskCountPerStatus[item.Status] = item.Count
	}
	summary.InReview = summary.TaskCountPerStatus[domain.TaskStatusReview]
	summary.ActiveProgress = activeProgress(
		summary.TaskCountPerStatus[domain.TaskStatusTodo],
		summary.TaskCountPerStatus[domain.TaskStatusInProgress],
	)

	dueSoon, overdue, err := s.dashboardDeadlineCounts(ctx, user)
	if err != nil {
		return domain.DashboardSummary{}, err
	}
	summary.DueSoon = dueSoon
	summary.Overdue = overdue
	summary.OverdueTaskCount = overdue

	workload, err := s.dashboardWorkloadPerMember(ctx, user)
	if err != nil {
		return domain.DashboardSummary{}, err
	}
	for _, item := range workload {
		summary.WorkloadPerMember[strconv.Itoa(item.MemberID)] = item.Count
	}

	if summary.RecentlyUpdatedTasks, err = s.dashboardTasks(ctx, user, "ta.updated_at DESC, ta.id DESC", ""); err != nil {
		return domain.DashboardSummary{}, err
	}
	if summary.RecentTasks, err = s.dashboardTasks(ctx, user, "ta.created_at DESC, ta.id DESC", ""); err != nil {
		return domain.DashboardSummary{}, err
	}
	if summary.UpcomingDeadlines, err = s.dashboardTasks(ctx, user, "ta.due_date ASC, ta.id ASC", "ta.status <> 'done' AND ta.due_date >= CURRENT_DATE"); err != nil {
		return domain.DashboardSummary{}, err
	}
	if summary.HighPriorityTasks, err = s.dashboardTasks(ctx, user, "CASE ta.priority WHEN 'urgent' THEN 0 WHEN 'high' THEN 1 ELSE 2 END, ta.due_date ASC NULLS LAST, ta.id ASC", "ta.status <> 'done' AND ta.priority IN ('urgent', 'high')"); err != nil {
		return domain.DashboardSummary{}, err
	}

	return summary, nil
}

func (s *Store) dashboardActiveProjectCount(ctx context.Context, user domain.AuthUser) (int, error) {
	where := dashboardProjectWhere(user)
	where.add("p.status = $%d", domain.ProjectStatusActive)
	var total int
	err := s.db.GetContext(ctx, &total, "SELECT count(*) FROM projects p "+where.clause(), where.values...)
	return total, err
}

func (s *Store) dashboardTaskCountPerStatus(ctx context.Context, user domain.AuthUser) ([]statusCount, error) {
	where := dashboardTaskWhere(user)
	var rows []statusCount
	err := s.db.SelectContext(ctx, &rows, `
		SELECT ta.status, count(*) AS count
		FROM tasks ta
		`+where.clause()+`
		GROUP BY ta.status
	`, where.values...)
	return rows, err
}

func (s *Store) dashboardDeadlineCounts(ctx context.Context, user domain.AuthUser) (int, int, error) {
	where := dashboardTaskWhere(user)
	where.add("ta.status <> $%d", domain.TaskStatusDone)
	var counts struct {
		DueSoon int `db:"due_soon"`
		Overdue int `db:"overdue"`
	}
	err := s.db.GetContext(ctx, &counts, `
		SELECT
			count(*) FILTER (WHERE ta.due_date >= CURRENT_DATE AND ta.due_date <= CURRENT_DATE + INTERVAL '7 days') AS due_soon,
			count(*) FILTER (WHERE ta.due_date < CURRENT_DATE) AS overdue
		FROM tasks ta
		`+where.clause(), where.values...)
	return counts.DueSoon, counts.Overdue, err
}

func (s *Store) dashboardWorkloadPerMember(ctx context.Context, user domain.AuthUser) ([]memberWorkload, error) {
	where := dashboardTaskWhere(user)
	where.add("ta.status <> $%d", domain.TaskStatusDone)
	where.add("ta.assignee_id IS NOT NULL")
	var rows []memberWorkload
	err := s.db.SelectContext(ctx, &rows, `
		SELECT ta.assignee_id AS member_id, count(*) AS count
		FROM tasks ta
		`+where.clause()+`
		GROUP BY ta.assignee_id
	`, where.values...)
	return rows, err
}

func (s *Store) dashboardTasks(ctx context.Context, user domain.AuthUser, orderBy string, extraCondition string) ([]domain.Task, error) {
	where := dashboardTaskWhere(user)
	if extraCondition != "" {
		where.add(extraCondition)
	}
	query := fmt.Sprintf(`
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name AS project_name,
			ta.creator_id, COALESCE(ta.assignee_id, 0) AS assignee_id, u.name AS assignee_name, NULL::text AS assignee_avatar_url,
			ta.due_date::text AS due_date, ta.position, ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		%s
		ORDER BY %s
		LIMIT 5
	`, where.clause(), orderBy)
	tasks, err := selectAll[domain.Task](s, ctx, query, where.values...)
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		if tasks[i].ProjectName != nil {
			tasks[i].Project = &domain.TaskProject{
				ID:   tasks[i].ProjectID,
				Name: *tasks[i].ProjectName,
			}
		}
	}
	return tasks, nil
}

func dashboardProjectWhere(user domain.AuthUser) whereBuilder {
	where := whereBuilder{}
	if user.Role != domain.RoleAdmin {
		where.add(`(
			p.owner_id = $%d OR
			EXISTS (
				SELECT 1
				FROM project_teams access_pt
				JOIN team_members access_tm ON access_tm.team_id = access_pt.team_id
				WHERE access_pt.project_id = p.id AND access_tm.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM team_members primary_tm
				WHERE primary_tm.team_id = p.team_id AND primary_tm.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM tasks access_t
				WHERE access_t.project_id = p.id AND (
					access_t.creator_id = $%d OR access_t.assignee_id = $%d OR
					EXISTS (
						SELECT 1
						FROM task_assignments access_ta
						WHERE access_ta.task_id = access_t.id AND access_ta.user_id = $%d
					)
				)
			)
		)`, user.ID, user.ID, user.ID, user.ID, user.ID, user.ID)
	}
	return where
}

func dashboardTaskWhere(user domain.AuthUser) whereBuilder {
	where := whereBuilder{}
	if user.Role != domain.RoleAdmin {
		where.add(`(
			ta.creator_id = $%d OR
			ta.assignee_id = $%d OR
			EXISTS (
				SELECT 1
				FROM task_assignments access_ta
				WHERE access_ta.task_id = ta.id AND access_ta.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				WHERE access_p.id = ta.project_id AND access_p.owner_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				JOIN project_teams access_pt ON access_pt.project_id = access_p.id
				JOIN team_members additional_tm ON additional_tm.team_id = access_pt.team_id
				WHERE access_p.id = ta.project_id AND additional_tm.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				JOIN team_members primary_tm ON primary_tm.team_id = access_p.team_id
				WHERE access_p.id = ta.project_id AND primary_tm.user_id = $%d
			)
		)`, user.ID, user.ID, user.ID, user.ID, user.ID, user.ID)
	}
	return where
}

func activeProgress(todo int, doing int) domain.ActiveProgress {
	total := todo + doing
	if total == 0 {
		return domain.ActiveProgress{Doing: doing, Todo: todo, Total: total}
	}
	ratio := float64(doing) / float64(total)
	return domain.ActiveProgress{
		Doing:      doing,
		Todo:       todo,
		Total:      total,
		Ratio:      math.Round(ratio*10000) / 10000,
		Percentage: math.Round(ratio*10000) / 100,
	}
}
