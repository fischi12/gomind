package services

import (
	"fmt"
	"github.com/mattlangl/gophe"
	"gomind/models"
	"gomind/repository"
	"gorm.io/gorm"
	"sort"
	"strings"
)

type Hand struct {
	HoleCards      []string
	CommunityCards []string
}
type handStrength struct {
	Wins  uint16
	Loss  uint16
	Draws uint16
}

func generateCombinations(deck []string, combinationSize int, start int, current []string) [][]string {
	if len(current) == combinationSize {
		return [][]string{append([]string{}, current...)}
	}

	var result [][]string
	for i := start; i < len(deck); i++ {
		newCurrent := append(current, deck[i])
		result = append(result, generateCombinations(deck, combinationSize, i+1, newCurrent)...)
	}
	return result
}

func initializeDeck() []string {
	var deck []string
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	suits := []string{"h", "d", "s", "c"}

	for _, rank := range ranks {
		for _, suit := range suits {
			deck = append(deck, rank+suit)
		}
	}
	sort.Strings(deck)
	return deck
}

func abstractHand(holeCards []string, communityCards []string) []string {
	suits := []string{"c", "d", "h", "s"}
	sort.Strings(holeCards)
	sort.Strings(communityCards)
	cards := append(holeCards, communityCards...)
	suitMapping := make(map[string]string)
	var result []string

	for _, card := range cards {
		cardSuit := string(card[len(card)-1])
		cardValue := card[:len(card)-1]
		mappedSuit, exists := suitMapping[cardSuit]
		if !exists {
			mappedSuit, suits = suits[0], suits[1:]
			suitMapping[cardSuit] = mappedSuit
		}
		result = append(result, cardValue+mappedSuit)
	}
	return result
}

func removeElement(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeDuplicateCombinations(combinations [][]string) [][]string {
	seen := make(map[string]struct{})
	var result [][]string

	for _, combo := range combinations {
		sortedCombo := append([]string{}, combo...)
		key := fmt.Sprintf("%v", sortedCombo)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, combo)
		}
	}
	return result
}

func GenerateHandCombinations(countCommunityCards int) []Hand {
	result := make([]Hand, 0)
	var holeCards = generateAbstractedHoleCards()

	for _, holeCard := range holeCards {
		boards := generateAbstractedBoards(holeCard, countCommunityCards)
		for _, board := range boards {
			result = append(result, Hand{HoleCards: board[:2], CommunityCards: board[2:]})
		}
	}
	return result
}

func generateAbstractedHoleCards() [][]string {
	deck := initializeDeck()
	var combinations = generateCombinations(deck, 2, 0, []string{})

	var abstractedHoleCards [][]string
	for _, card := range combinations {
		abstractedHand := abstractHand(card, []string{})
		abstractedHoleCards = append(abstractedHoleCards, abstractedHand)
	}

	return removeDuplicateCombinations(abstractedHoleCards)
}

func generateAbstractedBoards(holeCards []string, countCommunityCards int) [][]string {
	deck := initializeDeck()
	for _, card := range holeCards {
		deck = removeElement(deck, card)
	}
	communityCards := generateCombinations(deck, countCommunityCards, 0, []string{})

	var abstractedHands [][]string
	for _, communityCard := range communityCards {
		abstractedHand := abstractHand(holeCards, communityCard)
		abstractedHands = append(abstractedHands, abstractedHand)
	}
	return removeDuplicateCombinations(abstractedHands)
}

func calculateHandStrength(holeCards []string, board []string) handStrength {
	deck := initializeDeck()
	for _, holeCard := range holeCards {
		deck = removeElement(deck, holeCard)
	}
	for _, boardCard := range board {
		deck = removeElement(deck, boardCard)
	}

	if len(board) != 5 {
		countMissingBoardCards := 5 - len(board)
		var missingBoardCards = generateCombinations(deck, countMissingBoardCards, 0, []string{})
		result := handStrength{Wins: 0, Loss: 0, Draws: 0}
		for _, missingBoardCard := range missingBoardCards {
			strength := calculateHandStrength(holeCards, append(board, missingBoardCard...))
			result.Wins += strength.Wins
			result.Loss += strength.Loss
			result.Draws += strength.Draws
		}
		return result
	}

	newHand := gophe.NewHand(
		gophe.NewCard(holeCards[0]),
		gophe.NewCard(holeCards[1]),
		gophe.NewCard(board[0]),
		gophe.NewCard(board[1]),
		gophe.NewCard(board[2]),
		gophe.NewCard(board[3]),
		gophe.NewCard(board[4]),
	)

	rank := gophe.EvaluateHand(*newHand)
	playerRank := rank.GetValue()
	result := handStrength{Wins: 0, Loss: 0, Draws: 0}
	var enemyHoleCards = generateCombinations(deck, 2, 0, []string{})
	for _, enemyHoleCard := range enemyHoleCards {
		newHand := gophe.NewHand(
			gophe.NewCard(enemyHoleCard[0]),
			gophe.NewCard(enemyHoleCard[1]),
			gophe.NewCard(board[0]),
			gophe.NewCard(board[1]),
			gophe.NewCard(board[2]),
			gophe.NewCard(board[3]),
			gophe.NewCard(board[4]),
		)

		rank := gophe.EvaluateHand(*newHand)
		enemyRank := rank.GetValue()
		if enemyRank == playerRank {
			result.Draws++
		}
		if enemyRank < playerRank {
			result.Loss++
		}
		if enemyRank > playerRank {
			result.Wins++
		}
	}
	return result
}

func CalculateAndSaveHandStrengthFlop(batch []Hand, db *gorm.DB) error {
	result := make([]models.FlopHand, 0)
	for _, hand := range batch {
		strength := calculateHandStrength(hand.HoleCards, hand.CommunityCards)
		flopHand := models.FlopHand{
			Hand: strings.Join(hand.HoleCards, "") + strings.Join(
				hand.CommunityCards,
				"",
			), HoleCards: strings.Join(hand.HoleCards, ""),
			Board: strings.Join(
				hand.CommunityCards,
				"",
			),
			Wins: strength.Wins,
			Loss: strength.Loss, Draws: strength.Draws,
		}
		result = append(result, flopHand)
	}
	return repository.UpsertFlopHand(db, &result)
}

func CalculateAndSaveHandStrengthTurn(batch []Hand, db *gorm.DB) error {
	result := make([]models.TurnHand, 0)
	for _, hand := range batch {
		strength := calculateHandStrength(hand.HoleCards, hand.CommunityCards)
		turnHand := models.TurnHand{
			Hand: strings.Join(hand.HoleCards, "") + strings.Join(
				hand.CommunityCards,
				"",
			), HoleCards: strings.Join(hand.HoleCards, ""),
			Board: strings.Join(
				hand.CommunityCards,
				"",
			),
			Wins: strength.Wins,
			Loss: strength.Loss, Draws: strength.Draws,
		}
		result = append(result, turnHand)
	}
	return repository.UpsertTurnHand(db, &result)
}
