package landgrab

// Logic for how "attacking" is implemented in landgrab. It's similar to the board game Risk

import (
	"math/rand"
	"slices"
)

const diceSides = 6
const maxDefenders = 2

// always attack with 3 if possible
const preferredAttackers = 3

// Roll a number of dice, and return them in reverse order
func roll(numDice int, r *rand.Rand) []int {
	ret := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		ret[i] = r.Intn(diceSides) + 1
	}
	slices.Sort(ret)
	slices.Reverse(ret)
	return ret
}

// Simulate a single click of the "attack" button
func oneRound(attackers, defenders, dice int, r *rand.Rand) (int, int) {
	// Can only attack with as many attackers as are present
	attackerDice := min(dice, attackers)
	// Defender defends with as many as it can, up to maxDefenders
	defenderDice := min(maxDefenders, defenders)

	matchups := min(attackerDice, defenderDice)

	attackerRoll := roll(attackerDice, r)
	defenderRoll := roll(defenderDice, r)

	// Walk through each matchup, for every one "lost", that player decreases an army
	// Note tie goes to the defender
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

	// While attackers have more than one (can still attack) and defenders more than 0
	// (still own the territory) continue to attack
	for attackers > 1 && defenders > 0 {
		attackers, defenders = oneRound(attackers, defenders, preferredAttackers, r)
	}
	// Defender lost, leave one "behind" and move the rest
	// Possible TODO: Play around with what happens if you leave behind more than one
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
