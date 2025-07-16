# Task API

A RESTful task management API built with Go 1.18 and Gin framework.

## Live Demo

This project is automatically deployed to Google Cloud Run using GitHub Actions:
- **API Backend**: https://task-api.etrex.tw
- **Frontend Interface**: https://etrex.github.io/task-api/

The backend API is continuously deployed to GCP Cloud Run whenever changes are pushed to the main branch. The frontend is a single-page application hosted on GitHub Pages that provides a web interface for testing and interacting with the API.

## Features

- Create, read, update, and delete tasks
- High-performance in-memory storage with O(1) operations
- Optimized pagination with fixed page size (100 items)
- RESTful API design
- Comprehensive unit tests
- Docker support

## API Endpoints

- `GET /tasks?page=1` - List tasks with pagination (100 items per page)
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update a task
- `DELETE /tasks/{id}` - Delete a task

### Pagination

The API uses server-controlled pagination with a fixed page size of 100 items. Clients can only specify the page number:

```bash
# Get first page (default)
curl https://task-api.etrex.tw/tasks

# Get specific page
curl https://task-api.etrex.tw/tasks?page=2
```

Response format:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Task name",
      "status": 0
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 100,
    "total": 150,
    "pages": 2,
    "has_next": true,
    "has_prev": false
  }
}
```

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

You can test the API directly using the live deployment:

### Create a task
```bash
curl -X POST https://task-api.etrex.tw/tasks \
  -H "Content-Type: application/json" \
  -d '{"name":"Learn Go","status":0}'
```

### List tasks with pagination
```bash
# Get first page
curl https://task-api.etrex.tw/tasks

# Get specific page
curl https://task-api.etrex.tw/tasks?page=2
```

### Update a task
```bash
curl -X PUT https://task-api.etrex.tw/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Learn Go","status":1}'
```

### Delete a task
```bash
curl -X DELETE https://task-api.etrex.tw/tasks/{id}
```

> **Note**: Replace `{id}` with the actual task ID returned from the create or list operations.

## Storage Architecture

### High-Performance In-Memory Storage

The Task API uses a custom-designed in-memory storage system optimized for both performance and memory efficiency. The storage layer is built with two complementary data structures:

#### Data Structure Design

1. **Primary Storage**: A `slice` (dynamic array) that maintains insertion order
2. **Index Mapping**: A `map[string]int` that provides O(1) UUID-to-index lookups
3. **Concurrent Access**: Protected by `sync.RWMutex` for thread-safe operations

```go
type MemoryStorage struct {
    mu        sync.RWMutex
    tasks     []model.Task      // Preserves insertion order
    indexMap  map[string]int    // UUID -> slice index mapping
}
```

#### Time Complexity Analysis

| Operation | Time Complexity | Description |
|-----------|-----------------|-------------|
| **Create** | O(1) | Append to slice + update index map |
| **Read** | O(1) | Direct index lookup via map |
| **Update** | O(1) | Direct index access via map |
| **Delete** | O(1) | Swap-and-pop technique |
| **List (Paginated)** | O(limit) | Direct slice access (max 100 items) |

#### Key Optimizations

1. **Insertion Order Preservation**: Tasks are stored in the order they were created, eliminating the need for sorting
2. **Efficient Deletion**: Uses the "swap-and-pop" technique - moves the last element to the deleted position to avoid O(n) array shifting
3. **Fixed-Size Pagination**: Server-controlled pagination with 100 items per page ensures consistent O(100) performance
4. **Memory Locality**: Sequential data access in the slice provides excellent CPU cache performance

#### Storage Benefits

- **Predictable Performance**: All operations have guaranteed time complexity
- **Memory Efficient**: No duplicate data storage or complex tree structures
- **Scalable**: Performance remains consistent regardless of data size
- **Thread-Safe**: Full concurrent read/write support with minimal locking

### Pagination Strategy

The API implements server-controlled pagination to optimize performance:

- **Fixed Page Size**: 100 items per page (non-configurable by clients)
- **Stateless**: Each page request is independent
- **Efficient**: Direct slice access without scanning entire dataset
- **Consistent**: Response time remains constant regardless of total data size

## Performance Benchmarks

This API has been thoroughly tested for high-concurrency performance with excellent results.

**Note**: The performance metrics below were obtained from local testing with `kern.ipc.somaxconn=1024`. Cloud Run performance may vary due to different system configurations and network latency.

### Key Performance Metrics
- **Maximum Stable Concurrency**: 3,000 concurrent requests with 99.83% success rate
- **Average Response Time**: 1.67ms at 3,000 concurrent requests
- **Storage Layer**: 100,000 concurrent operations with 100% success rate
- **Memory Usage**: Stable at 50-58MB under maximum load
- **CPU Utilization**: Efficient multi-core usage (195% on 8-core system)
- **Pagination Performance**: Consistent O(100) response time regardless of dataset size

### Stress Testing

#### Local Performance Testing

Run comprehensive performance tests with detailed monitoring:

```bash
# Start the API server
go run main.go &

# Run stress test suite
go run benchmark/stress_benchmark.go
```

The stress test includes:
- Storage layer concurrent operations test
- HTTP API progressive load testing
- Mixed read/write workload simulation
- Real-time resource monitoring (CPU, memory, network connections)

#### Cloud Performance Testing

For testing the live GCP Cloud Run deployment, use the web-based benchmark tool:

**🌐 [Open Performance Benchmark Tool](https://etrex.github.io/task-api/benchmark.html)**

Features:
- **Real-time monitoring**: Live charts showing success rate and response times
- **Configurable testing**: Adjust concurrency levels, request types, and timeouts
- **Cloud Run optimized**: Designed for single-instance testing (max-instances=1)
- **Cross-platform**: Works in any modern browser without setup
- **Visual results**: Interactive charts and detailed performance summaries

The web benchmark automatically tests against the live API at `https://task-api.etrex.tw` and provides the same progressive load testing as the Go version, but with a visual interface.

### System Requirements for Optimal Performance

**⚠️ Important**: To achieve maximum performance, the system TCP listen queue must be increased:

```bash
# Check current setting (default is usually 128)
sysctl kern.ipc.somaxconn

# Increase to 1024 for optimal performance
sudo sysctl -w kern.ipc.somaxconn=1024
```

**Without this adjustment**, you'll encounter connection limits at around 240 concurrent requests due to TCP connection reset errors.

### Performance Bottlenecks Identified

1. **240 concurrent requests** - Limited by default `kern.ipc.somaxconn=128`
2. **3,000+ concurrent requests** - Limited by TCP connection establishment rate, not application performance
3. **Application layer** - No bottlenecks found; excellent performance characteristics

### Conclusion

The Task API demonstrates excellent performance characteristics with no application-level bottlenecks. The identified limits are system-level TCP constraints, indicating a well-optimized application design.