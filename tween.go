package tween

import "image/color"
import "time"
import "fmt"

type TweenData struct {
	duration  int64 // milliseconds
	accum     int64 // how many milliseconds have elapsed
	pace      int64 // how many milliseconds we should be at
	startTime time.Time
	running   bool
	done      bool
}

type ColorTweener struct {
	Tween       TweenData
	Starting    color.RGBA
	Ending      color.RGBA
	current     color.RGBA
	Channel     chan *color.RGBA
	stopChannel chan int
	doneChannel chan int
}

// Duration in seconds
func NewColorTweener(duration int64, start, stop color.RGBA) (*ColorTweener, <-chan *color.RGBA, <-chan int) {
	channel := make(chan *color.RGBA, 10)
	doneChannel := make(chan int)
	c := &ColorTweener{
		Tween: TweenData{
			duration:  duration, // save duration in nanoseconds
			accum:     0,
			pace:      0,
			startTime: time.Now(),
			running:   false,
			done:      false,
		},
		Starting:    start,
		Ending:      stop,
		current:     start,
		Channel:     channel,
		stopChannel: make(chan int),
		doneChannel: doneChannel,
	}

	return c, channel, doneChannel
}

func getNewColor(start, end color.RGBA, percent float64) color.RGBA {
	spanRed := int(end.R) - int(start.R)
	spanGreen := int(end.G) - int(start.G)
	spanBlue := int(end.B) - int(start.B)

	newC := color.RGBA{
		R: start.R + uint8(float64(spanRed)*percent),
		G: start.G + uint8(float64(spanGreen)*percent),
		B: start.B + uint8(float64(spanBlue)*percent),
		A: 255,
	}
	return newC
}

func (c *ColorTweener) Start() {
	if c.Tween.running != false || c.Tween.done != false {
		return
	}
	// can't stop this thread unless you call Stop() or let the timer
	// run out
	go func() {
		// set start time
		c.Tween.startTime = time.Now()

		// start initial timer
		timeChan := time.NewTimer(time.Millisecond).C

		// set pace what we should end up at
		c.Tween.pace += 1

		for !c.Tween.done {
			select {
			case <-timeChan:
				c.Tween.pace += 1
				c.Tween.accum = int64(time.Now().Sub(c.Tween.startTime) / time.Millisecond)
				// get new color and send color update
				percentageComplete := (float64(c.Tween.accum) / float64(c.Tween.duration))
				if percentageComplete > 1 {
					percentageComplete = 1
				}
				fmt.Printf("percentage complete: %f\n", percentageComplete)
				c.current = getNewColor(c.Starting, c.Ending, percentageComplete)
				c.Channel <- &c.current

				// see if we should keep going
				if c.Tween.pace != c.Tween.duration {
					// figure out how long we should run new timer for
					newTimerVal := c.Tween.pace - c.Tween.accum // in milliseconds
					if newTimerVal < 0 {
						newTimerVal = 0
					}

					// run timer
					timeChan = time.NewTimer(time.Duration(newTimerVal) * time.Millisecond).C
				} else {
					// done, send on stopChannel to complete
					c.Tween.done = true
				}
			case <-c.stopChannel:
				c.Tween.done = true
			}
		}

		// do cleanup
		c.Channel <- &c.Ending
		c.doneChannel <- 1
	}()
}

func (c *ColorTweener) Stop() {
	if c.Tween.running == true {
		c.stopChannel <- 1
	}
}
