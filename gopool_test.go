package gopool

import (
	"testing"
	"time"
)

var numWorkers = 4

func testFunc(args ...interface{}) interface{} {
	return args[0].(int) * args[0].(int)
}

func testFunc2(args ...interface{}) interface{} {
	time.Sleep(time.Duration(1) * time.Second)
	return nil
}

func Test1(t *testing.T) {
	p := New(numWorkers)
	if p.numWorkers != numWorkers {
		t.Errorf("Got %v numWorkers, want %v", p.numWorkers, numWorkers)
	}
	p.SetTaskTimeout(1)
	err := p.Add(nil, 1)
	if err != errNilFn {
		t.Errorf("Got %v error, want %v", err, errNilFn)
	}
	err = p.Add(testFunc, 1)
	if err != nil {
		t.Errorf("Got %v error, want %v", err, nil)
	}
	result := <-p.ResultChan
	if result.Result != 1 {
		t.Errorf("Got %v result, want %v", result.Result, 1)
	}
	for i := 0; i < numWorkers+2; i++ {
		p.Add(testFunc2)
	}
	p.tryGetTask()
	p.Quit()
}

func BenchmarkAccumulate(b *testing.B) {
	p := New(numWorkers)
	n := b.N
	for i := 0; i < n; i++ {
		p.Add(testFunc, i)
	}
	// for i := 0; i < n; i++ {
	// 	<-p.ResultChan
	// }
}
