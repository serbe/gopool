package gopool

import (
	"testing"
)

var numWorkers = 4

func testFunc(args ...interface{}) interface{} {
	return args[0].(int) * args[0].(int)
}

func Test1(t *testing.T) {
	p := New(numWorkers)
	if p.numWorkers != numWorkers {
		t.Fatalf("Found %v number of workers, want %v", p.numWorkers, numWorkers)
	}
	err := p.Run()
	if err != errWorkers {
		t.Fatal("No have error in already start workers")
	}
	var addedTasks, runningTasks, completeTasks int
	added, running, complete := p.Status()
	if added != addedTasks {
		t.Fatal("Wrong number of added tasks")
	}
	if running != runningTasks {
		t.Fatal("Wrong number of running tasks")
	}
	if complete != completeTasks {
		t.Fatal("Wrong number of complete tasks")
	}
	if p.Done() != false {
		t.Fatal("Wrong done status")
	}
	p.ResultChan(false)
	if p.useResultChan != false {
		t.Fatal("Wrong status of result chan")
	}
	p.ResultChan(true)
	if p.useResultChan != true {
		t.Fatal("Wrong status of result chan")
	}
	p.Add(testFunc, 1)
}
