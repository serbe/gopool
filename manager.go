package gopool

func (p *Pool) manager() {
runLoop:
	for {
		select {
		case work := <-p.inputChan:
			_ = p.queue.put(work)
			p.popTask()
		case <-p.doneTaskSignalChan:
			p.popTask()
		// case <-time.After(t10ms):
		// 	if p.free() > 0 {
		// 		if p.queue.length() > 0 {
		// 			work, _ := p.queue.get()
		// 			// if err == nil {
		// 			p.workChan <- work
		// 			// } else {
		// 			// log.Println("Error in p.queue.get", err)
		// 			// }
		// 		} else if p.completeTasks > 0 && p.completeTasks == p.addedTasks {
		// 			if p.addedTasks == 1 && time.Since(p.startTime) > timeout || p.addedTasks != 1 {
		// 				close(p.ResultChan)
		// 				break runLoop
		// 			}
		// 		}
		// 	}
		case <-p.managerQuitChan:
			break runLoop
		}
	}

	// 	p.managerWg.Add(1)
	// 	defer p.managerWg.Done()
	// managerLoop:
	// 	for {
	// 		select {
	// 		case <-p.addTaskSignal:
	// 			p.tasksWg.Add(1)
	// 		case task := <-p.doneTaskSignal:
	// 			p.runningTasks--
	// 			_ = p.completeTaskList.put(task)
	// 			p.tasksWg.Done()
	// 			p.completeTasks++
	// 		case <-p.managerQuitChan:
	// 			break managerLoop
	// 		}
	// 	}
}
