package gopool

func (p *Pool) worker(id int) {
	defer p.workersWg.Done()
	taskChan := make(chan *Task)
worker:
	for {
		p.wantedTaskChan <- taskChan

		select {
		case task := <-taskChan:
			p.runningTasks++
			p.exec(task)
			p.doneTaskChan <- task
			if p.useResultChan {
				p.resultChan <- task
			}
		case <-p.workersQuitChan:
			break worker
		default:
			// case <-time.After(10 * time.Millisecond):
		}
	}
}
