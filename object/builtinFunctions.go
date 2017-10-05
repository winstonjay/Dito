package object

import (
	"math"
)

// DitoIntAbs :
func DitoIntAbs(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

// DitoFloatAbs :
func DitoFloatAbs(value float64) float64 { return math.Abs(value) }

//

// DitoIntMax :
func DitoIntMax(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// DitoFloatMax :
func DitoFloatMax(x, y float64) float64 { return math.Max(x, y) }

//

// DitoIntMin :
func DitoIntMin(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// DitoFloatMin :
func DitoFloatMin(x, y float64) float64 { return math.Min(x, y) }
