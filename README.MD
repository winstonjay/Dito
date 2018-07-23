# Dito Programming Language

Dito is a work-in-progress toy interpeted language made to learn about writing parsers, interpeters etc. Things change around a bit so some docs arn't up to date. The syntax is pretty generic at the moment but probally will try out some more idiosyncratic concepts in the future.

This is the most up to date repo but its not that stable.

### Examples

For full list of examples see [`examples`](examples/) directory or the [stdlib](lib/std.dito).
Heres calcualting a square root using newtonws method.

```
# Calculate square root of x to a acuracy of delta.
func newton_sqrt(x, delta) {
    let mut z = delta
    for abs(x - z**2) > delta {
        z -= (z**2 - x) / (2 * x)
    }
    return z
}

print(newton_sqrt(25, 1/128))
```


### Having a go...

The most simple way to build with `go build` then use the executable or just `go run main.go`. You can try out the language either in a interactive shell with:

    ./dito

or to run a file:

    ./dito funcs.dito

