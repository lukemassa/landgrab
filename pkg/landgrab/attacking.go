package landgrab

import (
	"math/rand"
	"sort"
)

const diceSides = 6

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func roll(numDice int, r *rand.Rand) []int {
	ret := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		ret[i] = r.Intn(diceSides) + 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ret)))
	return ret
}

func oneRound(attackers, defenders, dice int, r *rand.Rand) (int, int) {
	if dice > attackers {
		dice = attackers
	}
	attackerDice := min(dice, attackers)
	defenderDice := min(2, defenders)

	matchups := min(attackerDice, defenderDice)

	attackerRoll := roll(attackerDice, r)
	defenderRoll := roll(defenderDice, r)

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
func invade(attackers, defenders int, r *rand.Rand) (int, int, bool) {

	for attackers > 1 && defenders > 0 {
		attackers, defenders = oneRound(attackers, defenders, 3, r)
	}
	if defenders == 0 {
		return 1, attackers - 1, true
	}
	return attackers, defenders, false
}

// Given an attacking territory w attackers, invade one territory after another
// Return how many attacking armies are in the last territory, or 0 if you never get there
func campaign(attackers int, defenderingTerritories []int, r *rand.Rand) int {
	for i := 0; i < len(defenderingTerritories); i++ {
		_, newAttackers, success := invade(attackers, defenderingTerritories[i], r)
		//fmt.Printf("Invading %d->%d leaves %d in defending w success %v\n", attackers, defenderingTerritories[i], newAttackers, success)
		attackers = newAttackers
		if !success {
			return 0
		}
	}
	return attackers
}
