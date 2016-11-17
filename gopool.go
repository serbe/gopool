package gopool

import (
	"container/list"
	"fmt"
	"sync"
)

// Pool - specification of gopool
type Pool struct {
	numWorkers int

	useResultChan bool

	workersRunning bool
	managerRunning bool

	addedTasks    int
	runningTasks  int
	completeTasks int

	workersWg sync.WaitGroup
	managerWg sync.WaitGroup
	tasksWg   sync.WaitGroup

	taskPool         *list.List
	completeTaskPool *list.List

	workersQuitChan chan bool
	managerQuitChan chan bool

	addTaskSignal  chan bool
	doneTaskSignal chan *Task
	resultChan     chan *Task
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int) *Pool {
	pool := new(Pool)
	pool.numWorkers = numWorkers
	pool.taskPool = list.New()
	pool.completeTaskPool = list.New()
	pool.workersQuitChan = make(chan bool)
	pool.managerQuitChan = make(chan bool)
	pool.addTaskSignal = make(chan bool)
	pool.doneTaskSignal = make(chan *Task)
	pool.resultChan = make(chan *Task)
	return pool
}

// Run - start pool
func (p *Pool) Run() error {
	if p.workersRunning {
		return fmt.Errorf("workers already running")
	}
	for i := 0; i < p.numWorkers; i++ {
		p.workersWg.Add(1)
		go p.worker(i)
	}
	p.workersRunning = true
	if p.managerRunning {
		return fmt.Errorf("wanager already running")
	}
	go p.manager()
	p.managerRunning = true
	return nil
}

// Status - return addedTasks, runningTasks ant completeTasks
func (p *Pool) Status() (int, int, int) {
	return p.addedTasks, p.runningTasks, p.completeTasks
}

// Wait - wait to finish all tasks
func (p *Pool) Wait() {
	p.tasksWg.Wait()
}

// Quit - send quit
func (p *Pool) Quit() {
	close(p.workersQuitChan)
	close(p.managerQuitChan)
}
