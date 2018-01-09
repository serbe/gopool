package gopool

import (
	"sync"
)

type taskList struct {
	sync.RWMutex
	list []*Task
	len  int64
}

func (tasks *taskList) put(task *Task) {
	tasks.Lock()
	tasks.list = append(tasks.list, task)
	tasks.len++
	tasks.Unlock()
}

func (tasks *taskList) get() (*Task, bool) {
	var task *Task
	tasks.Lock()
	if tasks.len > 0 {
		task = tasks.list[0]
		tasks.list = tasks.list[1:]
		tasks.len--
		tasks.Unlock()
		return task, true
	}
	tasks.Unlock()
	return task, false
}
