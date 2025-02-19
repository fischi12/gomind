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
	var contacts []models.FlopHand

	db.Find(&contacts)
	return contacts
}
