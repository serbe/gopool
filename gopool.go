package gopool

import "sync"

// Pool - specification of gopool
type Pool struct {
	taskQueueChan  chan interface{}
	taskWorkerChan chan interface{}
	handler        func(interface{})
	workers        []*worker
}

type worker struct {
	id        int
	isWorking bool
	handler   func(interface{})
	taskChan  *chan interface{}
}

type taskMap struct {
	mr    sync.RWMutex
	tasks []interface{}
}

type queue struct {
	toQueue   chan interface{}
	toWorkers chan interface{}
	taskMap   taskMap
}

// NewPool - create new gorourine pool
// num - max workers
func NewPool(num int, handler func(interface{})) *Pool {
	taskQueueChan := make(chan interface{}, 1)
	taskWorkerChan := make(chan interface{}, 1)
	pool := &Pool{
		taskQueueChan:  taskQueueChan,
		taskWorkerChan: taskWorkerChan,
		handler:        handler,
		workers:        make([]*worker, 4),
	}
	return pool
}

func (p *Pool) start(num int) {
	for i := 0; i < num; i++ {
		p.workers[i] = newWorker(i, p.handler, &p.taskWorkerChan)
		go p.workers[i].start()
	}
}

// AddTask - add new task to pool
func (p *Pool) AddTask(t interface{}) {
	p.taskQueueChan <- t
}

func newWorker(id int, handler func(interface{}), taskWorkerChan *chan interface{}) *worker {
	w := &worker{
		id:       id,
		handler:  handler,
		taskChan: taskWorkerChan,
	}
	return w
}

func (w *worker) start() {
	func() {
		for {
			select {
			case task := <-*w.taskChan:
				w.isWorking = true
				w.handler(task)
			}
		}
	}()
}

func newQueue(taskQueueChan chan interface{}, taskWorkerChan chan interface{}) *queue {
	q := &queue{
		toQueue:   taskQueueChan,
		toWorkers: taskWorkerChan,
	}
	return q
}
