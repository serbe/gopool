package gopool

import (
	"container/list"
	"log"
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

	waitTaskList     *list.List
	completeTaskList *list.List

	workersQuitChan chan bool
	managerQuitChan chan bool
	addTaskChan     chan *Task
	doneTaskChan    chan *Task
	wantedTaskChan  chan chan *Task
	resultChan      chan *Task
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
	pool.resultChan = make(chan *Task)

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
	task.F = f
	task.Args = args
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
