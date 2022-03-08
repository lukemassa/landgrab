package landgrab

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/montanaflynn/stats"
)

const diceSides = 6

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func roll(numDice int) []int {
	ret := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		ret[i] = rand.Intn(diceSides) + 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ret)))
	return ret
}

func oneRound(attackers, defenders, dice int) (int, int) {
	if dice > attackers {
		dice = attackers
	}
	attackerDice := min(dice, attackers)
	defenderDice := min(2, defenders)

	matchups := min(attackerDice, defenderDice)

	attackerRoll := roll(attackerDice)
	defenderRoll := roll(defenderDice)

	for i := 0; i < matchups; i++ {
		if attackerRoll[i] > defenderRoll[i] {
			defenders--
		} else {
			attackers--
		}
	}
	return attackers, defenders
}

// Attack with as much as you ahve until it's either taken or you can't attack
func invade(attackers, defenders int) (int, int, bool) {

	for attackers > 1 && defenders > 0 {
		attackers, defenders = oneRound(attackers, defenders, 3)
	}
	if defenders == 0 {
		return 1, attackers - 1, true
	}
	return attackers, defenders, false
}

// Given an attacking territory w attackers, invade one territory after another
// Return how many attacking armies are in the last territory, or 0 if you never get there
func campaign(attackers int, defenderingTerritories []int) int {
	for i := 0; i < len(defenderingTerritories); i++ {
		_, newAttackers, success := invade(attackers, defenderingTerritories[i])
		//fmt.Printf("Invading %d->%d leaves %d in defending w success %v\n", attackers, defenderingTerritories[i], newAttackers, success)
		attackers = newAttackers
		if !success {
			return 0
		}
	}
	return attackers
}

func probabilityDesiredRemaining(attackers int, defendingTerritories []int, trials int) (float64, float64, float64, float64) {
	successes := 0.0
	const batches = 1

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
	return p10, p50, p90, (successes / floatTrials) * 100
}

type attackerSummary struct {
	p10  float64
	p50  float64
	p90  float64
	prob float64
}

func (a attackerSummary) String() string {
	probStr := "100  "
	if a.prob < 99.999 {
		probStr = fmt.Sprintf("%05.2f", a.prob)
	}

	return fmt.Sprintf("%s    %-3d   %-3d   %-3d", probStr, int(a.p10), int(a.p50), int(a.p90))
}

func getAttackerSummary(attackers int, defendingTerritories []int, trialsPerAttacker int) attackerSummary {
	p10, p50, p90, prob := probabilityDesiredRemaining(attackers, defendingTerritories, trialsPerAttacker)
	return attackerSummary{
		p10:  p10,
		p50:  p50,
		p90:  p90,
		prob: prob,
	}
}

func DetermineAttackers(defendingTerritories []int) {
	attackers := 1
	trialsPerAttacker := 10_000
	totalDefendingArmies := 0
	for i := 0; i < len(defendingTerritories); i++ {
		totalDefendingArmies += defendingTerritories[i]
	}
	fmt.Printf("Calculating size of force needed to defeat %d armies to claim %d territories\n", totalDefendingArmies, len(defendingTerritories))
	fmt.Println("Attack Success  p10   p50   p90")
	for {
		summary := getAttackerSummary(attackers, defendingTerritories, trialsPerAttacker)
		fmt.Printf("%-7d%s\n", attackers, summary)
		if summary.prob > 99.99 {
			break
		}
		attackers++
	}
}
