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
				if runIteration(tt.args.history); !reflect.DeepEqual(InfoSetsPreFlop, tt.want) {
					t.Errorf("runPreflopIteration() = %v, want %v", InfoSetsPreFlop, tt.want)
				}
			},
		)
	}
}

func Test_runPostflopIteration(t *testing.T) {
	type args struct {
		history PostFlopHistory
	}
	tests := []struct {
		name string
		args args
		want map[string]InfoSet
	}{
		{
			name: "iterate once",
			args: args{
				history: PostFlopHistory{
					history:              []string{},
					playerHand:           "9h4d",
					opponentHand:         "2hAs",
					board:                "2c9c3sKhKs",
					winner:               "1",
					playerFlopCluster:    "1",
					opponentFlopCluster:  "1",
					playerTurnCluster:    "1",
					opponentTurnCluster:  "1",
					playerRiverCluster:   "1",
					opponentRiverCluster: "1",
				},
			},
			want: map[string]InfoSet{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if runIteration(tt.args.history); !reflect.DeepEqual(InfoSetsPreFlop, tt.want) {
					t.Errorf("runPreflopIteration() = %v, want %v", InfoSetsPreFlop, tt.want)
				}
			},
		)
	}
}

func Test_generatePreFlopHistory(t *testing.T) {
	tests := []struct {
		name string
		want PreFlopHistory
	}{
		{name: "Test", want: PreFlopHistory{}},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := generatePreFlopHistory(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generatePreFlopHistory() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

//func Test_trainPreFlopIteration(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		{name: "Train Preflop"},
//	}
//	for _, tt := range tests {
//		t.Run(
//			tt.name, func(t *testing.T) {
//				trainPostFlopIteration()
//			},
//		)
//	}
//}
