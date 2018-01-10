package gopool

import (
	"errors"
	"fmt"
	"time"
)

var (
	errNilFn   = errors.New("error: function is nil")
	errNotRun  = errors.New("error: pool is not running")
	errTimeout = errors.New("error: timed out")
)

// Task - task
type Task struct {
	ID       int64
	WorkerID int64
	Fn       func(...interface{}) interface{}
	Result   interface{}
	Args     []interface{}
	Error    error
}

// Add - add new task to pool
func (p *Pool) Add(fn func(...interface{}) interface{}, args ...interface{}) error {
	if fn == nil {
		return errNilFn
	}
	if !p.poolIsRunning() {
		return errNotRun
	}
	task := &Task{
		Fn:   fn,
		Args: args,
	}
	p.inputTaskChan <- task
	return nil
}

func (p *Pool) addTask(task *Task) {
	if p.getFreeWorkers() > 0 {
		p.decWorkers()
		p.workChan <- task
	} else {
		p.queue.put(task)
	}
}

func (p *Pool) tryGetTask() {
	task, ok := p.queue.get()
	if ok {
		p.decWorkers()
		p.workChan <- task
	}
}

// SetTaskTimeout - set task timeout in second before send quit signal
func (p *Pool) SetTaskTimeout(t int) {
	p.quitTimeout = time.Duration(t) * time.Second
	p.useTimeout = true
}

func (p *Pool) exec(task *Task) *Task {
	defer func() {
		err := recover()
		if err != nil {
			task.Result = nil
			task.Error = fmt.Errorf("Recovery %v", err.(string))
		}
	}()
	if p.useTimeout {
		ch := make(chan interface{}, 1)

		go func() {
			ch <- task.Fn(task.Args...)
			defer close(ch)
		}()

		select {
		case result := <-ch:
			task.Result = result
		case <-time.After(1 * time.Second):
			task.Error = errTimeout
		}
	} else {
		task.Result = task.Fn(task.Args...)
	}
	return task
}
