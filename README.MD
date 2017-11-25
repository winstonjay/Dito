# Dito Programming Language

Dito is a


### Examples:

`greet` lambda function as you
```go
greet = func(name) -> ("Hello, " + name)
print(greet("Dito!"))
# outputs Hello, Dito!
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


### Known Bugs

* improper ending of expressions.
* if else statements not working without semi colons linked ot above.
* if else expression precedence. need to find a new format.