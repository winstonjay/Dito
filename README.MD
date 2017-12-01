# Dito Programming Language

Dito is a work-in-progress toy interpeted language made to learn about writing parsers, interpeters etc. Things change around a bit so some docs arn't up to date. The syntax is
pretty generic at the moment but probally will try out some more idiosyncratic concepts in the future.

### Having a go...

The most simple way to run is with `go run`. You can try out the language either in a interactive shell with:

    go run main.go

Or you can run a file from with in the `testdata/` directory with:

    go run main.go <filename>


### Program Examples:

`greet` lambda function, as you might suspect, it outputs 'Hello, Dito!'
```python
greet = func(name) -> ("Hello, " + name)
print(greet("Dito!"))
```

`newton_sqrt` function, Calculate the square root of `x` to an accuracy of
`delta` using newtons method.
```python
func newton_sqrt(x, delta) {
    z = delta
    for abs(x - z**2) > delta {
        z -= (z**2 - x) / (2 * x)
    }
    return z
}
```

```python
func collatz_sequence(start) {
    # define an array with our start item.
    # define a lambda function that can only be a single expression.
    # here it uses a if else expression like what can be found in python.
    collatz = func(n) -> (
        n           if n <= 1 else
        n / 2       if n % 2 == 0 else
        3 * n + 1)
    sequence = [start]
    # use a while style for loop to keep adding to the sequence.
    # checking the item at the last index of the sequence.
    for sequence[-1] > 1 {
        # Get the next value in the sequence.
        newitem = collatz(sequence[-1])
        # Append each new item to array.
        sequence += [newitem]
    }
    return sequence
}
print(collatz_sequence(13))
# outputs [13, 40, 20, 10, 5, 16, 8, 4, 2, 1]
```

**For more examples see `testdata/` directory.**

