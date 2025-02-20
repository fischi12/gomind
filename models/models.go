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

//type InfoSetPreFlop struct {
//	InfoSetKey         string             `gorm:"primaryKey"`
//	Actions            []string           `gorm:"column:actions;type:jsonb"`
//	Player             int                `gorm:"column:player"`
//	Strategy           map[string]float64 `gorm:"column:strategy;type:jsonb"`
//	Regret             map[string]float64 `gorm:"column:regret;type:jsonb"`
//	CumulativeStrategy map[string]float64 `gorm:"column:cumulative_strategy;type:jsonb"`
//}
//
//type InfoSetPostFlop struct {
//	InfoSetKey         string             `gorm:"primaryKey"`
//	Actions            []string           `gorm:"column:actions;type:jsonb"`
//	Player             int                `gorm:"column:player"`
//	Strategy           map[string]float64 `gorm:"column:strategy;type:jsonb"`
//	Regret             map[string]float64 `gorm:"column:regret;type:jsonb"`
//	CumulativeStrategy map[string]float64 `gorm:"column:cumulative_strategy;type:jsonb"`
//}
