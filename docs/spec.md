# Working Dito Language Specification

A lot of stuff is under construction and changing but some rough features are presented here.

## Primitive datatypes
The following types are provided: `int`, `float`, `string`, `array`. Numerical types have a promotional hierachy where an `int` is converted to a `float` during divisions that produce decimal numbers. To perform purely `int` division use the `//` operator. `int`, `float` and `string` are the same as go's `int64`, `float64` and `string`. Here the `array` type can be a mixture of any of the other types.



