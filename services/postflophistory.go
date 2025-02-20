package services

import (
	"slices"
)

var infoSetsPostFlop = map[string]InfoSet{}

type PostFlopHistory struct {
	history              []string
	winner               string
	playerHand           string
	opponentHand         string
	board                string
	playerFlopCluster    string
	opponentFlopCluster  string
	playerTurnCluster    string
	opponentTurnCluster  string
	playerRiverCluster   string
	opponentRiverCluster string
}

func (history PostFlopHistory) getStage() int {
	count := 0
	for _, action := range history.history {
		if action == "/" {
			count++
		}
	}
	return count
}

func (history PostFlopHistory) isTerminal() bool {
	if len(history.history) == 0 {
		return false
	}
	folded := history.history[len(history.history)-1] == "f"
	isShowdown := history.getStage() == 3 && history.gameStageEnded()
	if folded || isShowdown {
		return true
	} else {
		return false
	}
}

func (history PostFlopHistory) actions() []string {
	if history.isChance() {
		return []string{}
	}

	if history.history[len(history.history)-1] == "k" {
		return []string{"k", "bMIN", "bMAX"}
	} else if slices.Equal(history.history[len(history.history)-2:], []string{"k", "bMIN"}) {
		return []string{"f", "c"}
	} else if history.history[len(history.history)-1] == "bMIN" {
		return []string{"bMAX", "f", "c"}
	} else if history.history[len(history.history)-1] == "bMAX" {
		return []string{"f", "c"}
	} else {
		return []string{"k", "bMIN", "bMAX"}
	}
}

func (history PostFlopHistory) player() int {
	if len(history.history) < 3 {
		return -1
	} else if history.gameStageEnded() {
		return -1
	} else if history.history[len(history.history)-1] == "/" {
		return -1
	} else {
		lastGameStage := history.getLastGameStage()
		return (len(lastGameStage) + 1) % 2
	}

}

func (history PostFlopHistory) gameStageEnded() bool {
	if history.history[len(history.history)-1] == "c" {
		return true
	}
	if history.history[len(history.history)-1] == "f" {
		return true
	}
	if slices.Equal(history.history[len(history.history)-2:], []string{"k", "k"}) {
		return true
	}
	return false
}

func (history PostFlopHistory) isChance() bool {
	return history.player() == -1
}

func (history PostFlopHistory) sampleChanceOutcome() string {
	if !history.isChance() {
		panic("Try to sample chance on no chance situation")
	}
	if len(history.history) == 0 {
		return history.playerHand
	} else if len(history.history) == 1 {
		return history.opponentHand
	} else if history.history[len(history.history)-1] != "/" {
		return "/"
	} else if history.getStage() == 1 {
		return history.board[:6]
	} else if history.getStage() == 2 {
		return history.board[6:8]
	} else if history.getStage() == 3 {
		return history.board[8:10]
	}
	panic("Something went wrong")
}

func (history PostFlopHistory) getLastGameStage() []string {
	var lastGameStageStartIdx = -1

	for i, val := range history.history {
		if val == "/" {
			lastGameStageStartIdx = i
		}
	}

	if lastGameStageStartIdx == -1 {
		return []string{}
	}

	return history.history[lastGameStageStartIdx+1:]
}

func (history PostFlopHistory) terminalUtility(player int) float64 {
	winner := history.winner
	potSize, latestBet := history.getTotalPotSize()
	lastGameStage := history.getLastGameStage()
	if history.history[len(history.history)-1] == "f" {
		histBeforeFlop := &PostFlopHistory{history: history.history[:len(history.history)-2]}
		potSize, latestBet = histBeforeFlop.getTotalPotSize()
		if history.history[len(history.history)-3] == "bMIN" {
			potSize += latestBet
		}
		if len(lastGameStage)%2 == player {
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

func (history PostFlopHistory) getTotalPotSize() (int, int) {
	total := 0
	stageTotal := 4
	latestBet := 0

	for _, a := range history.history {
		if a == "/" {
			total += stageTotal
			stageTotal = 0
			latestBet = 0
		} else if a == "bMIN" {
			latestBet = max(2, total/3)
			stageTotal += latestBet
		} else if a == "bMAX" {
			latestBet = total
			stageTotal += latestBet
		} else if a == "c" {
			stageTotal = 2 * latestBet
		}

	}
	total += stageTotal
	return total, latestBet
}

func (history PostFlopHistory) getInfoSetKey() string {
	if history.isChance() {
		panic("No InfoSet for chance node")
	}
	if history.isTerminal() {
		panic("No InfoSet for terminal node")
	}
	stageI := 0

	player := history.player()
	infoSetKey := ""
	for _, a := range history.history {
		if !slices.Contains(discreteActions, a) {
			if a == "/" {
				stageI += 1
				continue
			}
			if stageI == 1 {
				if player == 0 {
					infoSetKey += history.playerFlopCluster
				} else {
					infoSetKey += history.opponentFlopCluster
				}
			} else if stageI == 2 {
				if player == 0 {
					infoSetKey += history.playerTurnCluster
				} else {
					infoSetKey += history.opponentTurnCluster
				}
			} else if stageI == 3 {
				if player == 0 {
					infoSetKey += history.playerRiverCluster
				} else {
					infoSetKey += history.opponentRiverCluster
				}
			}
		} else {
			infoSetKey += a
		}
	}
	return infoSetKey
}

func (history PostFlopHistory) getInfoSet() InfoSet {
	infoSetKey := history.getInfoSetKey()
	actions := history.actions()
	player := history.player()
	_, exists := infoSetsPostFlop[infoSetKey]
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
		infoSetsPostFlop[infoSetKey] = infoSet
	}
	return infoSetsPostFlop[infoSetKey]
}

func (history PostFlopHistory) newHistory(newHistory []string) history {
	return PostFlopHistory{
		history:              newHistory,
		winner:               history.winner,
		playerHand:           history.playerHand,
		opponentHand:         history.opponentHand,
		board:                history.board,
		playerFlopCluster:    history.playerFlopCluster,
		opponentFlopCluster:  history.opponentFlopCluster,
		playerTurnCluster:    history.playerTurnCluster,
		opponentTurnCluster:  history.opponentTurnCluster,
		playerRiverCluster:   history.playerRiverCluster,
		opponentRiverCluster: history.opponentRiverCluster,
	}
}

func (history PostFlopHistory) getHistory() []string {
	return history.history
}
