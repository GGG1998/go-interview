# 🧩 Recruitment Task – Go (40 min)
## 📖 User Story
As a user of a task management system, I want to create and view my tasks so that I can better organize my work.

## 🎯 Functional Requirements
Implement a simple HTTP application in Go that supports:

1. Create a task
Endpoint: POST /tasks
Request (JSON):
{
  "title": "Buy groceries",
  "completed": false
}

Validation:
title must not be empty

Response:
201 Created
Returns the created task with an id

2. Get all tasks
Endpoint: GET /tasks
Response:
200 OK
Returns a list of all tasks

🧱 Technical Requirements
Use Go standard library (net/http)
Store data in memory (slice/map – no database)
Each task should have:
id (int or UUID)
title (string)
completed (bool)

⭐ Bonus (if time permits)
Endpoint: PATCH /tasks/{id}
update completed status
Proper HTTP error handling (e.g. 400, 404)
Basic separation of concerns (handler / service / storage)

🧪 What We Evaluate
Code readability and structure
Error handling
Correctness of REST implementation
Basic Go knowledge (structs, slices, JSON, HTTP)

### Check

```bash
curl -X -v POST http://localhost:8090/tasks/ \
     -H "Content-Type: application/json" \
     -d '{"name": "Learn GO", "completed": false}'
```

### Additional task

- Modify model and due time with notification
- Add Rate Limiter