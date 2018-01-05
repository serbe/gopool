package gopool

import (
	"sync/atomic"
)

func (p *Pool) worker(id int32) {
	for task := range p.workChan {
		task.WorkerID = id
		task = p.exec(task)
		p.ResultChan <- task
		p.endTaskChan <- true
	}
}

func (p *Pool) runWorkers() {
	for i := 0; i < int(p.numWorkers); i++ {
		go p.worker(int32(i))
	}
}

// GetFreeWorkers - get num of free workers
func (p *Pool) GetFreeWorkers() int32 {
	return atomic.LoadInt32(&p.freeWorkers)
}

func (p *Pool) incWorkers() {
	atomic.AddInt32(&p.freeWorkers, 1)
}

func (p *Pool) decWorkers() {
	atomic.AddInt32(&p.freeWorkers, -1)
}
