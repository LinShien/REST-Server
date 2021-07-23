package taskserver

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/shien/restserver/stdlib-REST-server/taskserver"
	"github.com/shien/restserver/taskstore"
)

// Backend server wraps the database like taskstore
type TaskServerForRouter struct {
	Datastore *taskstore.TaskStore
}

func NewTaskServerForRouter() *TaskServerForRouter {
	store := taskstore.New()

	return &TaskServerForRouter{Datastore: store}
}

// Handler function for routing and HTTP multiplexer in golang standard lib
func (ts *TaskServerForRouter) CreateTaskHandler(rsp http.ResponseWriter, req *http.Request) {
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

	taskserver.MarshalAndPrepareHTTPResponse(RequestTaskID{Id: id}, rsp)
}

func (ts *TaskServerForRouter) DeleteTaskHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling delete a task at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	err := ts.Datastore.DeleteTask(id)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusNotFound)
	}
}

func (ts *TaskServerForRouter) DeleteAllTasksHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling delete all tasks at %s\n", req.URL.Path)
	ts.Datastore.DeleteAllTasks()
}

func (ts *TaskServerForRouter) GetAllTasksHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling get all tasks at %s\n", req.URL.Path)

	allTasks := ts.Datastore.GetAllTasks() // 1. backend service

	taskserver.MarshalAndPrepareHTTPResponse(allTasks, rsp)
}

func (ts *TaskServerForRouter) GetTaskHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling get a task at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	task, err := ts.Datastore.GetTask(id)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusNotFound)
		return
	}

	taskserver.MarshalAndPrepareHTTPResponse(task, rsp) // 2. Prepare the HTTP response to client
}

func (ts *TaskServerForRouter) TagHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling get a task by tag at %s\n", req.URL.Path)

	tag := mux.Vars(req)["tag"]

	tasks := ts.Datastore.GetTaskByTag(tag)

	taskserver.MarshalAndPrepareHTTPResponse(tasks, rsp)
}

func (ts *TaskServerForRouter) DueHandler(rsp http.ResponseWriter, req *http.Request) {
	log.Printf("Handling get a task by due at %s\n", req.URL.Path)

	vars := mux.Vars(req)

	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	day, _ := strconv.Atoi(vars["day"])

	tasks := ts.Datastore.GetTaskByDueDate(year, time.Month(month), day)

	taskserver.MarshalAndPrepareHTTPResponse(tasks, rsp)
}
