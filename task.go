package gopool

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var errNilFn = errors.New("error: function is nil")

// Task - task
type Task struct {
	ID       int
	WorkerID int
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
	task := Task{
		Fn:   fn,
		Args: args,
	}
	p.inputTaskChan <- task
	if !p.isRunning && p.autorun {
		p.TryGetTask()
	}
	return nil
}

func (p *Pool) addTask(task Task) {
	if p.getFreeWorkers() > 0 {
		if p.timerIsRunning {
			p.timer.Stop()
		}
		p.decWorkers()
		p.workChan <- task
	} else {
		p.queue.put(task)
	}
}

// TryGetTask - try to get task from queue
func (p *Pool) TryGetTask() {
	if p.freeWorkers > 0 {
		p.isRunning = true
		task, ok := p.queue.get()
		if ok {
			if p.timerIsRunning {
				p.timer.Stop()
			}
			p.decWorkers()
			p.workChan <- task
		}
	}
}

// SetTaskTimeout - set task timeout in second before send quit signal
func (p *Pool) SetTaskTimeout(t int) {
	p.quitTimeout = time.Duration(t) * time.Second
	p.timer = time.NewTimer(p.quitTimeout)
	p.timerIsRunning = true
	go func() {
		<-p.timer.C
		p.quit <- true
	}()
}

func (p *Pool) exec(task Task) Task {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Panic while running task:", err)
			task.Result = nil
			task.Error = fmt.Errorf("Recovery %v", err.(string))
		}
	}()
	task.Result = task.Fn(task.Args...)
	return task
}
