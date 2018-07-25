# Dito Programming Language

Dito is a work-in-progress toy interpeted language made to learn about writing parsers, interpeters etc. Things change around a bit so some docs arn't up to date. The syntax is pretty generic at the moment but probally will try out some more idiosyncratic concepts in the future.

This is the most up to date repo but its not that stable.

### Examples

For full list of examples see [`examples`](examples/) directory or the [stdlib](lib/std.dito).
Heres example of dynamic programming used to find the total number of ways a target number can
be made given a set of coins. This can be found in full in the exampls dir described.

```swift
import std

def coinSums(coins, target) {
    let sack = zeros(target+1)
    sack[0] = 1
    for coin in coins {
        for i in range(coin, target+1) {
            sack[i] += sack[i-coin]
        }
    }
    return sack[target]
}

let zeros = def(t) -> map(def(_) -> 0, range(0, t))

let coins  = [1, 2, 5, 10, 20, 50, 100, 200]
let target = 200

print(coinSums(coins, target))
```


### Having a go...

The most simple way to build with `go build` then use the executable or just `go run main.go`. You can try out the language either in a interactive shell with:

    ./dito

or to run a file:

    ./dito funcs.dito

