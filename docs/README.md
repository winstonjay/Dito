# Language specification

This is currently in progress. What is here is limited but should be reasonably accurate.

### Keywords

The following keywords are reserved and may not be used as identifiers.

```
and     if      else        for         in          def
or      mut     return      not         import      let
```

### Operators

The following character sequences represent operators (including assignment operators) and punctuation:

```
+   -   *   /   //  %   **  <<  >>  &   |   ^
==  !=  <   >   <=  >=  ++  ->  (   )   [   ]
{   }   ;   ,   =   +=  -=  /=  %=  "
```

**Operator precedence**

precedence levels dictate the order of operations in expression. Here they are ranked lowest (least binding) to highghest (most binding). Operators with the same precedence are executed left to right.

```
1           if
2           or and
3           == != <= >= < >
4           + - ^ | &
5           * / % // << >>
6           **
7           not + -
8           bracketed expressions and function calls.
```

**Arithmetic operators**

```
+       ADD         int, float
-       SUB         int, float
*       MUL         int, float
/       DIV         int, float
//      IDIV        int
%       MOD         int
**      POW         int, float
<<      LSHIFT      int
>>      RSHIFT      int
&       BITAND      int
|       BITOR       int
^       BITXOR      int
```

**Prefix operators**

```
+       ADD         int, float
-       SUB         int, float
not     NOT         int, float, bool
```

**Comparison operators**

```
==      EQUALS      int, float, string, bool
!=      NEQUALS     int, float, string, bool
<       LTHAN       int, float, bool
>       GTHAN       int, float, bool
<=      LEQUALS     int, float, bool
>=      GEQUALS     int, float, bool
or      OR          bool
and     AND         bool
```

**Array operators**

```
++      CAT         array
in      IN          array
```

```
dito: let arr = [1, 4, 9, 25, 36]
dito: arr[1]
4
dito: arr[1:3]
[4, 9]
dito: arr[0] *= 5
dito: arr
[5, 4, 9, 25, 36]
```

**Assignment operators**

```
=       EQUAL       any
+=      ADDEQUAL    int, float
-=      SUBEQUAL    int, float
*=      MULEQUAL    int, float
/=      DIVEQUAL    int, float
%=      MODEQUAL    int
```

## Assignments

Varibles are immutable by defualt and created via the `let` keyword.

```
let x = 100
```

trying to update this varible this would cause an error.

```
x = 300
```

To create a mutable varible the keyword `mut` must also be added.

```
let mut y = 200
```

As we have defined y as mutable this would be fine to do or even any of the inplace operators.

```
y = 3 - y
y += x
y %= 100
```

## Functions

functions are first class citizens and they are like the functions of
most C like languages. They have their own scope, parameted etc. There are block functions, expression functions and also a limitied number of built in functions. There is not any difference between the function types under the hood except the limitation to single expressions for lambdas.

**Block functions**

Block functions define a subroutine that is contained within a block statements. Functions can contain other functions, the result of the last statement is allways returned from the function unless an explicit `return` statement is given.

```
def newton_sqrt(x, delta) {
    let mut z = delta
    for abs(x - z**2) > delta {
        z -= (z**2 - x) / (2 * x)
    }
    return z
}
```

**Expression functions (lambdas)**

Expression functions or lambdas are un-named functions that are limited to a single expression. These allways return a value.

```
let sqr = def(x) -> x * x
```

```
let collatz = def(n) -> (
    n           if n <= 1 else
    n / 2       if n % 2 == 0 else
    3 * n + 1)
```

**Composition**

both types of functions are first class citizens so they can be passed
and returned by other functions. For example...

```
dito: let mul = def(a) -> def(b) -> a * b
dito: mul(10)(10)
100
```

This uses the `std` library function map, which takes a function and a
array as arguements.

```
dito: map(def(x) -> x * x, [1, 2, 3, 4, 5])
[1, 4, 9, 16, 25]
```

## The rest

TODO...

