package main

import (
	"fmt"
	"os"
	"strconv"
    "strings"

	"github.com/lukemassa/landgrab/pkg/landgrab"
)

func usage() {
	fmt.Println("Usage: ./landgrab N [N...]")
	fmt.Println("Example: './landgrab 3 2 1' would tell you how many armies to place to win a campaign against three territories, with 3, 2, 1 armies in them")
	os.Exit(1)
}

func main() {
    // If no args, or arg begins with '-', just print the usage
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		usage()
	}
    // Convert args into slice of ints, panic if it can't parse
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
