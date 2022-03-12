package landgrab

import (
	"fmt"

	"github.com/montanaflynn/stats"
)

type attackerSummary struct {
	attackers int
	p10       float64
	p50       float64
	p90       float64
	prob      float64
}

func (a attackerSummary) String() string {
	probStr := "100  "
	if a.prob < 99.999 {
		probStr = fmt.Sprintf("%05.2f", a.prob)
	}

	return fmt.Sprintf("%-7d%s    %-3d   %-3d   %-3d", a.attackers, probStr, int(a.p10), int(a.p50), int(a.p90))
}

func getAttackerSummary(attackers int, defendingTerritories []int, trials int) attackerSummary {
	successes := 0.0

	results := make([]int, trials)

	for i := 0; i < trials; i++ {
		//fmt.Println(i)
		remaining := campaign(attackers, defendingTerritories)
		if remaining != 0 {
			successes += 1
		}

		results[i] = remaining
		//	fmt.Println(i)
	}
	resultsForStats := stats.LoadRawData(results)
	p10, _ := resultsForStats.Percentile(10)
	p50, _ := resultsForStats.Percentile(50)
	p90, _ := resultsForStats.Percentile(90)
	floatTrials := float64(trials)
	return attackerSummary{
		attackers: attackers,
		p10:       p10,
		p50:       p50,
		p90:       p90,
		prob:      (successes / floatTrials) * 100,
	}
}

func DetermineAttackers(defendingTerritories []int) {
	//attackers := 1
	trialsPerAttacker := 10_000
	totalDefendingArmies := 0
	for i := 0; i < len(defendingTerritories); i++ {
		totalDefendingArmies += defendingTerritories[i]
	}

	fmt.Printf("Calculating size of force needed to defeat %d armies to claim %d territories\n", totalDefendingArmies, len(defendingTerritories))
	fmt.Println("Attack Success  p10   p50   p90")
	for attackers := 0; ; attackers++ {
		summary := getAttackerSummary(attackers, defendingTerritories, trialsPerAttacker)
		fmt.Printf("%s\n", summary)
		if summary.prob > 99.99 {
			return
		}
	}
}
