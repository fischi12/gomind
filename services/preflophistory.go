package services

import (
	"slices"
	"strconv"
)

var discreteActions []string = []string{"k", "bMIN", "bMID", "bMAX", "c", "f"}
var InfoSetsPreFlop = map[string]InfoSet{}

type PreFlopHistory struct {
	history      []string
	winner       string
	playerHand   string
	opponentHand string
	board        string
}

func (history PreFlopHistory) isTerminal() bool {
	if len(history.history) == 0 {
		return false
	}
	if len(history.history[len(history.history)-1]) == 10 {
		return true
	}
	return false
}

func (history PreFlopHistory) actions() []string {
	if history.isChance() {
		return []string{}
	}

	if len(history.history) == 2 {
		return []string{"c", "bMIN", "bMID", "bMAX", "f"}
	} else if history.history[len(history.history)-1] == "bMIN" {
		return []string{"bMID", "bMAX", "f", "c"}
	} else if history.history[len(history.history)-1] == "bMID" {
		return []string{"bMAX", "f", "c"}
	} else if history.history[len(history.history)-1] == "bMAX" {
		return []string{"f", "c"}
	} else {
		return []string{"k", "bMIN", "bMID", "bMAX"}
	}
}

func (history PreFlopHistory) player() int {
	if len(history.history) < 2 {
		return -1
	} else if history.gameStageEnded() {
		return -1
	} else if history.history[len(history.history)-1] == "/" {
		return -1
	} else {
		return (len(history.history) + 1) % 2
	}

}

func (history PreFlopHistory) gameStageEnded() bool {
	if history.history[len(history.history)-1] == "c" && len(history.history) > 3 {
		return true
	}
	if history.history[len(history.history)-1] == "f" {
		return true
	}
	if slices.Equal(history.history[len(history.history)-2:], []string{"c", "k"}) {
		return true
	}
	return false
}

func (history PreFlopHistory) isChance() bool {
	return history.player() == -1
}

func (history PreFlopHistory) sampleChanceOutcome() string {
	if !history.isChance() {
		panic("Try to sample chance on no chance situation")
	}
	if len(history.history) == 0 {
		return history.playerHand
	} else if len(history.history) == 1 {
		return history.opponentHand
	} else if history.history[len(history.history)-1] != "/" {
		return "/"
	} else {
		return history.board
	}
}

func (history PreFlopHistory) terminalUtility(player int) float64 {
	winner := history.winner
	potSize, latestBet := history.getTotalPotSize()
	if slices.Contains(history.history, "f") {
		foldIdx := slices.Index(history.history, "f")
		histBeforeFlop := &PreFlopHistory{history: history.history[:foldIdx-1]}
		potSize, latestBet = histBeforeFlop.getTotalPotSize()
		if history.history[len(history.history)-3] == "bMIN" || history.history[len(history.history)-3] == "bMID" {
			potSize += latestBet
		}
		if len(history.history)%2 == player {
			return float64(-potSize) / 2.0
		} else {
			return float64(potSize) / 2
		}
	}
	if winner == "0" {
		return 0
	}
	if (winner == "1" && player == 0) || (winner == "-1" && player == 1) {
		return float64(potSize) / 2
	} else {
		return float64(-potSize) / 2
	}
}

func (history PreFlopHistory) getTotalPotSize() (int, int) {
	stageTotal := 3
	latestBet := 2

	for _, a := range history.history {
		if a == "bMIN" {
			oldStageTotal := stageTotal
			stageTotal = latestBet + stageTotal
			latestBet = oldStageTotal
		} else if a == "bMID" {
			oldStageTotal := stageTotal
			stageTotal = latestBet + 2*stageTotal
			latestBet = 2 * oldStageTotal
		} else if a == "bMAX" {
			stageTotal = latestBet + 100
			latestBet = 100
		} else if a == "c" {
			stageTotal = 2 * latestBet
		}

	}
	return stageTotal, latestBet
}

func (history PreFlopHistory) getInfoSetKey() string {
	if history.isChance() {
		panic("No InfoSet for chance node")
	}
	if history.isTerminal() {
		panic("No InfoSet for terminal node")
	}
	player := history.player()
	infoSetKey := ""
	infoSetKey += strconv.Itoa(getPreFlopClusterId(history.history[player]))
	for _, v := range history.history {
		if slices.Contains(discreteActions, v) {
			infoSetKey += v
		}
	}
	return infoSetKey
}

func (history PreFlopHistory) getInfoSet() InfoSet {
	infoSetKey := history.getInfoSetKey()
	actions := history.actions()
	player := history.player()
	_, exists := InfoSetsPreFlop[infoSetKey]
	if !exists {
		infoSet := InfoSet{
			InfoSetKey:         infoSetKey,
			Actions:            actions,
			Player:             player,
			Strategy:           map[string]float64{},
			Regret:             map[string]float64{},
			CumulativeStrategy: map[string]float64{},
		}
		infoSet.GetStrategy()
		InfoSetsPreFlop[infoSetKey] = infoSet
	}
	return InfoSetsPreFlop[infoSetKey]
}

func (history PreFlopHistory) newHistory(newHistory []string) history {
	return PreFlopHistory{
		history: newHistory, winner: history.winner, playerHand: history.playerHand,
		opponentHand: history.opponentHand, board: history.board,
	}
}

func (history PreFlopHistory) getHistory() []string {
	return history.history
}
