package gopool

import (
	"math"
	"testing"
)

func testFunc(args ...interface{}) interface{} {
	x := args[0].(float64)
	j := 0.
	for i := 1.0; i < 10000000; i++ {
		j += math.Sqrt(i)
	}
	return x*x + j
}

func TestNew(t *testing.T) {
	pool := New(2)
	pool.Run()

	numWorkers := pool.numWorkers
	if numWorkers != 2 {
		t.Errorf("%v != %v", 2, numWorkers)
	}

	pool.Add(testFunc, float64(2))
	task := pool.GetTask()

	result := task.Result.(float64)
	if result != float64(2.1081849490439312e+10) {
		t.Errorf("%v != %v", float64(2.1081849490439312e+10), result)
	}

	addedTasks, runningTasks, completeTasks := pool.Status()
	if addedTasks != 1 || runningTasks != 0 || completeTasks != 1 {
		t.Errorf("%v != %v || %v != %v || %v != %v", 1, addedTasks, 0, runningTasks, 1, completeTasks)
	}

	pool.Run()
	pool.Add(testFunc, float64(3))

	allresult := pool.Results()
	if len(allresult) != 1 {
		t.Errorf("%v != %v", 1, len(allresult))
	}
	pool.Wait()
	pool.Quit()
}

func TestPool_Run(t *testing.T) {
	tests := []struct {
		name string
		p    *Pool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.p.Run()
	}
}

func TestPool_Add(t *testing.T) {
	type args struct {
		f    func(...interface{}) interface{}
		args []interface{}
	}
	tests := []struct {
		name string
		p    *Pool
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.p.Add(tt.args.f, tt.args.args...)
	}
}

func TestPool_Status(t *testing.T) {
	tests := []struct {
		name  string
		p     *Pool
		want  int
		want1 int
		want2 int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, got1, got2 := tt.p.Status()
		if got != tt.want {
			t.Errorf("%q. Pool.Status() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. Pool.Status() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
		if got2 != tt.want2 {
			t.Errorf("%q. Pool.Status() got2 = %v, want %v", tt.name, got2, tt.want2)
		}
	}
}

func TestPool_Wait(t *testing.T) {
	tests := []struct {
		name string
		p    *Pool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.p.Wait()
	}
}
