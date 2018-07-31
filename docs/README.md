# Language specification

This is currently in progress. What is here is limited but should be reasonably accurate.

### Keywords

The following keywords are reserved and may not be used as identifiers.

```
and     if      else        for         in          func
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

```
1           if
2           or and
3           == != <= >= < >
4           + - ^ | &
0           * / % // << >>
5           **
6           not + -
7           bracketed expressions and function calls.
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

## The rest

TODO...

