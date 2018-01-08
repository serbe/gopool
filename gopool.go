package gopool

import (
	"sync"
	"sync/atomic"
	"time"
)

// var t50ms = time.Duration(50) * time.Millisecond

// Pool - specification of gopool
type Pool struct {
	useTimeout  bool
	runningPool uint32
	// runningTimer  uint32
	numWorkers    int64
	freeWorkers   int64
	inputJobs     int64
	workChan      chan *Task
	inputTaskChan chan *Task
	ResultChan    chan *Task
	quit          chan bool
	endTaskChan   chan bool
	queue         *taskList
	quitTimeout   time.Duration
	list          sync.Pool
	// timer         *time.Timer
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.workChan = make(chan *Task)
	p.inputTaskChan = make(chan *Task)
	p.ResultChan = make(chan *Task)
	p.endTaskChan = make(chan bool)
	p.quit = make(chan bool)
	p.queue = new(taskList)
	p.list = sync.Pool{}
	// 	New: func() interface{} {
	// 		t := new(Task)
	// 		return t
	// 	},
	// }
	go p.runBroker()
	go p.runWorkers()
	p.runningPool = 1
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
			p.TryGetTask()
		case <-p.endTaskChan:
			p.incWorkers()
			// if p.timerIsRunning() && p.GetFreeWorkers() == p.numWorkers {
			// 	p.timer.Reset(p.quitTimeout)
			// }
			p.TryGetTask()
		case <-p.quit:
			close(p.workChan)
			close(p.ResultChan)
			break loopPool
			// case <-time.After(t50ms):
			// 	p.TryGetTask()
		}
	}
}

func (p *Pool) getJobs() int64 {
	return atomic.LoadInt64(&p.inputJobs)
}

func (p *Pool) incJobs() {
	atomic.AddInt64(&p.inputJobs, 1)
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	atomic.StoreUint32(&p.runningPool, 0)
	p.quit <- true
}

func (p *Pool) poolIsRunning() bool {
	return atomic.LoadUint32(&p.runningPool) != 0
}

// func (p *Pool) timerIsRunning() bool {
// 	return atomic.LoadUint32(&p.runningTimer) != 0
// }
