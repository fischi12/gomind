package repository

import (
	"gomind/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpsertFlopHand(db *gorm.DB, flopHand *[]models.FlopHand) error {
	return db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		},
	).CreateInBatches(flopHand, 100).Error
}

func FindAllFlopHand(db *gorm.DB) []models.FlopHand {
	var hands []models.FlopHand

	db.Find(&hands)
	return hands
}

func UpsertTurnHand(db *gorm.DB, turnHand *[]models.TurnHand) error {
	return db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		},
	).CreateInBatches(turnHand, 100).Error
}

func FindAllTurnHand(db *gorm.DB) []models.TurnHand {
	var hands []models.TurnHand

	db.Find(&hands)
	return hands
}

//
//func UpsertInfoSetPreFlop(db *gorm.DB, infoSet *[]models.InfoSetPreFlop) error {
//	return db.Clauses(
//		clause.OnConflict{
//			UpdateAll: true,
//		},
//	).CreateInBatches(infoSet, 100).Error
//}
//
//func UpsertInfoSetPostFlop(db *gorm.DB, infoSet *[]models.InfoSetPostFlop) error {
//	return db.Clauses(
//		clause.OnConflict{
//			UpdateAll: true,
//		},
//	).CreateInBatches(infoSet, 100).Error
//}
