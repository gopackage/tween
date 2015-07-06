package tween

import "time"

// Tween drives timed transitions between two states. Tween currently supports
// start and end states that are color.Color objects.
type Tween struct {
	FrameRate int           // FrameRate is number of "display" frames per second
	Duration  time.Duration // Duration is the time that the tween takes place
	Start     interface{}   // Start is the starting state of the tween
	End       interface{}   // End is the ending state of the tween
}
