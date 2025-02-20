package services

type history interface {
	isTerminal() bool
	terminalUtility(player int) float64
	isChance() bool
	sampleChanceOutcome() string
	getInfoSet() InfoSet
	player() int
	newHistory([]string) history
	getHistory() []string
}

func vanillaCfr(history history, player int, pi_0 float64, pi_1 float64) float64 {
	if history.isTerminal() {
		return history.terminalUtility(player)
	} else if history.isChance() {
		a := history.sampleChanceOutcome()
		return vanillaCfr(history.newHistory(append(history.getHistory(), a)), player, pi_0, pi_1)
	}
	infoSet := history.getInfoSet()

	v := 0.0
	va := map[string]float64{}

	for _, a := range infoSet.Actions {
		if history.player() == 0 {
			va[a] = vanillaCfr(
				history.newHistory(append(history.getHistory(), a)),
				player,
				infoSet.Strategy[a]*pi_0,
				pi_1,
			)
		} else {
			va[a] = vanillaCfr(
				history.newHistory(append(history.getHistory(), a)),
				player,
				pi_0,
				infoSet.Strategy[a]*pi_1,
			)
		}
		v += infoSet.Strategy[a] * va[a]
	}
	if history.player() == player {
		for _, a := range infoSet.Actions {
			if player == 0 {
				infoSet.Regret[a] += pi_1 * (va[a] - v)
				infoSet.CumulativeStrategy[a] += pi_0 * infoSet.Strategy[a]
			} else {
				infoSet.Regret[a] += pi_0 * (va[a] - v)
				infoSet.CumulativeStrategy[a] += pi_1 * infoSet.Strategy[a]
			}

		}
		infoSet.GetStrategy()
	}
	return v

}
