package main

import "testing"

func Test_enqueuTask(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "enqueuTask"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				enqueuTask()
			},
		)
	}
}
