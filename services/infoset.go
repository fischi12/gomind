package services

type InfoSet struct {
	InfoSetKey         string
	Actions            []string
	Player             int
	Strategy           map[string]float64
	Regret             map[string]float64
	CumulativeStrategy map[string]float64
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
