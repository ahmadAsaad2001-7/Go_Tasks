package models

import (
	"task-manager-api/db"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// AddTag creates a new tag or returns the existing one if it already exists
func (t Tag) AddTag() (Tag, error) {
	query := `
	INSERT INTO tags (name)
	VALUES ($1)
	ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
	RETURNING id, name
	`
	err := db.DB.QueryRow(query, t.Name).Scan(&t.ID, &t.Name)
	if err != nil {
		return Tag{}, err
	}
	return t, nil
}

// LinkTagToTask inserts a record into the join table
func LinkTagToTask(taskID, tagID int) error {
	query := `INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.DB.Exec(query, taskID, tagID)
	return err
}

// GetTagsByTask fetches all tags for a specific task
func GetTagsByTask(taskID int) ([]Tag, error) {
	query := `
	SELECT t.id, t.name
	FROM tags t
	JOIN task_tags tt ON t.id = tt.tag_id
	WHERE tt.task_id = $1
	`
	rows, err := db.DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}
