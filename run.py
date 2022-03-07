#!/usr/bin/env python3

import random

DICE_SIDES = 6

def roll(num_dice):
    ret = []
    for i in range(num_dice):
        ret.append(random.randint(0, DICE_SIDES) + 1)
    return sorted(ret, reverse=True)



def one_round(attackers, defenders, dice):
    if dice > attackers:
        dice = attackers
    attacker_dice = min(dice, attackers)
    defender_dice = min(2, defenders)

    matchups = min(attacker_dice, defender_dice)

    attacker_role = roll(attacker_dice)
    defender_role = roll(defender_dice)

    for i in range(matchups):
        # Attacker's role must be larger than defenders,
        # tie goes to defender
        if attacker_role[i] > defender_role[i]:
            defenders-=1
        else:
            attackers-=1
    return attackers, defenders



def main():
    for i in range(100):
        attackers = 5
        defenders = 1
        dice = 3
        after_attackers, after_defenders = one_round(attackers, defenders, dice)
        print(f"{attackers} attacked {defenders} with {dice}, resulting in {after_attackers} to {after_defenders}")

if __name__ == "__main__":
    main()
