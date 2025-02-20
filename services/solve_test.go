package services

import (
	"reflect"
	"testing"
)

func Test_runPreflopIteration(t *testing.T) {
	type args struct {
		history PreFlopHistory
	}
	tests := []struct {
		name string
		args args
		want map[string]InfoSet
	}{
		{
			name: "iterate once",
			args: args{
				history: PreFlopHistory{
					history:      []string{},
					playerHand:   "9h4d",
					opponentHand: "2hAs",
					board:        "2c9c3sKhKs",
					winner:       "1",
				},
			},
			want: map[string]InfoSet{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if runPreflopIteration(tt.args.history); !reflect.DeepEqual(InfoSetsPreFlop, tt.want) {
					t.Errorf("runPreflopIteration() = %v, want %v", InfoSetsPreFlop, tt.want)
				}
			},
		)
	}
}
