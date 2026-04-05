package task

import (
	"encoding/json"
	"io"
	"net/http"

	"example.com/task-manager/internal/db"
	"github.com/google/uuid"
)

func renderJSON(response http.ResponseWriter, body []byte, status int) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	response.Write(body)
}

type Response[T any] struct {
	ErrorMsg string `json:"error"`
	Data     T      `json:"data"`
}

type CreateTaskResponse struct {
	Id string `json:"id"`
}

type Request struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type TaskController struct {
	db *db.MemoryDb[Task]
}

func (t *TaskController) CreateTasks(response http.ResponseWriter, request *http.Request) {

	body, err := io.ReadAll(request.Body)
	if err != nil {
		value, _ := json.Marshal(Response[struct{}]{ErrorMsg: "Can't read body"})
		renderJSON(response, value, http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var r Request
	if err := json.Unmarshal(body, &r); err != nil {
		value, _ := json.Marshal(Response[struct{}]{ErrorMsg: "Can't decode"})
		renderJSON(response, value, http.StatusInternalServerError)
		return
	}

	id := uuid.New()
	t.db.Insert(&Task{
		Id:   id.String(),
		Name: r.Name,
	})

	value, _ := json.Marshal(Response[CreateTaskResponse]{Data: CreateTaskResponse{
		id.String(),
	}})

	renderJSON(response, value, http.StatusCreated)
}

func (t *TaskController) GetTask(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	if id == "" {
		value, _ := json.Marshal(Response[struct{}]{ErrorMsg: "id param is empty"})
		renderJSON(response, value, http.StatusBadRequest)
		return
	}

	element, err := t.db.SelectById(id)
	if err != nil {
		value, _ := json.Marshal(Response[struct{}]{ErrorMsg: "Not found"})
		renderJSON(response, value, http.StatusNotFound)
		return
	}

	value, _ := json.Marshal(Response[Task]{Data: *element})
	renderJSON(response, value, http.StatusOK)
}

func NewTaskController(router *http.ServeMux) {
	if router == nil {
		return
	}
	t := TaskController{
		db: db.NewMemoryDb[Task](),
	}
	router.HandleFunc("POST /tasks/", t.CreateTasks)
	router.HandleFunc("GET /tasks/{id}/", t.GetTask)
}
