package gopool

import (
	"container/list"
	"fmt"
	"log"
	"sync"
)

// Pool - specification of gopool
type Pool struct {
	numWorkers     int
	isRunning      bool
	workersRunning bool
	managerRunning bool

	addedTasks    int
	runningTasks  int
	completeTasks int

	workersWg sync.WaitGroup
	managerWg sync.WaitGroup
	tasksWg   sync.WaitGroup

	waitTaskList     *list.List
	completeTaskList *list.List

	workersQuitChan chan bool
	managerQuitChan chan bool
	addTaskChan     chan *Task
	doneTaskChan    chan *Task
	wantedTaskChan  chan chan *Task
	ResultChan      chan interface{}
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int) *Pool {
	pool := new(Pool)
	pool.numWorkers = numWorkers
	pool.waitTaskList = list.New()
	pool.completeTaskList = list.New()
	pool.workersQuitChan = make(chan bool)
	pool.managerQuitChan = make(chan bool)
	pool.addTaskChan = make(chan *Task)
	pool.doneTaskChan = make(chan *Task)
	pool.wantedTaskChan = make(chan chan *Task)
	pool.ResultChan = make(chan interface{})

	return pool
}

// Run - start pool
func (p *Pool) Run() {
	if p.workersRunning {
		log.Println("Workers already running")
	} else {
		for i := 0; i < p.numWorkers; i++ {
			p.workersWg.Add(1)
			go p.worker(i)
		}
		p.workersRunning = true
	}
	if p.managerRunning {
		log.Println("Manager already running")
	} else {
		go p.manager()
		p.managerRunning = true
	}
}

// Add - add task to pool
func (p *Pool) Add(f func(...interface{}) interface{}, args ...interface{}) {
	task := new(Task)
	task.f = f
	task.args = args
	task.confirm = make(chan bool)
	p.addTaskChan <- task
	<-task.confirm
}

// Status - return addedTasks, runningTasks ant completeTasks
func (p *Pool) Status() (int, int, int) {
	return p.addedTasks, p.runningTasks, p.completeTasks
}

// Wait - wait to finish all tasks
func (p *Pool) Wait() {
	p.tasksWg.Wait()
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

func (p *Pool) exec(t *Task) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic while running job:", err)
			t.result = nil
			t.err = fmt.Errorf(err.(string))
		}
	}()
	t.result = t.f(t.args...)
}
