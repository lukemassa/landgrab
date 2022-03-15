# landgrab

Script to help determine how many pieces to place in order to take over a number of territories in https://landgrab.net/

## Usage

### Mac/Linux

1. `git clone git@github.com:lukemassa/landgrab.git`
2. `cd landgrab`
3. `./landgrab -h` and follow the instructions.

### Windows

1. Make sure go is installed (https://go.dev/)
1. `git clone git@github.com:lukemassa/landgrab.git`
1. `dir landgrab`
1. `go run main.go -h` and follow the instructions.

## Understanding the output

```
lukemassa@Lukes-MacBook-Pro landgrab % ./landgrab 3 2 1          
Calculating size of force needed to defeat 6 armies to claim 3 territories
Attack Success  p10   p50   p90   trials
2      00.00    0     0     0     100
3      00.00    0     0     0     100
4      06.00    0     0     0     100
5      13.79    0     0     1     116
6      32.34    0     0     2     303
7      44.41    0     0     3     653
8      57.69    0     1     4     1054
9      68.21    0     2     5     1532
10     77.39    0     3     6     1942
11     84.76    0     4     7     2271
12     88.77    0     5     8     2805
13     93.23    2     6     9     2896
14     94.71    2     7     10    3117
15     97.29    4     8     11    3206
16     97.98    5     9     12    3410
17     98.38    5     10    13    3705
18     99.46    6     11    14    3538
19     99.54    8     12    15    3505
20     99.81    9     13    16    3679
21     99.94    10    14    17    3458
22     99.89    11    15    18    3560
23     100      12    16    19    3391
Finished 48441 trials in 315.358621ms
```

This is telling you that, say, if you start attacking three territories in a row with 14 armies, you'd have a 94.71% chance of getting to the end, that there's a 10% chance you'd have 2 armies remaining, 50% chance you'd have 7, and 90% chance you'd have 13 (technically these are percentiles), having run 3117 trials (mostly for debugging).

The assumption is you attack in a single straight line, leave exactly one army in every conquered territory, and always attack with 3 if possible.

## How does the code work?

`pkg/landgrab/attacking.go` has all the logic for how the game itself actually works (dice rolls, etc.)

`pkg/landgrab/trials.go` runs a number of "trials" per attacker and then does some averages/etc on the results. The number of trials it runs per attacker is a function of how long it takes to get "confident" (which is why the trials go up as you increase attackers, because there's more variability), but it always runs at least 100 trials.

The code runs in parallel as fast as it possibly can, so if you enter a large number like `./landgrab 100` you'll probably see a noticeable slowdown in your computer. Simply run `ctrl+c` to kill it.

It doesn't have any temp files or anything to cleanup, so always safe to stop it.
