package main

import (
	"encoding/gob"
	"fmt"
	"github.com/mattlangl/gophe"
	"github.com/schollz/progressbar/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type handStrength struct {
	Wins  uint16
	Loss  uint16
	Draws uint16
}

type turnHand struct {
	Hand      string `gorm:"primaryKey"`
	HoleCards string `gorm:"index"`
	Board     string `gorm:"index"`
	Wins      uint16
	Loss      uint16
	Draws     uint16
}

type flopHand struct {
	Hand      string `gorm:"primaryKey"`
	HoleCards string `gorm:"index"`
	Board     string `gorm:"index"`
	Wins      uint16
	Loss      uint16
	Draws     uint16
}

type hand struct {
	HoleCards      []string
	CommunityCards []string
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

func generateHandCombinations(countCommunityCards int) []hand {
	result := make([]hand, 0)
	var holeCards = generateAbstractedHoleCards()

	for _, holeCard := range holeCards {
		boards := generateAbstractedBoards(holeCard, countCommunityCards)
		for _, board := range boards {
			result = append(result, hand{HoleCards: board[:2], CommunityCards: board[2:]})
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

func worker(jobs <-chan hand, results chan<- flopHand, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for job := range jobs {
		strength := calculateHandStrength(job.HoleCards, job.CommunityCards)
		results <- flopHand{
			Hand: strings.Join(job.HoleCards, "") + strings.Join(
				job.CommunityCards,
				"",
			), HoleCards: strings.Join(job.HoleCards, ""),
			Board: strings.Join(
				job.HoleCards,
				"",
			),
			Wins: strength.Wins,
			Loss: strength.Loss, Draws: strength.Draws,
		}
		bar.Add(1)
	}
}

func main() {
	fmt.Println(runtime.NumCPU())
	dsn := "host=localhost user=postgres password=example  port=5555"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err := db.AutoMigrate(&flopHand{})
	if err != nil {
		return
	}
	hands := generateHandCombinations(3)

	bar := progressbar.Default(int64(len(hands)))

	const numWorkers = 8

	jobs := make(chan hand, len(hands))
	results := make(chan flopHand, len(hands))
	var wg sync.WaitGroup

	// Worker starten
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, bar)
	}

	// Jobs in den Kanal senden
	for _, item := range hands {
		jobs <- item

	}
	close(jobs) // Keine weiteren Jobs mehr

	wg.Wait()

	toSave := make([]flopHand, 0)
	for result := range results {
		toSave = append(toSave, result)
	}

	db.CreateInBatches(toSave, 100)
	file, _ := os.Create("flopAbstraction.gob")

	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(toSave)

}
