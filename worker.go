package gopool

import "time"

var ms10 = time.Duration(10) * time.Millisecond

func (p *Pool) runWorker(id int) {
	for task := range p.workChan {
		p.dec()
		task.WorkerID = id
		p.exec(task)
		p.doneTaskSignalChan <- true
		p.ResultChan <- *task
		p.inc()
	}
}

func (p *Pool) free() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.freeWorkers
}

func (p *Pool) inc() {
	p.m.Lock()
	p.freeWorkers++
	p.completeTasks++
	p.m.Unlock()
}

func (p *Pool) dec() {
	p.m.Lock()
	p.freeWorkers--
	p.m.Unlock()
}
