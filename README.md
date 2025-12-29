# Temporal Workflow Example

This project demonstrates a simple Temporal workflow for job processing with an API interface.

## 1. Prerequisites

- Go 1.25

- Temporal

- Gin

## 2. How to start Temporal server

```bash

# Start Temporal server in development mode
$ temporal server start-dev
```

This command starts:

- Temporal Frontend Service

- Temporal History Service

- Temporal Matchmaker Service

- Temporal Worker Service

## 3. How to start the worker and the API

### Start Worker

```bash

# Navigate to worker directory
$ cd worker

# Run worker
$ go run main.go
```

### Start API Server

```bash

# Navigate to watch directory
$ cd watch

# Run API server
$ go run main.go
```

The API server will start on port 8080 by default.

## 4. How to verify all core requirements

### Using Postman

1. **Create a new job**

- Endpoint: `POST http://localhost:8080/jobs`

- Request Body:

```json

{
    "input": "greeting-workflow", // jobid
    "options": "0"      // 1 means failure
}
```

- Response:

```json

{
    "job_id": "greeting-workflow"
}
```

2.**Get job status**

- Endpoint: `GET http://localhost:8080/jobs/greeting-workflow`

- Response:

```json

{
    "attempt": 1,
    "error": null,
    "job_id": "greeting-workflow",
    "result": "Hello ",
    "status": "finished"
}
```

## 5. Key Design Decisions

### API Endpoints

- **POST /jobs** - Create a new job
        

    - Accepts job parameters and returns job ID

    - Validates input parameters

    - Handles job submission and state management

- **GET /jobs/{jobId}** - Get job status
        

    - Retrieves job details from the workflow state

    - Returns current status and result (if completed)

### Workflow Implementation

- **workerflow.go**

    - Maintains a map of job IDs to their state

    - Uses read-write locks to prevent concurrent modifications

    - Implements job scheduling and execution logic

    - Updates job status before and after execution

- **activity.go**

    - Simulates a time-consuming operation (e.g., data processing)

    - Represents a single task that can be executed by a worker

    - Returns processing result and duration

### Worker Implementation

- **worker/main.go**

    - Registers workflows and activities with Temporal client

    - Starts Temporal worker to process tasks

    - Handles workflow execution and state management

    - Controlled failure can set first attemp failed 
    s
    - Automatic retry for most 3 times if failed

### API Implementation

- **watch/main.go**

    - Uses Gin framework for HTTP server

    - Provides RESTful endpoints for job management

    - Communicates with Temporal worker for job status updates

## 6. Project Structure

```bash

.
├── activity.go         # Activity definitions
├── docker-compose.yml
├── dockerfile
├── go.mod              # Go module file
├── go.sum 
├── query
│   └── query.go        # deal with query
├── readme.md
├── start
│   └── main.go         # start in terminal
├── watch
│   └── main.go         # API server
├── worker
│   └── main.go         # Worker
└── workflow.go         # Workflow definitions
```

## 7. AI Usage Disclosure

- Used AI tools to learn Temporal framework concepts

- Received guidance on Dockerfile creation

- Gained insights on workflow design patterns

- Received suggestions for API endpoint design
