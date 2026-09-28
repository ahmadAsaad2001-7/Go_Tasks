package models

import (
	"task-manager-api/db"
)

// Go "Enum" using custom types and constants
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	ProjectID   int        `json:"project_id"`
	AssigneeID  int        `json:"assignee_id"`
	Tags        []Tag      `json:"tags,omitempty"` // We will populate this later
}

// ---------- Queries ----------

func GetTasksByProject(projectID int, statusFilter string, limit, offset int) ([]Task, error) {
	// This query supports Filtering and Pagination!
	query := `
	SELECT id, title, description, status, project_id, assignee_id
	FROM tasks
	WHERE project_id = $1 AND deleted_at IS NULL
	AND ($2 = '' OR status = $2)
	ORDER BY id DESC
	LIMIT $3 OFFSET $4
	`
	rows, err := db.DB.Query(query, projectID, statusFilter, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssigneeID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTask(id string) (Task, error) {
	query := `SELECT id, title, description, status, project_id, assignee_id FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Task{}, err
	}
	defer stmt.Close()

	var t Task
	err = stmt.QueryRow(id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssigneeID)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func GetTasksByUser(userID int) ([]Task, error) {
	query := `SELECT id, title, description, status, project_id, assignee_id FROM tasks WHERE assignee_id = $1 AND deleted_at IS NULL`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssigneeID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// ---------- Commands ----------

func (t Task) AddTask() (Task, error) {
	query := `
	INSERT INTO tasks (title, description, status, project_id, assignee_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, description, status, project_id, assignee_id
	`
	err := db.DB.QueryRow(query, t.Title, t.Description, t.Status, t.ProjectID, t.AssigneeID).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssigneeID,
	)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (t Task) UpdateTask() (Task, error) {
	query := `
	UPDATE tasks
	SET title = $1, description = $2, status = $3, assignee_id = $4
	WHERE id = $5 AND deleted_at IS NULL
	RETURNING id, title, description, status, project_id, assignee_id
	`
	err := db.DB.QueryRow(query, t.Title, t.Description, t.Status, t.AssigneeID, t.ID).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.AssigneeID,
	)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (t Task) SoftDeleteTask() error {
	query := `UPDATE tasks SET deleted_at = NOW() WHERE id = $1`
	_, err := db.DB.Exec(query, t.ID)
	return err
}
