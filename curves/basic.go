// Package curves adds standard Transitioner curve implementations for
// the most commonly used tweens.
package curves

// Linear is a simple one-to-one linear transition.
type Linear struct {
}

// Transition produces the simplest linear transition curve.
func (l *Linear) Transition(completed float64) float64 {
	return completed
}
