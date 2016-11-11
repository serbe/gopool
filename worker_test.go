package gopool

import "testing"

func TestPool_worker(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		p    *Pool
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.p.worker(tt.args.id)
	}
}
