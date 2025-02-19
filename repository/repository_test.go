package repository

import (
	"gomind/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"reflect"
	"testing"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der DB: %v", err)
	}

	db.AutoMigrate(&models.FlopHand{})

	return db
}

func Test_createFlopHand(t *testing.T) {
	db := SetupTestDB()
	hand := []models.FlopHand{{Hand: "Hand", HoleCards: "HoleCards", Board: "Board", Wins: 1, Draws: 2, Loss: 3}}
	err := UpsertFlopHand(db, &hand)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var expected []models.FlopHand
	db.Find(&expected)

	if !reflect.DeepEqual(hand, expected) {
		t.Errorf("%v, want %v", expected, hand)
	}
}

func Test_upsertFlopHand(t *testing.T) {
	db := SetupTestDB()
	hand := []models.FlopHand{{Hand: "Hand", HoleCards: "HoleCards", Board: "Board", Wins: 1, Draws: 2, Loss: 3}}

	duplicate := []models.FlopHand{
		{
			Hand: "Hand", HoleCards: "duplicate", Board: "duplicate", Wins: 2, Draws: 3,
			Loss: 4,
		},
	}
	err := UpsertFlopHand(db, &hand)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = UpsertFlopHand(db, &duplicate)
	if err != nil {
		t.Fatalf("error duplicate: %v", err)
	}
	var expected []models.FlopHand
	db.Find(&expected)

	if !reflect.DeepEqual(duplicate, expected) {
		t.Errorf("%v, want %v", expected, duplicate)
	}
}

func Test_findAllFlopHand(t *testing.T) {
	db := SetupTestDB()
	hand := []models.FlopHand{{Hand: "Hand", HoleCards: "HoleCards", Board: "Board", Wins: 1, Draws: 2, Loss: 3}}

	handSecond := []models.FlopHand{
		{
			Hand: "HandSecond", HoleCards: "duplicate", Board: "duplicate", Wins: 2, Draws: 3,
			Loss: 4,
		},
	}
	_ = UpsertFlopHand(db, &hand)
	_ = UpsertFlopHand(db, &handSecond)

	actual := FindAllFlopHand(db)

	expected := []models.FlopHand{hand[0], handSecond[0]}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v, want %v", actual, expected)
	}
}
