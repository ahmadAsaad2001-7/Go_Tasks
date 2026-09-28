markdown

# 🚀 Task Manager API - Go Learning Notes

Welcome to your Task Manager API! This document serves as a cheat sheet and reference guide for the concepts, architecture, and implementation details of this Go project.

---

## 📖 Project Overview
This is a RESTful API built in Go that allows users to manage projects and tasks. It includes authentication, role-based access control (RBAC), pagination, filtering, and soft deletes.

### 🛠️ Tech Stack
*   **Language:** Go (Golang)
*   **Web Framework:** Gin (`github.com/gin-gonic/gin`)
*   **Database:** PostgreSQL (with `github.com/lib/pq` driver)
*   **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`)
*   **Password Hashing:** Bcrypt (`golang.org/x/crypto/bcrypt`)
*   **Config Management:** `.env` files (`github.com/joho/godotenv`)

---

## 📁 Project Structure
```text
task-manager-api/
├── main.go                 # Entry point (starts the server)
├── .env                    # Environment variables (DB credentials, JWT secret)
├── db/
│   └── db.go               # Database connection and table creation
├── models/                 # Data structures and SQL queries
│   ├── user.go
│   ├── project.go
│   ├── task.go
│   └── tag.go
├── routes/                 # HTTP handlers (controllers)
│   ├── routes.go           # Registers all URLs to handlers
│   ├── auth_routes.go      # Signup and Login
│   ├── project_routes.go   # Project CRUD
│   └── task_routes.go      # Task CRUD
├── middlewares/            # Request interceptors
│   └── auth.go             # JWT validation and Role checks
├── utils/                  # Helper functions
│   ├── password.go         # Bcrypt hashing/comparison
│   └── jwt.go              # Token generation and verification
└── api/                    # VS Code REST Client test files (.http)

🗄️ Database Schema

The database uses 5 tables to support a Many-to-Many relationship between Tasks and Tags.

    users: id, email, password, role (user/admin), created_at

    projects: id, name, description, owner_id (FK), deleted_at, created_at

    tasks: id, title, description, status, project_id (FK), assignee_id (FK), deleted_at, created_at

    tags: id, name

    task_tags (Join Table): task_id (FK), tag_id (FK)

🧠 Core Concepts Implemented
1. Authentication (JWT)

    Signup: Hashes the password using Bcrypt and saves the user.

    Login: Compares passwords. If valid, generates a JWT containing userId, email, role, and an expiration time (exp).

    Authorization Middleware: Extracts the token from the Authorization header (trims "Bearer "), verifies it, and stores the user data in the Gin context using c.Set("user", user).

    Role-Based Access Control (RBAC): The AdminOnly middleware checks if user.Role == "admin". If not, it returns a 403 Forbidden.

2. Soft Deletes

Instead of actually deleting rows with DELETE FROM, we set deleted_at = NOW().

    CRITICAL RULE: Every SELECT query MUST include WHERE deleted_at IS NULL to avoid showing deleted records.

3. Pagination & Filtering

Implemented in getTasks:

    Reads query params: ?status=todo&page=1&limit=10

    Calculates offset = (page - 1) * limit

    SQL uses LIMIT $2 OFFSET $3

4. Many-to-Many Relationships

    Tasks and Tags are linked via the task_tags join table.

    We don't use a struct for TaskTag; instead, we use a simple function like LinkTagToTask(taskID, tagID).

5. Go "Enums"

Go doesn't have enums. We use a custom type and constants:
go

type TaskStatus string
const (
    StatusTodo       TaskStatus = "todo"
    StatusInProgress TaskStatus = "in_progress"
    StatusDone       TaskStatus = "done"
)

6. Gin Context & Type Assertions

When reading the user from the middleware:
go

userValue, exists := c.Get("user")
if !exists { ... }
user := userValue.(middlewares.User) // Type assertion

Tip: Use c.MustGet("user") if you are 100% sure it exists (like in protected routes).
🌐 API Endpoints Reference
Public Routes
Method	Endpoint	Description
POST	/auth/signup	Register a new user
POST	/auth/login	Login and receive a JWT
Protected Routes (Require Authorization: Bearer <token>)
Method	Endpoint	Description
POST	/projects	Create a new project
GET	/projects	Get all projects owned by the user
GET	/projects/:id	Get a specific project
PUT	/projects/:id	Update a project (Owner only)
DELETE	/projects/:id	Soft delete a project (Owner only)
POST	/projects/:id/tasks	Create a task inside a project
GET	/projects/:id/tasks	Get tasks (Supports ?status=todo&page=1&limit=10)
PUT	/tasks/:id	Update a task
DELETE	/tasks/:id	Soft delete a task
Admin Only Routes
Method	Endpoint	Description
DELETE	/admin/projects/:id	Soft delete ANY project (Bypasses owner check)
🚀 How to Run the Project
1. Database Setup

Ensure PostgreSQL is running. Create the database:
bash

sudo -i -u postgres
psql -c "CREATE DATABASE taskmanager;"
exit

2. Environment Variables (.env)

Create a .env file in the root folder:
env

DB_URL=postgres://postgres:yourpassword@localhost:5432/taskmanager?sslmode=disable
JWT_SECRET=super_secret_key_change_this_in_production

3. Install Dependencies
bash

go mod tidy

4. Run the Server
bash

go run main.go

The server will start on port 4001.
🧪 Testing the API (Using .http files)

Create a folder named api and add files like signup.http:

api/signup.http
http

POST http://localhost:4001/auth/signup
content-type: application/json

{
    "email": "admin@test.com",
    "password": "password123",
    "role": "admin"
}

api/create-project.http
http

POST http://localhost:4001/projects
content-type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
    "name": "My First Project",
    "description": "Learning Go"
}

⚠️ Common Gotchas & Lessons Learned

    Gin Route Conflicts: You cannot use different wildcard names for the same URL position (e.g., /projects/:id and /projects/:projectId/tasks will panic). Always standardize to :id.

    Missing } : A missing closing brace in Go will often cause confusing "undefined function" errors in other files. Always run go build to catch syntax errors early.

    Struct Tags: The backticks (`json:"name"`) are essential. They map Go struct fields to JSON keys and database columns.

    c.Param vs c.Query: c.Param("id") reads from the URL path (/projects/1). c.Query("page") reads from the query string (?page=1).

    Database Connection Pool: We use a global var DB *sql.DB to share the connection pool across the app. Never open a new connection inside a handler.
