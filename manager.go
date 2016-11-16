package gopool

func (p *Pool) manager() {
	p.managerWg.Add(1)
	defer p.managerWg.Done()
manager:
	for {
		select {
		case task := <-p.addTaskChan:
			p.waitTaskList.PushBack(task)
			p.tasksWg.Add(1)
			p.addedTasks++
			task.confirm <- true
		case taskChan := <-p.wantedTaskChan:
			elem := p.waitTaskList.Front()
			if elem != nil {
				task := elem.Value.(*Task)
				p.waitTaskList.Remove(elem)
				taskChan <- task
			}
		case doneTask := <-p.doneTaskChan:
			p.runningTasks--
			p.completeTaskList.PushBack(doneTask)
			p.tasksWg.Done()
			p.completeTasks++
		case <-p.managerQuitChan:
			break manager
		default:
			// case <-time.After(10 * time.Millisecond):
		}
	}
}
