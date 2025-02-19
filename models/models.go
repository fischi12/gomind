package models

type TurnHand struct {
	Hand      string `gorm:"primaryKey"`
	HoleCards string `gorm:"index"`
	Board     string `gorm:"index"`
	Wins      uint16
	Loss      uint16
	Draws     uint16
}

type FlopHand struct {
	Hand      string `gorm:"primaryKey"`
	HoleCards string `gorm:"index"`
	Board     string `gorm:"index"`
	Wins      uint16
	Loss      uint16
	Draws     uint16
}
