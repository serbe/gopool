package gopool

import (
	"time"
)

var ms10 = time.Duration(10) * time.Millisecond

func (p *Pool) worker(id int) {
	defer p.workersWg.Done()
workerLoop:
	for {
		select {
		case <-time.After(ms10):
			if p.taskPool.Len() > 0 {
				elem := p.taskPool.Front()
				if elem != nil {
					task := elem.Value.(*Task)
					p.taskPool.Remove(elem)
					task.WorkerID = id
					p.runningTasks++
					p.exec(task)
					p.doneTaskSignal <- task
					if p.useResultChan {
						p.resultChan <- task
					}
				} else {
					p.taskPool.Remove(elem)
					p.doneTaskSignal <- nil
					if p.useResultChan {
						p.resultChan <- nil
					}
				}
			}
		case <-p.workersQuitChan:
			break workerLoop
		}
	}
}
