package gopool

import (
	"sync"
)

type taskList struct {
	sync.RWMutex
	list []Task
}

func (tasks *taskList) put(task Task) {
	tasks.Lock()
	tasks.list = append(tasks.list, task)
	tasks.Unlock()
}

func (tasks *taskList) get() (Task, bool) {
	tasks.Lock()
	var task Task
	if len(tasks.list) > 0 {
		task = tasks.list[0]
		tasks.list = tasks.list[1:]
		tasks.Unlock()
		return task, true
	}
	tasks.Unlock()
	return task, false
}
