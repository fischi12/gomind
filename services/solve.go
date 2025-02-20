package services

import (
	"github.com/mattlangl/gophe"
	"github.com/schollz/progressbar/v3"
	"math/rand"
	"strconv"
	"strings"
)

func runIteration(history history) {
	vanillaCfr(history, 0, 1, 1)
	vanillaCfr(history, 1, 1, 1)
}

func runPostFlopIteration() {
	flopHistory := generatePostFlopHistory()
	runIteration(flopHistory)
}

func runPreFlopIteration() {
	flopHistory := generatePreFlopHistory()
	runIteration(flopHistory)
}

func trainPreFlopIteration() {
	iterations := 1000000
	bar := progressbar.New(iterations)
	for range iterations {
		bar.Add(1)
		runPreFlopIteration()
	}

}

func trainPostFlopIteration() {
	iterations := 1000000
	bar := progressbar.New(iterations)
	for range iterations {
		bar.Add(1)
		runPostFlopIteration()
	}

}

func generatePreFlopHistory() PreFlopHistory {
	deck := initializeDeck()

	rand.Shuffle(
		len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		},
	)

	playerHand := gophe.NewHand(
		gophe.NewCard(deck[0]),
		gophe.NewCard(deck[1]),
		gophe.NewCard(deck[4]),
		gophe.NewCard(deck[5]),
		gophe.NewCard(deck[6]),
		gophe.NewCard(deck[7]),
		gophe.NewCard(deck[8]),
	)
	opponentHand := gophe.NewHand(
		gophe.NewCard(deck[2]),
		gophe.NewCard(deck[3]),
		gophe.NewCard(deck[4]),
		gophe.NewCard(deck[5]),
		gophe.NewCard(deck[6]),
		gophe.NewCard(deck[7]),
		gophe.NewCard(deck[8]),
	)

	playerRank := gophe.EvaluateHand(*playerHand).GetValue()
	opponentRank := gophe.EvaluateHand(*opponentHand).GetValue()

	winner := "1"

	if playerRank == opponentRank {
		winner = "0"
	}

	if opponentRank < playerRank {
		winner = "-1"
	}

	return PreFlopHistory{
		history:      []string{},
		playerHand:   strings.Join(deck[:2], ""),
		opponentHand: strings.Join(deck[2:4], ""),
		board:        strings.Join(deck[4:9], ""),
		winner:       winner,
	}

}

func generatePostFlopHistory() PostFlopHistory {
	deck := initializeDeck()

	rand.Shuffle(
		len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		},
	)

	playerHand := gophe.NewHand(
		gophe.NewCard(deck[0]),
		gophe.NewCard(deck[1]),
		gophe.NewCard(deck[4]),
		gophe.NewCard(deck[5]),
		gophe.NewCard(deck[6]),
		gophe.NewCard(deck[7]),
		gophe.NewCard(deck[8]),
	)
	opponentHand := gophe.NewHand(
		gophe.NewCard(deck[2]),
		gophe.NewCard(deck[3]),
		gophe.NewCard(deck[4]),
		gophe.NewCard(deck[5]),
		gophe.NewCard(deck[6]),
		gophe.NewCard(deck[7]),
		gophe.NewCard(deck[8]),
	)

	playerRank := gophe.EvaluateHand(*playerHand).GetValue()
	opponentRank := gophe.EvaluateHand(*opponentHand).GetValue()

	winner := "1"

	if playerRank == opponentRank {
		winner = "0"
	}

	if opponentRank < playerRank {
		winner = "-1"
	}

	return PostFlopHistory{
		history:              []string{},
		playerHand:           strings.Join(deck[:2], ""),
		opponentHand:         strings.Join(deck[2:4], ""),
		board:                strings.Join(deck[4:9], ""),
		winner:               winner,
		playerFlopCluster:    strconv.Itoa(rand.Intn(10) + 1),
		opponentFlopCluster:  strconv.Itoa(rand.Intn(10) + 1),
		playerTurnCluster:    strconv.Itoa(rand.Intn(10) + 1),
		opponentTurnCluster:  strconv.Itoa(rand.Intn(10) + 1),
		playerRiverCluster:   strconv.Itoa(rand.Intn(10) + 1),
		opponentRiverCluster: strconv.Itoa(rand.Intn(10) + 1),
	}

}
