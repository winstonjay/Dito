# Dito Programming Language

Dito is a


### Examples:

`greet` lambda function, as you might suspect, it outputs 'Hello, Dito!'
```go
greet = func(name) -> ("Hello, " + name)
print(greet("Dito!"))
```

```go
func newton_sqrt(x, delta=1/128) {
    z = delta
    for abs(x - z**2) > delta {
        z -= (z**2 - x) / (2 * x)
    }
    return z
}
```

```go
func coin_sums(coins, target) {
    sack = [0] * len(coins)
    sack[0] = 1
    for coin in coins {
        for i in iota(coin, target) {
            sack[i] += sack[i-coin]
        }
    }
    return sack[target]
}
```

### Known Bugs