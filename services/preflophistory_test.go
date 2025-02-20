package services

import "testing"

func TestIsNotTerminal(t *testing.T) {
	hist := PreFlopHistory{history: []string{"ABC"}}
	result := hist.isTerminal()
	if !(false == result) {
		t.Fatalf(`IsTerminal(%q) = %t, want match for %t, nil`, hist, result, false)
	}
}

func TestIsTerminal(t *testing.T) {
	hist := PreFlopHistory{history: []string{"1h2n3m4l5n"}}
	result := hist.isTerminal()
	if !(true == result) {
		t.Fatalf(`IsTerminal(%q) = %t, want match for %t, nil`, hist, result, true)
	}
}

func TestGameStageEnded(t *testing.T) {
	histLastCall := PreFlopHistory{history: []string{"1h2n3m4l5n", "asd", "a2", "c"}}
	histLastFold := PreFlopHistory{history: []string{"1h2n3m4l5n", "f"}}
	histLastCheck := PreFlopHistory{history: []string{"1h2n3m4l5n", "asd", "c", "k"}}
	if result := histLastCall.gameStageEnded(); !result {
		t.Fatalf(`IsGameStageEnded(%q) = %t, want match for %t, nil`, histLastCall, result, true)
	}
	if result := histLastFold.gameStageEnded(); !result {
		t.Fatalf(`IsGameStageEnded(%q) = %t, want match for %t, nil`, histLastFold, result, true)
	}
	if result := histLastCheck.gameStageEnded(); !result {
		t.Fatalf(`IsGameStageEnded(%q) = %t, want match for %t, nil`, histLastCheck, result, true)
	}
}

func TestGameStageNotEnded(t *testing.T) {
	histLastCall := PreFlopHistory{history: []string{"1h2n3m4l5n", "c"}}
	if result := histLastCall.gameStageEnded(); result {
		t.Fatalf(`IsGameStageNotEnded(%q) = %t, want match for %t, nil`, histLastCall, result, true)
	}
}

func TestPlayer(t *testing.T) {
	histIsChance := PreFlopHistory{history: []string{"1h2n3m4l5n"}}
	histLastFold := PreFlopHistory{history: []string{"1h2n3m4l5n", "f"}}
	histEnded := PreFlopHistory{history: []string{"1h2n3m4l5n", "/"}}
	histPlayerOneTurn := PreFlopHistory{history: []string{"3s4s", "KhKs", "c"}}
	histPlayerTwoTurn := PreFlopHistory{history: []string{"3s4s", "KhKs"}}
	if result := histIsChance.player(); result != -1 {
		t.Fatalf(`Player(%q) = %b, want match for %b, nil`, histIsChance, result, -1)
	}
	if result := histLastFold.player(); result != -1 {
		t.Fatalf(`Player(%q) = %b, want match for %b, nil`, histLastFold, result, -1)
	}
	if result := histEnded.player(); result != -1 {
		t.Fatalf(`Player(%q) = %b, want match for %b, nil`, histEnded, result, -1)
	}
	if result := histPlayerOneTurn.player(); result != 0 {
		t.Fatalf(`Player(%q) = %b, want match for %b, nil`, histPlayerOneTurn, result, 0)
	}
	if result := histPlayerTwoTurn.player(); result != 1 {
		t.Fatalf(`Player(%q) = %b, want match for %b, nil`, histPlayerTwoTurn, result, 1)
	}
}

func TestGetTotalPotSize(t *testing.T) {
	betMin := PreFlopHistory{history: []string{"3s4s", "KhKs", "bMIN"}}
	if stageTotal, latestBet := betMin.getTotalPotSize(); stageTotal != 5 && latestBet != 3 {
		t.Fatalf(`GetTotalPotSize(%q) = %d, %d, want match for %d, nil`, betMin, stageTotal, latestBet, 4)
	}
	betMinBetMid := PreFlopHistory{history: []string{"3s4s", "KhKs", "bMIN", "bMID"}}
	if stageTotal, latestBet := betMinBetMid.getTotalPotSize(); stageTotal != 13 && latestBet != 10 {
		t.Fatalf(`GetTotalPotSize(%q) = %d, %d, want match for %d, nil`, betMinBetMid, stageTotal, latestBet, 4)
	}
	betMax := PreFlopHistory{history: []string{"3s4s", "KhKs", "bMIN", "c"}}
	if stageTotal, latestBet := betMax.getTotalPotSize(); stageTotal != 5 && latestBet != 3 {
		t.Fatalf(`GetTotalPotSize(%q) = %d, %d, want match for %d, nil`, betMax, stageTotal, latestBet, 4)
	}

}
