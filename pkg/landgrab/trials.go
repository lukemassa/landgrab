package landgrab

import (
	"fmt"
	"math"
	"math/rand"
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

// Runs one trial; for a specific number of attackers, return statistics about
// how successful a campaign would be
func (a attackerTrial) run(r *rand.Rand) attackerSummary {
	successes := 0.0
	//fmt.Printf("Working on attacker %d\n", a.attackers)

	var results stats.Float64Data
	trials := 0
	minTrials := 100
	zBar := 1.96 // 95% confidence

	for ; ; trials++ {
		remaining := float64(campaign(a.attackers, a.defendingTerritories, r))
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

// oneBroker runs sequentially and simply grabs a trial from a channel
// calculates its results and puts it in the new channel
func oneBroker(trials <-chan attackerTrial, results chan<- attackerSummary) {

	// We create a source here instead of relying on the default source
	// since the default source is thread-safe, and thus must spend a lot of time
	// locking and unlocking.
	// A newSource is not threadsafe (https://pkg.go.dev/math/rand), which is why it is created here
	// since this method is run sequentually
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for {
		trial := <-trials
		result := trial.run(r)
		results <- result
	}
}

// reorderResults takes results from a buffer channel and puts them
// into a results channel
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

// Coordinates the brokers that do the trials, reorders the results, and returns them
func brokers(trials <-chan attackerTrial, results chan<- attackerSummary) {
	numWorkers := 10
	resultBufferChan := make(chan attackerSummary, numWorkers)
	go reorderResults(resultBufferChan, results)
	for i := 0; i < numWorkers; i++ {
		oneBrokerAttackerChan := make(chan attackerTrial)
		go func() {
			// All this anonymous function is doing is firing off the broker then
			// reodering the results by sending them to the reordering channel
			go oneBroker(oneBrokerAttackerChan, resultBufferChan)
			for {

				oneBrokerAttackerChan <- (<-trials)
			}
		}()
	}
}

// Given a list of defending territories, prints to stdout a summary of how many attackers
// should be placed to overtake them
func DetermineAttackers(defendingTerritories []int) {
	margin := .1
	totalDefendingArmies := 0
	for i := 0; i < len(defendingTerritories); i++ {
		totalDefendingArmies += defendingTerritories[i]
	}

	fmt.Printf("Calculating size of force needed to defeat %d armies to claim %d territories\n", totalDefendingArmies, len(defendingTerritories))
	fmt.Println("Shows the percent chance of overtaking all territories, as well as percentile of how many attacking armies are expected to be left")
	fmt.Println()
	fmt.Println("Attack Success  p10   p50   p90   trials")
	// Start attackers at 2 since that's how many you need to attack
	start := time.Now()
	totalTrials := 0
	trials := make(chan attackerTrial)
	results := make(chan attackerSummary)
	brokers(trials, results)
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
	fmt.Println()
	fmt.Printf("Finished %d trials in %v\n", totalTrials, duration)
}
