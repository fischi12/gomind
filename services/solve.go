package services

import "fmt"

func runPreflopIteration(history PreFlopHistory) {
	cfr := vanillaCfr(history, 0, 1, 1)
	vanillaCfr(history, 1, 1, 1)
	strategy := InfoSetsPreFlop["14"].Strategy["f"]
	regret := InfoSetsPreFlop["14"].Regret["f"]
	cumulativeStrategy := InfoSetsPreFlop["14"].CumulativeStrategy["f"]
	fmt.Println(cfr)
	fmt.Println(strategy)
	fmt.Println(regret)
	fmt.Println(cumulativeStrategy)
	//Double strategyFold = infoSets.get("14").getStrategy().get(Action.FOLD);
	//Double regretFold = infoSets.get("14").getRegret().get(Action.FOLD);
	//Double cumulativeStrategyFold = infoSets.get("14").getCumulativeStrategy().get(Action.FOLD);
	//assertEquals(0.9912767045168166, strategyFold, 0.00001);
	//assertEquals(49.52140287769784, regretFold, 0.00001);
	//assertEquals(0.2, cumulativeStrategyFold, 0.00001);
	//
	//assertEquals(22.093750000000004, acctualPlayer, 0.00001);
	//assertEquals(-51.02140287769784, acctualOpponent, .00001);
}
