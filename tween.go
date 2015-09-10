package tween

import "time"

// Transitioner is the interface for calculating tweening transitions.
type Transitioner interface {
	// Transition calculates the percentage of the transition between the start
	// and end values based tween (elapsed time) completion status.
	// For example, a linear tween simply has a 1:1 ratio between completed
	// and transition percentages and returns the completed value unchanged.
	// An "ease in" transition may create a logorithmic relationship between
	// the completion time and transition value for the first third, then a
	// linear relationship for the remainder.
	Transition(completed float64) float64
}

// Updater is the interface for updating the current value as it tweens between
// start and end values of a Tween.
type Updater interface {
	// Start signals the beginning of a tween and is sent before the tweening
	// begins. Start may be used to setup or pre-calculate updates.
	//
	// framerate is the number of frames per second in the tween
	// frames is the total number of frames that be generated
	// frameTime is the duration for each frame
	// runningTime is the total duration for the entire tween
	Start(framerate, frames int, frameTime, runningTime time.Duration)
	// Update receives information about the current Tween Frame and should be
	// used to update output or state.
	Update(Frame Frame)
	// End signals the end of the tween and is called after all updates.
	// End may be used to clean up resources (e.g. update channels).
	End()
}

// Frame captures information about the current "frame" of a tween transition.
type Frame struct {
	Completed    float64       // Completed is the percentage 0.0 - 1.0 of elapsed time.
	Transitioned float64       // Transitioned is the percentage 0.0 - 1.0 of transition between start and end values of the tween.
	Index        int           // Index is the current frame index
	Elapsed      time.Duration // Elapsed is the current elapsed time in the tween.
}

// NewEngine creates a basic tween Engine with a framerate of 60fps.
func NewEngine(duration time.Duration, transitioner Transitioner, updater Updater) *Engine {
	return &Engine{
		Duration:     duration,
		Transitioner: transitioner,
		Updater:      updater,
		Framerate:    60,
	}
}

// Engine runs a tween relying on transitioner and updater.
type Engine struct {
	Duration     time.Duration // The total duration of the tween.
	Framerate    int           // The number of tween data points per second (defaults to 60 fps - like the real gamers use).
	Transitioner Transitioner  // Transitioner calculates the transition curve for the tween.
	Updater      Updater       // Updater updates the tween values for each frame.

	running bool     // True if the tween is running
	done    chan int // Internal channel used to terminate the tween early
}

// Start begins the tween running.
func (e *Engine) Start() {
	e.done = make(chan int)
	// can't stop this thread unless you call Stop() or let the timer
	// run out
	go func() {
		// Setup internal done channel
		e.done = make(chan int)

		// Based on fps we can calculate how long a frame is:
		frameDuration := time.Second / time.Duration(e.Framerate) // The duration in a frame
		cutoff := e.Duration - frameDuration                      // The cutoff point where elapsed time is considered "done"
		frames := int(e.Duration / frameDuration)                 // The number of frames in the duration

		// start ticker
		e.running = true
		e.Updater.Start(e.Framerate, frames, frameDuration, e.Duration)

		// Send initial frame
		frame := Frame{}
		e.Updater.Update(frame)

		// set start time
		ticker := time.NewTicker(frameDuration)
		timeChan := ticker.C
		started := time.Now()

		for e.running {
			select {
			case <-timeChan:
				frame.Elapsed = time.Since(started)

				// Calculate the frame index - some frames can be skipped so
				// must find correct time slot for this elapsed time
				frame.Index = int(frame.Elapsed / frameDuration)

				// Calculate the completed percentage of time
				frame.Completed = ((float64(frame.Index) * float64(frameDuration)) / float64(e.Duration))
				if frame.Completed > 1 {
					go e.Stop()
					break
				}

				// Calulate the completed percentage of the transition
				frame.Transitioned = e.Transitioner.Transition(frame.Completed)

				// Update the value
				e.Updater.Update(frame)

				// see if we should keep going
				if frame.Elapsed > cutoff {
					go e.Stop() // terminate ourself
				}
			case <-e.done:
				ticker.Stop()
				e.running = false
			}
		}

		// cleanup
		frame.Elapsed = e.Duration
		frame.Completed = 1
		frame.Transitioned = 1
		e.Updater.Update(frame)
		e.Updater.End()
	}()
}

// Stop terminates the tween immediately.
func (e *Engine) Stop() {
	if e.running == true {
		e.done <- 1
	}
}
