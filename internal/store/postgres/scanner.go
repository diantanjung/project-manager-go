package postgres

import (
	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"

	"github.com/jackc/pgx/v5"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner, includePassword bool) (domain.User, error) {
	var user domain.User
	var err error
	if includePassword {
		err = row.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.AvatarStorageKey,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	} else {
		err = row.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.AvatarStorageKey,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	}
	return user, notFound(err)
}

func scanTeam(row scanner) (domain.Team, error) {
	var team domain.Team
	err := row.Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)
	return team, notFound(err)
}

func scanProject(row scanner, joined bool) (domain.Project, error) {
	var project domain.Project
	var err error
	if joined {
		err = row.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.TeamID,
			&project.TeamName,
			&project.OwnerID,
			&project.OwnerName,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
	} else {
		err = row.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.TeamID,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
	}
	return project, notFound(err)
}

func scanProjectPage(
	rows pgx.Rows,
	filter service.PageFilter,
	total int,
	err error,
) (domain.Paginated[domain.Project], error) {
	if err != nil {
		return domain.Paginated[domain.Project]{}, err
	}
	defer rows.Close()

	out := []domain.Project{}
	for rows.Next() {
		project, err := scanProject(rows, true)
		if err != nil {
			return domain.Paginated[domain.Project]{}, err
		}
		out = append(out, project)
	}
	return paginated(out, filter, total), rows.Err()
}

func scanTask(row scanner) (domain.Task, error) {
	var task domain.Task
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ProjectID,
		&task.ProjectName,
		&task.CreatorID,
		&task.AssigneeID,
		&task.AssigneeName,
		&task.AssigneeAvatarURL,
		&task.DueDate,
		&task.Position,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, notFound(err)
}

func (s *Store) scanTaskPage(
	rows pgx.Rows,
	filter service.PageFilter,
	total int,
	err error,
) (domain.Paginated[domain.Task], error) {
	if err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	defer rows.Close()

	out := []domain.Task{}
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return domain.Paginated[domain.Task]{}, err
		}
		out = append(out, task)
	}
	return paginated(out, filter, total), rows.Err()
}
