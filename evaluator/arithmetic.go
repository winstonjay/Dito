package evaluator

import (
	"dito/object"
	"math"
)

const (
	// MaxInt64 : 9223372036854775807
	MaxInt64 = 1<<63 - 1
	// MinInt64 : -9223372036854775807
	MinInt64 = -(MaxInt64 + 1)
)

// FloatObjPow : Just use go's math power function.
func FloatObjPow(base, exponent float64) *object.Float {
	return object.NewDitoFloat(math.Pow(base, exponent))
}

// IntegerObjPow :
func IntegerObjPow(a, b int64) object.Object {
	base, exponent := a, b
	var result int64 = 1
	for exponent > 0 {
		if exponent%2 == 1 {
			if Int64MulOverflows(result, base) {
				goto OverlowError
			}
			result *= base
		}
		exponent /= 2
		if Int64MulOverflows(base, base) {
			goto OverlowError
		}
		base *= base
	}
	return object.NewDitoInteger(result)

OverlowError:
	return newError("OverlowError in operation: %d ** %d", a, b)
}

// Int64SumOverflows :
func Int64SumOverflows(a, b int64) bool {
	return (a > 0 && a > MaxInt64-b) || (a < 0 && a < MinInt64-b)
}

// Int64DiffOverflows :
func Int64DiffOverflows(a, b int64) bool {
	if a < 0 && a > MaxInt64+b {
		return true
	}
	if a > 0 && a < MinInt64+b {
		return true
	}
	return false
}

// Int64MulOverflows :
func Int64MulOverflows(a, b int64) bool {
	if a == 0 || b == 0 || a == 1 || b == 1 {
		return false
	}
	if a == MinInt64 || b == MinInt64 {
		return true
	}
	c := a * b
	return c/b != a
}
