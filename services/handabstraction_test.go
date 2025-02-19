package services

import (
	"reflect"
	"testing"
)

func Test_calculateHandStrength(t *testing.T) {
	type args struct {
		holeCards []string
		board     []string
	}
	tests := []struct {
		name string
		args args
		want handStrength
	}{
		{
			name: "pocket_asse",
			args: args{holeCards: []string{"Ah", "As"}, board: []string{"2s", "3d", "Ks", "Qc", "7c"}},
			want: handStrength{Wins: 884, Loss: 105, Draws: 1},
		},
		{
			name: "royal_flush",
			args: args{holeCards: []string{"Ah", "Kh"}, board: []string{"Jh", "Th", "Qh", "Qc", "7c"}},
			want: handStrength{Wins: 990, Loss: 0, Draws: 0},
		},
		{
			name: "royal_flush_turn",
			args: args{holeCards: []string{"Ah", "Kh"}, board: []string{"Jh", "Th", "Qh", "Qc"}},
			want: handStrength{Wins: 45540, Loss: 0, Draws: 0},
		},
		{
			name: "pocket_asse_turn",
			args: args{holeCards: []string{"Ah", "As"}, board: []string{"2s", "3d", "Ks", "Qc"}},
			want: handStrength{Wins: 40062, Loss: 5434, Draws: 44},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := calculateHandStrength(tt.args.holeCards, tt.args.board); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("calculateHandStrength() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_generateCombinations(t *testing.T) {
	type args struct {
		deck            []string
		combinationSize int
		start           int
		current         []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "basic combination size 2",
			args: args{deck: []string{"A", "B", "C", "D"}, combinationSize: 2, start: 0, current: []string{}},
			want: [][]string{
				{"A", "B"},
				{"A", "C"},
				{"A", "D"},
				{"B", "C"},
				{"B", "D"},
				{"C", "D"},
			},
		},
		{
			name: "basic combination size 3",
			args: args{deck: []string{"A", "B", "C", "D"}, combinationSize: 3, start: 0, current: []string{}},
			want: [][]string{
				{"A", "B", "C"},
				{"A", "B", "D"},
				{"A", "C", "D"},
				{"B", "C", "D"},
			},
		},
		{
			name: "basic combination move start",
			args: args{deck: []string{"A", "B", "C", "D"}, combinationSize: 2, start: 1, current: []string{}},
			want: [][]string{
				{"B", "C"},
				{"B", "D"},
				{"C", "D"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := generateCombinations(
					tt.args.deck,
					tt.args.combinationSize,
					tt.args.start,
					tt.args.current,
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generateCombinations() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_generateAbstractedHoleCards(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "generate hole cards", want: 169},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := len(generateAbstractedHoleCards()); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generateAbstractedHoleCards() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_generateAbstractedBoards(t *testing.T) {
	type args struct {
		holeCards           []string
		countCommunityCards int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "abstract board",
			args: args{holeCards: []string{"3c", "3d"}, countCommunityCards: 3},
			want: 11124,
		},
		{
			name: "abstract board suited hole cards",
			args: args{holeCards: []string{"3d", "4d"}, countCommunityCards: 3},
			want: 5035,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := len(
					generateAbstractedBoards(
						tt.args.holeCards,
						tt.args.countCommunityCards,
					),
				); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generateAbstractedBoards() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_generateHandCombinations(t *testing.T) {
	type args struct {
		countCommunityCards int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "flop hand combinations", args: args{countCommunityCards: 3}, want: 1374984},
		{name: "turn hand combinations", args: args{countCommunityCards: 4}, want: 15633514},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := len(GenerateHandCombinations(tt.args.countCommunityCards)); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generateHandCombinations() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_abstractHand(t *testing.T) {
	type args struct {
		holeCards      []string
		communityCards []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "sort hole and community cards", args: args{
				holeCards: []string{"2c", "2d"}, communityCards: []string{
					"2h", "3h", "3c",
				},
			}, want: []string{"2c", "2d", "2h", "3c", "3h"},
		},
		{
			name: "basic abstraction", args: args{
				holeCards: []string{"2c", "2d"}, communityCards: []string{
					"3c", "3d", "5s",
				},
			}, want: []string{"2c", "2d", "3c", "3d", "5h"},
		},
		{
			name: "basic abstraction", args: args{
				holeCards: []string{"2c", "2d"}, communityCards: []string{
					"3c", "3d", "5h",
				},
			}, want: []string{"2c", "2d", "3c", "3d", "5h"},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := abstractHand(tt.args.holeCards, tt.args.communityCards); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("abstractHand() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
