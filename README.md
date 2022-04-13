# Landgrab

Script to help determine how many pieces to place in order to take over a number of territories in https://landgrab.net/

## Usage

The code has been tested on windows, MacOS, and CentOS. It doesn't use any non-portable features of go so should theoretically work anywhere go does.

1. Make sure go is installed (https://go.dev)
1. `git clone git@github.com:lukemassa/landgrab.git`
2. `cd landgrab`
3. `go run main.go -h` and follow the instructions.

## Understanding the output

```
lukemassa@Lukes-MacBook-Pro landgrab % ./landgrab 3 2 1
Calculating size of force needed to defeat 6 armies to claim 3 territories
Shows the percent chance of overtaking all territories, as well as percentile of how many attacking armies are expected to be left

Attack Success  p10   p50   p90   trials
2      00.00    0     0     0     100
3      00.00    0     0     0     100
4      03.00    0     0     0     100
5      14.00    0     0     1     100
6      27.64    0     0     2     275
7      41.98    0     0     3     667
8      55.62    0     1     4     1077
9      68.88    0     2     5     1536
10     79.60    0     3     6     1873
11     84.56    0     4     7     2377
12     90.20    1     5     8     2611
13     92.72    1     6     9     2938
14     95.68    3     7     10    3102
15     97.41    4     8     11    3171
16     97.82    4     9     12    3584
17     98.80    5     10    13    3496
18     99.31    6     11    14    3496
19     99.73    8     12    15    3283
20     99.89    8     13    16    3488
21     99.81    9     14    17    3644
22     99.94    10    15    18    3608
23     99.94    12    16    19    3610
24     100      12    17    20    3608

Finished 51844 trials in 55.364767ms
```

This is telling you that, say, if you start attacking three territories in a row with 14 armies, you'd have a 95.68% chance of getting to the end, that there's a 10% chance you'd have 3 armies remaining, 50% chance you'd have 7, and 90% chance you'd have 10 (technically these are percentiles), having run 3117 test trials (this value included mostly for debugging).

The assumption is you attack in a single straight line, leave exactly one army in every conquered territory, and always attack with 3 if possible.

## How does the code work?

`pkg/landgrab/attacking.go` has all the logic for how the game itself actually works (dice rolls, etc.)

`pkg/landgrab/trials.go` runs a number of "trials" per attacker and then does some averages/etc on the results. The number of trials it runs per attacker is a function of how long it takes to get "confident" (which is why the trials go up as you increase attackers, because there's more variability), but it always runs at least 100 trials.

The code runs in parallel as fast as it possibly can, so if you enter a large number like `./landgrab 100` you might see a noticeable slowdown in your computer/fans start to turn on. Simply run `ctrl+c` (or equivalent) to kill it.

It doesn't have any temp files or anything to cleanup, so always safe to stop it midway.
