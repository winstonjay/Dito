# Dito Programming Language

Dito is a work-in-progress toy interpeted language made to learn about writing parsers, interpeters etc. Things change around a bit so some docs arn't up to date. The syntax is pretty generic at the moment but probally will try out some more idiosyncratic concepts in the future.

This is the most up to date repo but its not that stable.

### Having a go...

The most simple way to build with `go build` then use the executable or just `go run main.go`. You can try out the language either in a interactive shell with:

    ./dito

Or you can run a file from `exampledata/` directory when in the main directory with. (The default path is `exampledata/`)

    ./dito funcs.dito


**For more examples see `exampledata/` directory.**


TODO: idk how on earth this has happend.

    (dito)> x
    [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
    (dito)> x[0] * 2
    0
    (dito)> x[0] *= 2
    (dito)> x
    [2, 1, 2, 3, 4, 5, 6, 7, 8, 9]
    (dito)> 0 * 2
    0

sort this: inplace operators in general are a mess.