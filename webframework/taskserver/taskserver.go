package taskserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shien/restserver/taskstore"
)

// Backend server wraps the database like taskstore
type TaskServerForWebFramework struct {
	Datastore *taskstore.TaskStore
}

func NewTaskServerForWebFramework() *TaskServerForWebFramework {
	store := taskstore.New()

	return &TaskServerForWebFramework{Datastore: store}
}

func (ts *TaskServerForWebFramework) CreateTaskHandler(context *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask

	// Gin will bind the request to the GO struct for us
	// parsing, validating requests and assigning their values to Go structs.
	if err := context.ShouldBindJSON(&rt); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.Datastore.CreateTask(rt.Text, rt.Tags, rt.Due)

	context.JSON(http.StatusOK, gin.H{"Id": id})
}

func (ts *TaskServerForWebFramework) DeleteAllTasksHandler(context *gin.Context) {
	ts.Datastore.DeleteAllTasks()
}

func (ts *TaskServerForWebFramework) DeleteTaskHandler(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))

	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = ts.Datastore.DeleteTask(id); err != nil {
		context.String(http.StatusNotFound, err.Error())
		return
	}
}

func (ts *TaskServerForWebFramework) GetAllTasksHandler(context *gin.Context) {
	tasks := ts.Datastore.GetAllTasks()

	context.JSON(http.StatusOK, tasks)
}

func (ts *TaskServerForWebFramework) GetTaskHandler(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))

	if err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := ts.Datastore.GetTask(id)

	if err != nil {
		context.String(http.StatusNotFound, err.Error())
		return
	}

	context.JSON(http.StatusOK, task)
}

func (ts *TaskServerForWebFramework) TagHandler(context *gin.Context) {
	tag := context.Params.ByName("tag")

	tasks := ts.Datastore.GetTaskByTag(tag)

	context.JSON(http.StatusOK, tasks)
}

func (ts *TaskServerForWebFramework) DueHandler(context *gin.Context) {
	year, err := strconv.Atoi(context.Params.ByName("year"))

	prepareBadRequestError := func() {
		context.String(http.StatusBadRequest,
			fmt.Sprintf("Expect method GET at /due/<year>/<month>/<day>, got %v", context.FullPath()))
	}

	// no regexp support
	if err != nil {
		prepareBadRequestError()
		return
	}

	month, err := strconv.Atoi(context.Params.ByName("month"))

	if err != nil {
		prepareBadRequestError()
		return
	}

	day, err := strconv.Atoi(context.Params.ByName("day"))

	if err != nil {
		prepareBadRequestError()
		return
	}

	// validate the date from client
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	tasks := ts.Datastore.GetTaskByDueDate(date.Year(), date.Month(), date.Day())

	context.JSON(http.StatusOK, tasks)
}
