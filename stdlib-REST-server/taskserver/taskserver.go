package taskserver

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shien/restserver/taskstore"
)

// Backend server wraps the database like taskstore
type TaskServer struct {
	Datastore *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	store := taskstore.New()
	return &TaskServer{Datastore: store}
}

// Handler function for routing and HTTP multiplexer in golang standard lib
func (ts *TaskServer) createTaskHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling create a task at %s\n", req.URL.Path)

	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type RequestTaskID struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		http.Error(rsp, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var rt RequestTask

	if err := decoder.Decode(&rt); err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.Datastore.CreateTask(rt.Text, rt.Tags, rt.Due)

	MarshalAndPrepareHTTPResponse(RequestTaskID{Id: id}, rsp)
}

func (ts *TaskServer) deleteTaskHandler(rsp http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Handling delete a task at %s\n", req.URL.Path)

	err := ts.Datastore.DeleteTask(id)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusNotFound)
	}
}

func (ts *TaskServer) deleteAllTasksHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling delete all tasks at %s\n", req.URL.Path)
	ts.Datastore.DeleteAllTasks()
}

func (ts *TaskServer) getAllTasksHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling get all tasks at %s\n", req.URL.Path)

	allTasks := ts.Datastore.GetAllTasks() // 1. backend service

	MarshalAndPrepareHTTPResponse(allTasks, rsp)
}

func (ts *TaskServer) getTaskHandler(rsp http.ResponseWriter, req *http.Request, id int) {
	log.Printf("Handling get a task at %s\n", req.URL.Path)

	task, err := ts.Datastore.GetTask(id)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusNotFound)
		return
	}

	MarshalAndPrepareHTTPResponse(task, rsp) // 2. Prepare the HTTP response to client
}

// handler that sees what REST API should be provided and pass the request to the low-level handlers
func (ts *TaskServer) TaskHandler(rsp http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		if req.Method == http.MethodPost {
			ts.createTaskHandler(rsp, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(rsp, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTasksHandler(rsp, req)
		} else {
			http.Error(rsp,
				fmt.Sprintf("Expect method GET, DELETE or POST at /task/, got %v", req.Method),
				http.StatusMethodNotAllowed)
			return
		}

	} else { // handler requests like /task/<id>
		pathParts := TrimAndParseRequestPath(*req)

		if len(pathParts) < 2 {
			http.Error(rsp, "Expect /task/<id> in task handler function", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(pathParts[1])

		if err != nil {
			http.Error(rsp, "Expect /task/<id> in task handler function", http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTaskHandler(rsp, req, id)
		} else if req.Method == http.MethodGet {
			ts.getTaskHandler(rsp, req, id)
		} else {
			http.Error(rsp,
				fmt.Sprintf("Expect method GET, DELETE or POST at /task/<id>, got %v", req.Method),
				http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *TaskServer) TagHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling tasks by tag at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(rsp,
			fmt.Sprintf("Expect method GET at /tag/<tag>, got %v", req.Method),
			http.StatusMethodNotAllowed)
		return
	}

	pathParts := TrimAndParseRequestPath(*req)

	if len(pathParts) < 2 {
		http.Error(rsp, "Expect /tag/<tag> in tag handler function", http.StatusBadRequest)
		return
	}

	tag := pathParts[1]

	task := ts.Datastore.GetTaskByTag(tag)

	MarshalAndPrepareHTTPResponse(task, rsp)
}

func (ts *TaskServer) DueHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling tasks by due at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(rsp,
			fmt.Sprintf("Expect method GET at /due/<date>, got %v", req.Method),
			http.StatusMethodNotAllowed)
		return
	}

	pathParts := TrimAndParseRequestPath(*req)

	prepareBadRequestError := func() {
		http.Error(rsp,
			fmt.Sprintf("Expect method GET at /due/<year>/<month>/<day>, got %v", req.Method),
			http.StatusBadRequest)
	}

	if len(pathParts) != 4 {
		prepareBadRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])

	if err != nil {
		prepareBadRequestError()
		return
	}

	month, err := strconv.Atoi(pathParts[2])

	if err != nil || (month > 12) || (month < 1) {
		prepareBadRequestError()
		return
	}

	day, err := strconv.Atoi(pathParts[3])

	if err != nil {
		prepareBadRequestError()
		return
	}

	tasks := ts.Datastore.GetTaskByDueDate(year, time.Month(month), day)
	MarshalAndPrepareHTTPResponse(tasks, rsp)
}

func TrimAndParseRequestPath(req http.Request) []string {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	return pathParts
}

func MarshalAndPrepareHTTPResponse(task interface{}, rsp http.ResponseWriter) {
	js, err := json.Marshal(task)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
		return
	}

	rsp.Header().Set("Content-Type", "application/json")
	rsp.Write(js)
}
