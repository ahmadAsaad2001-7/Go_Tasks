package models

import (
	"task-manager-api/db"
)

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
}

// ---------- Queries ----------

func GetProjectsByUser(userID int) ([]Project, error) {
	query := `SELECT id, name, description, owner_id FROM projects WHERE owner_id = $1 AND deleted_at IS NULL`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func GetProject(id string) (Project, error) {
	query := `SELECT id, name, description, owner_id FROM projects WHERE id = $1 AND deleted_at IS NULL`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	var p Project
	err = stmt.QueryRow(id).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
	if err != nil {
		return Project{}, err
	}
	return p, nil
}

func GetProjectByName(name string) (Project, error) {
	query := `SELECT id, name, description, owner_id FROM projects WHERE name = $1 AND deleted_at IS NULL`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	var p Project
	err = stmt.QueryRow(name).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
	if err != nil {
		return Project{}, err
	}
	return p, nil
}

// ---------- Commands ----------

func (p Project) AddProject() (Project, error) {
	query := `
	INSERT INTO projects (name, description, owner_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, owner_id
	`
	err := db.DB.QueryRow(query, p.Name, p.Description, p.OwnerID).Scan(
		&p.ID, &p.Name, &p.Description, &p.OwnerID,
	)
	if err != nil {
		return Project{}, err
	}
	return p, nil
}

func (p Project) UpdateProject() (Project, error) {
	query := `
	UPDATE projects
	SET name = $1, description = $2
	WHERE id = $3 AND owner_id = $4 AND deleted_at IS NULL
	RETURNING id, name, description, owner_id
	`
	err := db.DB.QueryRow(query, p.Name, p.Description, p.ID, p.OwnerID).Scan(
		&p.ID, &p.Name, &p.Description, &p.OwnerID,
	)
	if err != nil {
		return Project{}, err
	}
	return p, nil
}

func (p Project) SoftDeleteProject() error {
	// We don't actually DELETE, we just set the timestamp
	query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1 AND owner_id = $2`
	_, err := db.DB.Exec(query, p.ID, p.OwnerID)
	return err
}
