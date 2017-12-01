# Dito Programming Language

Dito is a work-in-progress toy interpeted language made to learn about writing parsers, interpeters etc. Things change around a bit so some docs arn't up to date. The syntax is
pretty generic at the moment but probally will try out some more idiosyncratic concepts in the future.

### Having a go...

The most simple way to run is with `go run`. You can try out the language either in a interactive shell with:

    go run main.go

Or you can run a file from `testdata/` directory when in the main directory with. (The default path is `testdata/`)

    go run main.go <filename>


### Program Examples:

`greet` lambda function, as you might suspect, it outputs 'Hello, Dito!'
```go
greet = func(name) -> ("Hello, " + name)
print(greet("Dito!"))
```

`newton_sqrt` function, Calculate the square root of `x` to an accuracy of
`delta` using newtons method.
```go
func newton_sqrt(x, delta) {
    z = delta
    for abs(x - z**2) > delta {
        z -= (z**2 - x) / (2 * x)
    }
    return z
}
```

```go
func collatz_sequence(start) {
    collatz = func(n) -> (
        n           if n <= 1 else
        n / 2       if n % 2 == 0 else
        3 * n + 1)
    sequence = [start]
    for sequence[-1] > 1 {
        newitem = collatz(sequence[-1])
        sequence += [newitem]
    }
    return sequence
}
print(collatz_sequence(13))
```
outputs `[13, 40, 20, 10, 5, 16, 8, 4, 2, 1]`

**For more examples see `testdata/` directory.**

