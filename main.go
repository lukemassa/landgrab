package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lukemassa/landgrab/pkg/landgrab"
)

func usage() {
	fmt.Println("Usage: ./landgrab N [N...]")
	fmt.Println("Example: './landgrab 3 2 1' would tell you how many armies to place to win a campaign against three territories, with 3, 2, 1 armies in them")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
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
