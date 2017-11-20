package gopool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		task.WorkerID = id
		task = p.exec(task)
		p.ResultChan <- task
		p.endTaskChan <- true
	}
}

func (p *Pool) runWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go p.worker(i)
	}
}

func (p *Pool) getFreeWorkers() int {
	p.RLock()
	freeWorkers := p.freeWorkers
	p.RUnlock()
	return freeWorkers
}

func (p *Pool) incWorkers() {
	p.Lock()
	p.freeWorkers++
	p.Unlock()
}

func (p *Pool) decWorkers() {
	p.Lock()
	p.freeWorkers--
	p.Unlock()
}
