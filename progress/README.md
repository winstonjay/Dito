# Language Todo list.


## Language Features

- implement `and` and `or` operators for all atomic types.
- implement while style `for` loops like go.
- implement `hashmaps`.
- implement simple slice operations like python/go (e.g. `arr[:4]`). building block for array operations.
- Change `Numeric` base type to `Atomic` and enforce that all these types can perform math ops.
- Resolve type heirarchy issues.
- fix tracebacks so that column position is correct and at the start of the token.


# Environment Features
- command-line help. (requires docs to be written)
- docs.

## Tests
See `cover.html` for code coverage.



# Bugs

### 1

Trailing comments like this produce an error. This shouldnt be a problem.

```
let MAX_INT =  (1 << 63) - 1 #  9223372036854775807
let MIN_INT =  (1 << 63)     # -9223372036854775808
```

```
PARSE ERROR:
Traceback line 10 column 2:
    let
      ^ Is your problem here?
Expected next token is 'NEWLINE'. got 'let' instead
```

### 2

Scope resolution issues. can get but cannot update varibles in outer scopes.
need to decide if this should be able to be done. e.g.:

```
let mut counter = 0
func inc() {
    counter += 1
    return counter
}
```

if so the result should be [1, 2, 3], instead we get:

```
(dito)> [inc(), inc(), inc()]
[1, 1, 1]
```

This seems inconsistent to me, varibles should either be
inaccessible for re assignement be able to be updated by a child scope.
Having constants as we do should stop any unintended assignment issues.

### 3

operation causes runtime panic.

```
(dito)> "x" ++ ["hello", "sup"]
panic: interface conversion: object.Object is *object.String, not *object.Array
```