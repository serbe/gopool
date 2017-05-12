package gopool

import (
	"fmt"
	"sync"
)

type tList struct {
	m   sync.Mutex
	len int
	val []*Task
}

func (t *tList) put(task *Task) error {
	t.m.Lock()
	defer t.m.Unlock()
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	t.val = append(t.val, task)
	t.len++
	return nil
}

func (t *tList) get() (*Task, error) {
	t.m.Lock()
	defer t.m.Unlock()
	var task *Task
	if t.len > 0 {
		task = t.val[0]
		t.len--
		t.val = t.val[1:]
		return task, nil
	}
	return nil, fmt.Errorf("tList is empty")
}
