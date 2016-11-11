package gopool

func (p *Pool) worker(id int) {
	taskChan := make(chan *Task)
worker:
	for {
		p.wantedTaskChan <- taskChan

		select {
		case task := <-taskChan:
			p.runningTasks++
			p.exec(task)
			p.doneTaskChan <- task
		case <-p.workersQuitChan:
			break worker
		default:
		}
	}
}
