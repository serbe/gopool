package gopool

import (
	"reflect"
	"testing"
)

func TestPool_Results(t *testing.T) {
	tests := []struct {
		name string
		p    *Pool
		want []*Task
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.p.Results(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Pool.Results() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPool_GetTask(t *testing.T) {
	tests := []struct {
		name string
		p    *Pool
		want *Task
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.p.GetTask(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Pool.GetTask() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPool_exec(t *testing.T) {
	type args struct {
		t *Task
	}
	tests := []struct {
		name string
		p    *Pool
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.p.exec(tt.args.t)
	}
}
