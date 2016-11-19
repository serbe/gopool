package gopool

import "time"

var ms10 = time.Duration(10) * time.Millisecond

func (p *Pool) worker(id int) {
	defer p.workersWg.Done()
workerLoop:
	for {
		select {
		case <-time.After(ms10):
			if p.waitingTaskList.len > 0 {
				task, err := p.waitingTaskList.get()
				if err == nil {
					task.WorkerID = id
					p.runningTasks++
					p.exec(task)
					p.doneTaskSignal <- task
					if p.useResultChan {
						p.resultChan <- task
					}
				}
			}
		case <-p.workersQuitChan:
			break workerLoop
		}
	}
}
