package gopool

import (
	"sync"
	"time"
)

var t50ms = time.Duration(50) * time.Millisecond

// Pool - specification of gopool
type Pool struct {
	sync.RWMutex
	timerIsRunning bool
	autorun        bool
	isRunning      bool
	numWorkers     int
	freeWorkers    int
	inputJobs      int
	workChan       chan Task
	inputTaskChan  chan Task
	ResultChan     chan Task
	quit           chan bool
	endTaskChan    chan bool
	queue          taskList
	quitTimeout    time.Duration
	timer          *time.Timer
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.workChan = make(chan Task)
	p.inputTaskChan = make(chan Task)
	p.ResultChan = make(chan Task)
	p.endTaskChan = make(chan bool)
	p.quit = make(chan bool)
	go p.runBroker()
	go p.runWorkers()
	return p
}

func (p *Pool) runBroker() {
loopPool:
	for {
		select {
		case task := <-p.inputTaskChan:
			p.incJobs()
			task.ID = p.getJobs()
			p.addTask(task)
		case <-p.endTaskChan:
			p.incWorkers()
			if p.timerIsRunning && p.getFreeWorkers() == p.numWorkers {
				p.timer.Reset(p.quitTimeout)
			}
		case <-p.quit:
			close(p.workChan)
			close(p.ResultChan)
			break loopPool
		case <-time.After(t50ms):
			p.TryGetTask()
		}
	}
}

func (p *Pool) getJobs() int {
	p.RLock()
	inputJobs := p.inputJobs
	p.RUnlock()
	return inputJobs
}

func (p *Pool) incJobs() {
	p.Lock()
	p.inputJobs++
	p.Unlock()
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	p.quit <- true
}

// Autorun - set auto run get tasks
func (p *Pool) Autorun(flag bool) {
	p.autorun = flag
}
