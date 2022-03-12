package landgrab

import (
	"fmt"
	"math"
	"time"

	"github.com/montanaflynn/stats"
)

type attackerSummary struct {
	attackers int
	p10       float64
	p50       float64
	p90       float64
	prob      float64
	trials    int
}

type attackerTrial struct {
	defendingTerritories []int
	attackers            int
	margin               float64
}

func (a attackerSummary) String() string {
	probStr := "100  "
	if a.prob < 99.999 {
		probStr = fmt.Sprintf("%05.2f", a.prob)
	}

	return fmt.Sprintf("%-7d%s    %-3d   %-3d   %-3d   %-3d", a.attackers, probStr, int(a.p10), int(a.p50), int(a.p90), a.trials)
}

func (a attackerTrial) run() attackerSummary {
	successes := 0.0

	var results stats.Float64Data
	trials := 0
	minTrials := 100
	zBar := 1.96 // 95% confidence

	for ; ; trials++ {
		remaining := float64(campaign(a.attackers, a.defendingTerritories))
		if remaining != 0 {
			successes += 1
		}

		results = append(results, remaining)
		if trials >= minTrials {
			stddev, _ := results.StandardDeviation()
			if zBar*stddev/math.Sqrt(float64(trials)) < a.margin {
				break
			}
		}
	}
	p10, _ := results.Percentile(10)
	p50, _ := results.Percentile(50)
	p90, _ := results.Percentile(90)
	floatTrials := float64(trials)
	return attackerSummary{
		attackers: a.attackers,
		p10:       p10,
		p50:       p50,
		p90:       p90,
		prob:      (successes / floatTrials) * 100,
		trials:    trials,
	}
}

func broker(trials <-chan attackerTrial, results chan<- attackerSummary) {
	for {
		trial := <-trials
		result := trial.run()
		results <- result
	}
}

func DetermineAttackers(defendingTerritories []int) {
	//attackers := 1
	margin := .1
	totalDefendingArmies := 0
	for i := 0; i < len(defendingTerritories); i++ {
		totalDefendingArmies += defendingTerritories[i]
	}

	fmt.Printf("Calculating size of force needed to defeat %d armies to claim %d territories\n", totalDefendingArmies, len(defendingTerritories))
	fmt.Println("Attack Success  p10   p50   p90 trials")
	// Start attackers at 2 since that's how many you need to attack
	start := time.Now()
	totalTrials := 0
	trials := make(chan attackerTrial)
	results := make(chan attackerSummary)
	go broker(trials, results)
	finished := false
	for attackers := 2; !finished; attackers++ {
		select {

		case summary := <-results:

			totalTrials += summary.trials
			fmt.Printf("%s\n", summary)
			if summary.prob > 99.99 {
				finished = true
			}
		default:
			trial := attackerTrial{attackers: attackers, defendingTerritories: defendingTerritories, margin: margin}
			trials <- trial
		}
	}
	duration := time.Since(start)
	fmt.Printf("Finished %d trials in %v\n", totalTrials, duration)
}
