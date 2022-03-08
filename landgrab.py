#!/usr/bin/env python3

import random

DICE_SIDES = 6

def roll(num_dice):
    ret = []
    for i in range(num_dice):
        ret.append(random.randint(0, DICE_SIDES) + 1)
    return sorted(ret, reverse=True)



def one_round(attackers, defenders, dice):
    """
    One round of attacks, returning new army levels
    """
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


def invade(attackers, defenders):
    """
    Attack with as much as you have until it's either taken over or you can't attack
    """
    while attackers >1 and defenders > 0:
        attackers, defenders = one_round(attackers, defenders, 3)
    # Decision here to march all troops in after winning.
    # If we decide to only march some, you have to keep track of
    # how many attacking dice were rolled, because at least that many
    # have to march
    if defenders == 0:
        # Defender lost, move all troops
        return 1, attackers-1, True

    return attackers, defenders, False

def campaign(attackers, defending_territories):
    """
    Given an attacking territory w attackers, invade one territory after another
    Return how many attacking armies are in the last territory, or 0 if you never get there
    """
    for defending_territory in defending_territories:
        # After invasion, we will leave behind everyone, and whoever is in the previously
        # defending territory will become next attackers. If invasion failed, we quit now
        _, attackers, success = invade(attackers, defending_territory)
        if not success:
            return 0
    return attackers


def probability_desired_remaining(attackers, defending_territories, desired_remaining_attackers, trials):
    # TODO: Make the code figure out num of correct trials based on some tolerance
    success = 0.0
    total = 0.0
    for i in range(trials):
        remaining = campaign(attackers, defending_territories)
        if remaining >= desired_remaining_attackers:
            success+=1.0
        total += remaining
    return success/trials, total/trials



def determine_attackers(defending_territories, desired_remaining_attackers):
    attackers = 1
    trials_per_attacker = 10_000
    while True:
        prob, avg = probability_desired_remaining(attackers, defending_territories, desired_remaining_attackers, trials_per_attacker)
        attackers+=1
        print("%-5s%.4f %.2f" % (attackers, prob, avg))
        if prob > .99:
            break


def main():
    determine_attackers([4, 5, 1, 1, 1, 1], 10)


if __name__ == "__main__":
    main()
