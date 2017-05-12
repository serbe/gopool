package gopool

import (
	"errors"
	"sync"
)

var (
	errWorkers = errors.New("workers already running")
)

// Pool - specification of gopool
type Pool struct {
	numWorkers     int
	useResultChan  bool
	workersRunning bool
	managerRunning bool
	addedTasks     int
	runningTasks   int
	completeTasks  int

	workersWg sync.WaitGroup
	managerWg sync.WaitGroup
	tasksWg   sync.WaitGroup

	waitingTaskList  *tList
	completeTaskList *tList

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
	pool.waitingTaskList = new(tList)
	pool.completeTaskList = new(tList)
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
		return errWorkers
	}
	for i := 0; i < p.numWorkers; i++ {
		p.workersWg.Add(1)
		go p.worker(i)
	}
	p.workersRunning = true
	// if p.managerRunning {
	// 	return errManager
	// }
	go p.manager()
	p.managerRunning = true
	return nil
}

// Status - return addedTasks, runningTasks ant completeTasks
func (p *Pool) Status() (int, int, int) {
	return p.addedTasks, p.runningTasks, p.completeTasks
}

// WaitAll - wait to finish all tasks
func (p *Pool) WaitAll() {
	p.tasksWg.Wait()
}

// Quit - send quit
func (p *Pool) Quit() {
	close(p.workersQuitChan)
	close(p.managerQuitChan)
}

// ResultChan - return chan of result tasks
func (p *Pool) ResultChan(open bool) *chan *Task {
	if open {
		p.useResultChan = true
	} else {
		p.useResultChan = false
	}
	return &p.resultChan
}

// Done - check all task done
func (p *Pool) Done() bool {
	return p.addedTasks > 0 && p.runningTasks == 0 && p.addedTasks == p.completeTasks
}
