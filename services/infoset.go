package services

type InfoSet struct {
	InfoSetKey         string             `gorm:"primaryKey"`
	Actions            []string           `gorm:"column:actions;type:jsonb"`
	Player             int                `gorm:"column:player"`
	Strategy           map[string]float64 `gorm:"column:strategy;type:jsonb"`
	Regret             map[string]float64 `gorm:"column:regret;type:jsonb"`
	CumulativeStrategy map[string]float64 `gorm:"column:cumulative_strategy;type:jsonb"`
}

func (infoSet *InfoSet) GetStrategy() {
	regret := map[string]float64{}
	regretSum := 0.0
	for key, value := range infoSet.Regret {
		maxValue := max(value, 0)
		regret[key] = maxValue
		regretSum += maxValue
	}
	if regretSum > 0 {
		for key, value := range regret {
			infoSet.Strategy[key] = value / regretSum
		}
	} else {
		for _, action := range infoSet.Actions {
			infoSet.Strategy[action] = 1 / float64(len(infoSet.Actions))
		}
	}
}
