package gopool

import (
	"container/list"
	"fmt"
	"log"
)

// Task - task
type Task struct {
	F       func(...interface{}) interface{}
	Result  interface{}
	Args    []interface{}
	Err     error
	confirm chan bool
}

// Results - return all complete tasks and clear old results
func (p *Pool) Results() []*Task {
	results := make([]*Task, p.completeTaskList.Len())
	i := 0
	for elem := p.completeTaskList.Front(); elem != nil; elem = elem.Next() {
		results[i] = elem.Value.(*Task)
		i++
	}
	p.completeTaskList = list.New()
	return results
}

// GetTask - get next result
func (p *Pool) GetTask() *Task {
	p.useResultChan = true
	task := new(Task)
gettask:
	for {
		select {
		case task = <-p.resultChan:
			break gettask
		default:
		}
	}
	p.useResultChan = false
	return task
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
