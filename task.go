package gopool

import (
	"fmt"
	"log"
)

// Task - task
type Task struct {
	WorkerID int
	F        func(...interface{}) interface{}
	Result   interface{}
	Args     []interface{}
	Err      error
}

// Add - add task to pool
func (p *Pool) Add(f func(...interface{}) interface{}, args ...interface{}) {
	if f != nil && args != nil {
		task := new(Task)
		task.F = f
		task.Args = args
		_ = p.waitingTaskList.put(task)
		p.addedTasks++
		p.addTaskSignal <- true
	}
}

// Results - return all complete tasks and clear old results
func (p *Pool) Results() []*Task {
	results := p.completeTaskList.val
	p.completeTaskList = new(tList)
	return results
}

func (p *Pool) exec(t *Task) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Panic while running task:", err)
			t.Result = nil
			t.Err = fmt.Errorf("Recovery %v", err.(string))
		}
	}()
	t.Result = t.F(t.Args...)
}
