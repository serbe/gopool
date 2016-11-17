package gopool

func (p *Pool) manager() {
	p.managerWg.Add(1)
	defer p.managerWg.Done()
	for {
		select {
		case <-p.addTaskSignal:
			p.tasksWg.Add(1)
		case task := <-p.doneTaskSignal:
			p.runningTasks--
			p.completeTaskPool.PushBack(task)
			p.tasksWg.Done()
			p.completeTasks++
		case <-p.managerQuitChan:
			break
			// default:
			// case <-time.After(10 * time.Millisecond):
		}
	}
}
