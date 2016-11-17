package gopool

import (
	"container/list"
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
	task := new(Task)
	task.F = f
	task.Args = args
	p.taskPool.PushBack(task)
	p.addedTasks++
	p.addTaskSignal <- true
}

// Results - return all complete tasks and clear old results
func (p *Pool) Results() []*Task {
	results := make([]*Task, p.completeTaskPool.Len())
	i := 0
	for elem := p.completeTaskPool.Front(); elem != nil; elem = elem.Next() {
		results[i] = elem.Value.(*Task)
		i++
	}
	p.completeTaskPool = list.New()
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
