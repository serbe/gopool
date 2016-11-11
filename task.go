package gopool

// Task - task
type Task struct {
	f       func(...interface{}) interface{}
	result  interface{}
	args    []interface{}
	err     error
	confirm chan bool
}
