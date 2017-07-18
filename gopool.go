package gopool

import (
	"errors"
	"sync"
	"time"
)

var (
	errWorkers = errors.New("workers already running")
	errInput   = errors.New("error input values")
	t10ms      = time.Duration(10) * time.Millisecond
	timeout    = time.Duration(5) * time.Second
)

// Pool - specification of gopool
type Pool struct {
	m sync.RWMutex

	numWorkers  int
	freeWorkers int

	startTime time.Time

	inputChan  chan *Task
	workChan   chan *Task
	ResultChan chan Task

	// useResultChan  bool
	workersIsRunning bool
	// managerRunning bool
	addedTasks    int
	completeTasks int

	// runningTasks   int

	// workersWg sync.WaitGroup
	// managerWg sync.WaitGroup
	// tasksWg   sync.WaitGroup

	// waitingTaskList  *tList
	// completeTaskList *tList

	// workersQuitChan chan bool
	managerQuitChan chan bool

	// addTaskSignal  chan bool
	doneTaskSignalChan chan bool
	// resultChan     chan *Task
	queue taskList
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.startTime = time.Now()

	// p.workChan = make(chan *Task, numWorkers)
	p.inputChan = make(chan *Task)
	p.ResultChan = make(chan Task)

	// pool.waitingTaskList = new(tList)
	// pool.completeTaskList = new(tList)
	// pool.workersQuitChan = make(chan bool)
	p.managerQuitChan = make(chan bool)
	// pool.addTaskSignal = make(chan bool)
	p.doneTaskSignalChan = make(chan bool)
	// pool.resultChan = make(chan *Task)
	return p
}

// Run - start pool
func (p *Pool) Run() error {
	if p.workersIsRunning {
		return errWorkers
	}
	for i := 0; i < p.numWorkers; i++ {
		go p.runWorker(i)
	}
	p.workersIsRunning = true
	// if p.managerRunning {
	// 	return errManager
	// }
	go p.manager()
	// p.managerRunning = true
	return nil
}

// Status - return addedTasks, runningTasks ant completeTasks
// func (p *Pool) Status() (int, int, int) {
// 	return p.addedTasks, p.runningTasks, p.completeTasks
// }

// WaitAll - wait to finish all tasks
// func (p *Pool) WaitAll() {
// 	p.tasksWg.Wait()
// }

// Quit - send quit
func (p *Pool) Quit() {
	// close(p.workersQuitChan)
	close(p.managerQuitChan)
}

// ResultChan - return chan of result tasks
// func (p *Pool) ResultChan(open bool) *chan *Task {
// 	if open {
// 		p.useResultChan = true
// 	} else {
// 		p.useResultChan = false
// 	}
// 	return &p.resultChan
// }

// Done - check all task done
// func (p *Pool) Done() bool {
// 	return p.addedTasks > 0 && p.runningTasks == 0 && p.addedTasks == p.completeTasks
// }
