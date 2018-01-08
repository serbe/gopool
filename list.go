package gopool

import (
	"sync"
	"sync/atomic"
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

func (p *Pool) put(task *Task) {
	p.list.Put(task)
}

func (p *Pool) get() (*Task, bool) {
	t := p.list.Get()
	if t != nil {
		return t.(*Task), true
	}
	return &Task{}, false
}

func (tasks *taskList) get() (*Task, bool) {
	var task *Task
	if tasks.length() > 0 {
		tasks.Lock()
		task = tasks.list[0]
		tasks.list = tasks.list[1:]
		tasks.len--
		tasks.Unlock()
		return task, true
	}
	return task, false
}

func (tasks *taskList) length() int64 {
	return atomic.LoadInt64(&tasks.len)
}
