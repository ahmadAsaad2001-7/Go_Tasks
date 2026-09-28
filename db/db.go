package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB // Fixed: Uppercase DB to match usage below

func InitDB() { // Fixed: Added space before brace
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	var errOpen error
	DB, errOpen = sql.Open("postgres", connStr) // Fixed: Use = instead of := because DB is global
	if errOpen != nil {
		panic(fmt.Sprintf("Error connecting to PostgreSQL: %v", errOpen))
	}

	DB.SetMaxOpenConns(10) // Fixed: Indentation
	DB.SetMaxIdleConns(3)

	if err := DB.Ping(); err != nil {
		panic(fmt.Sprintf("Error pinging PostgreSQL: %v", err))
	}
	fmt.Println("Connected to PostgreSQL!")

	createTable() // Fixed: Call the function
}

func createTable() { // Fixed: Added func keyword
	CreateUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`
		_, err := DB.Exec(CreateUsersTable) // Fixed: First time using := is fine
		if err != nil {
			panic(fmt.Sprintf("Error creating users table: %v", err))
		}

		createProjectsTable := `
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			owner_id INTEGER NOT NULL,
			deleted_at TIMESTAMP DEFAULT NULL, -- Added for soft deletes
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
			)
			`
			_, err = DB.Exec(createProjectsTable) // Fixed: Use = instead of :=
			if err != nil {
				panic(fmt.Sprintf("Error creating projects table: %v", err))
			}

			createTasksTable := `
			CREATE TABLE IF NOT EXISTS tasks (
				id SERIAL PRIMARY KEY,
				title TEXT NOT NULL,
				status TEXT NOT NULL,
				description TEXT NOT NULL,
				project_id INTEGER NOT NULL,
				assignee_id INTEGER NOT NULL, -- Fixed: Typo (removed extra 's')
				deleted_at TIMESTAMP DEFAULT NULL, -- Added for soft deletes
				FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE CASCADE, -- Fixed: Added comma
				FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
				)
				`
				_, err = DB.Exec(createTasksTable) // Fixed: Use = and pass the correct table
				if err != nil {
					panic(fmt.Sprintf("Error creating tasks table: %v", err))
				}

				// NEW: Tags Table
				createTagsTable := `
				CREATE TABLE IF NOT EXISTS tags (
					id SERIAL PRIMARY KEY,
					name TEXT NOT NULL UNIQUE
					)
					`
					_, err = DB.Exec(createTagsTable)
					if err != nil {
						panic(fmt.Sprintf("Error creating tags table: %v", err))
					}

					// NEW: Task Tags Join Table (Many-to-Many)
					createTaskTagsTable := `
					CREATE TABLE IF NOT EXISTS task_tags (
						task_id INTEGER NOT NULL,
						tag_id INTEGER NOT NULL,
						PRIMARY KEY (task_id, tag_id),
						FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
						FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
						)
						`
						_, err = DB.Exec(createTaskTagsTable)
						if err != nil {
							panic(fmt.Sprintf("Error creating task_tags table: %v", err))
						}

						fmt.Println("Tables created successfully!")
} // Fixed: Missing closing brace for createTable()
