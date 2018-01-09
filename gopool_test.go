package gopool

import (
	"testing"
	"time"
)

var numWorkers int64 = 4

func testFunc(args ...interface{}) interface{} {
	return args[0].(int) * args[0].(int)
}

func testFunc2(args ...interface{}) interface{} {
	time.Sleep(time.Duration(args[0].(int)) * time.Second)
	return nil
}

func Test1(t *testing.T) {
	p := New(numWorkers)
	if p.numWorkers != numWorkers {
		t.Errorf("Got %v numWorkers, want %v", p.numWorkers, numWorkers)
	}
	// p.SetTaskTimeout(1)
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
	for i := 0; i < int(numWorkers+2); i++ {
		err = p.Add(testFunc, i)
		if err != nil {
			t.Errorf("Got %v error, want %v", err, nil)
		}
	}
	for i := 0; i < int(numWorkers+2); i++ {
		task := <-p.ResultChan
		if task.Error != nil {
			t.Errorf("Got %v error, want %v", task.Error, nil)
		}
	}
	p.Quit()
	err = p.Add(testFunc2, 3)
	if err != errNotRun {
		t.Errorf("Got %v error, want %v", err, errNotRun)
	}
}

func TestTimeout(t *testing.T) {
	p := New(numWorkers)
	p.SetTaskTimeout(1)
	p.Add(testFunc2, 2)
	p.Add(testFunc2, 2)
	p.Add(testFunc2, 4)
	for i := 0; i < 3; i++ {
		task := <-p.ResultChan
		if task.Error != errTimeout {
			t.Errorf("Got %v error, want %v", task.Error, errTimeout)
		}
	}
}

func BenchmarkAccumulate(b *testing.B) {
	p := New(numWorkers)
	n := b.N
	for i := 0; i < n; i++ {
		err := p.Add(testFunc, i)
		if err != nil {
			println("Error", err)
		}
	}
	for i := 0; i < n; i++ {
		task := <-p.ResultChan
		_ = task.Result.(int)
	}
}

func BenchmarkTimeout(b *testing.B) {
	p := New(numWorkers)
	p.SetTaskTimeout(1)
	n := b.N
	for i := 0; i < n; i++ {
		err := p.Add(testFunc, i)
		if err != nil {
			println("Error", err)
		}
	}
	for i := 0; i < n; i++ {
		task := <-p.ResultChan
		_ = task.Result.(int)
	}
}

func BenchmarkParallel(b *testing.B) {
	p := New(numWorkers)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := p.Add(testFunc, 1)
			if err != nil {
				b.Errorf("Got %v error, want %v", err, nil)
			}
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := <-p.ResultChan
			res := task.Result.(int)
			if res != 1 {
				b.Errorf("Got %v error, want %v", res, 1)
			}
		}
	})
}
