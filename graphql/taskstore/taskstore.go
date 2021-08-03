package taskstore

import (
	"fmt"
	"sync"
	"time"

	"github.com/shien/restserver/graphql/graph/model"
)

type TaskStore struct {
	sync.Mutex

	tasks  map[int]*model.Task
	nextID int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]*model.Task)
	ts.nextID = 0

	return ts
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time, attachments []*model.Attachment) int {
	ts.Lock()
	defer ts.Unlock()

	task := &model.Task{
		ID:          ts.nextID,
		Text:        text,
		Due:         due,
		Attachments: attachments}

	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextID] = task
	ts.nextID++

	return task.ID
}

func (ts *TaskStore) GetTask(id int) (*model.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]

	if ok {
		return task, nil
	} else {
		return nil, fmt.Errorf("task with id = %d not found", id)
	}
}

func (ts *TaskStore) GetAllTasks() []*model.Task {
	ts.Lock()
	defer ts.Unlock()

	// var allTasks []*model.Task
	allTasks := make([]*model.Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id = %d not found", id)
	}

	delete(ts.tasks, id)

	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]*model.Task)

	return nil
}

func (ts *TaskStore) GetTaskByTag(tag string) []*model.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []*model.Task

TaskLoop:
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue TaskLoop
			}
		}
	}

	return tasks
}

func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []*model.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []*model.Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()

		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
