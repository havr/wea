package util

import "math"

// Round rounds the given number keeping given number of digits after coma.
func Round(n float64, precision int) float64 {
	m := math.Pow10(precision)
	return math.Round(n*m) / m
}
