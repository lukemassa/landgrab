package landgrab

import (
	"fmt"
	"math"
	"time"

	"github.com/montanaflynn/stats"
)

const minAttackers int = 2

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
	//fmt.Printf("Working on attacker %d\n", a.attackers)

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

func oneBroker(trials <-chan attackerTrial, results chan<- attackerSummary) {
	for {
		trial := <-trials
		result := trial.run()
		results <- result
	}
}

func reorderResults(resultBufferChan <-chan attackerSummary, results chan<- attackerSummary) {

	// Pretend like you saw one before the min so we can start sending right away
	lastSent := minAttackers - 1
	resultBuffer := make([]attackerSummary, 0)
	for {
		result := <-resultBufferChan
		//fmt.Printf("Got result for %d\n", result.attackers)
		resultBuffer = append(resultBuffer, result)

		// Loop trying to send until there is nothing in the buffer to send
		for {
			sent := false
			// Walk through the buffer backwards, if a result is the
			// one to send, send it, and clip it from the end of the buffer
			for i := len(resultBuffer) - 1; i >= 0; i-- {
				if resultBuffer[i].attackers == lastSent+1 {

					results <- resultBuffer[i]
					lastSent = resultBuffer[i].attackers
					resultBuffer = append(resultBuffer[:i], resultBuffer[i+1:]...)
					sent = true
				}
			}
			if !sent {
				break
			}
		}
	}
}

func broker(trials <-chan attackerTrial, results chan<- attackerSummary) {
	numWorkers := 10
	resultBufferChan := make(chan attackerSummary, numWorkers)
	go reorderResults(resultBufferChan, results)
	for i := 0; i < numWorkers; i++ {
		oneBrokerAttackerChan := make(chan attackerTrial)
		go func() {

			go oneBroker(oneBrokerAttackerChan, resultBufferChan)
			for {
				trial := <-trials
				oneBrokerAttackerChan <- trial
			}
		}()
	}
}

func DetermineAttackers(defendingTerritories []int) {
	margin := .1
	totalDefendingArmies := 0
	for i := 0; i < len(defendingTerritories); i++ {
		totalDefendingArmies += defendingTerritories[i]
	}

	fmt.Printf("Calculating size of force needed to defeat %d armies to claim %d territories\n", totalDefendingArmies, len(defendingTerritories))
	fmt.Println("Attack Success  p10   p50   p90   trials")
	// Start attackers at 2 since that's how many you need to attack
	start := time.Now()
	totalTrials := 0
	trials := make(chan attackerTrial)
	results := make(chan attackerSummary)
	broker(trials, results)
	finished := false
	for attackers := minAttackers; !finished; {
		select {

		case summary := <-results:

			totalTrials += summary.trials
			fmt.Printf("%s\n", summary)
			if summary.prob > 99.99 {
				finished = true
			}
		case trials <- attackerTrial{attackers: attackers, defendingTerritories: defendingTerritories, margin: margin}:
			attackers++
		}
	}
	duration := time.Since(start)
	fmt.Printf("Finished %d trials in %v\n", totalTrials, duration)
}
