# Task API

A RESTful task management API built with Go 1.18 and Gin framework.

## Live Demo

This project is automatically deployed to Google Cloud Run using GitHub Actions:
- **API Backend**: https://task-api-fdrvii3bmq-de.a.run.app
- **Frontend Interface**: https://etrex.github.io/task-api/

The backend API is continuously deployed to GCP Cloud Run whenever changes are pushed to the main branch. The frontend is a single-page application hosted on GitHub Pages that provides a web interface for testing and interacting with the API.

## Features

- Create, read, update, and delete tasks
- In-memory storage
- RESTful API design
- Comprehensive unit tests
- Docker support

## API Endpoints

- `GET /tasks` - List all tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update a task
- `DELETE /tasks/{id}` - Delete a task

## Task Model

```json
{
  "id": "string (UUID)",
  "name": "string (required)",
  "status": "integer (0 or 1, required)"
}
```

- `status: 0` - Incomplete task
- `status: 1` - Completed task

## Running with Docker

### Build the image
```bash
docker build -t task-api .
```

### Run the container
```bash
docker run -d -p 8080:8080 --name task-api task-api
```

### Stop the container
```bash
docker stop task-api && docker rm task-api
```

## Running locally

### Prerequisites
- Go 1.18 or later

### Install dependencies
```bash
go mod download
```

### Run the application
```bash
go run .
```

### Run tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test ./... -cover
```

## Example Usage

### Create a task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name":"Learn Go","status":0}'
```

### List all tasks
```bash
curl http://localhost:8080/tasks
```

### Update a task
```bash
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Learn Go","status":1}'
```

### Delete a task
```bash
curl -X DELETE http://localhost:8080/tasks/{id}
```