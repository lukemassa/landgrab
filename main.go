package main

import (
	"os"
	"strconv"

	"github.com/lukemassa/landgrab/pkg/landgrab"
)

func main() {
	if len(os.Args) < 2 {
		panic("Must enter at least one territory")
	}
	attackers := make([]int, len(os.Args)-1)
	for i := 0; i < len(os.Args)-1; i++ {
		attacker, err := strconv.Atoi(os.Args[i+1])
		if err != nil {
			panic(err)
		}
		attackers[i] = attacker
	}
	landgrab.DetermineAttackers(attackers)
}
