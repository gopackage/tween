// Package curves adds standard transition curve implementations for
// the most commonly used tweens.
package curves

import "math"

//go:generate go run gen/gen.go

// Linear calculates a linear transition function.
func Linear(completed float64) float64 {
	return completed
}

// Swing is a simple ease-in-ease-out transition that provides minimal curvature
// at the beginning and end of the transition.
func Swing(completed float64) float64 {
	return 0.5 - math.Cos(completed*math.Pi)/2
}
