package gopool

import (
	"sync/atomic"
)

func (p *Pool) worker(id int64) {
	for task := range p.workChan {
		task.WorkerID = id
		task = p.exec(task)
		if p.poolIsRunning() {
			p.ResultChan <- task
			p.endTaskChan <- true
		} else {
			break
		}
	}
}

func (p *Pool) runWorkers() {
	var i int64
	for i = 0; i < p.numWorkers; i++ {
		go p.worker(i)
	}
}

// GetFreeWorkers - get num of free workers
func (p *Pool) GetFreeWorkers() int64 {
	return atomic.LoadInt64(&p.freeWorkers)
}

func (p *Pool) incWorkers() {
	atomic.AddInt64(&p.freeWorkers, 1)
}

func (p *Pool) decWorkers() {
	atomic.AddInt64(&p.freeWorkers, -1)
}
