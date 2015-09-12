package tween_test

import (
	"time"

	. "github.com/gopackage/tween"
	"github.com/gopackage/tween/curves"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Recorder struct {
	Frames      []Frame
	FPS         int
	TotalFrames int
	FTime       time.Duration
	Running     time.Duration
	Done        chan int
}

func (u *Recorder) Start(framerate, frames int, frameTime, runningTime time.Duration) {
	u.FPS = framerate
	u.TotalFrames = frames
	u.FTime = frameTime
	u.Running = runningTime
}

func (u *Recorder) Update(frame Frame) {
	u.Frames = append(u.Frames, frame)
}

func (u *Recorder) End() {
	u.Done <- 1
}

var _ = Describe("Core", func() {
	Describe("Engine", func() {
		It("should generate frames", func(done Done) {
			d := make(chan int)
			recorder := &Recorder{Done: d}
			engine := NewEngine(time.Second, curves.Linear, recorder)
			engine.Start()
			<-d
			Ω(recorder.FPS).Should(Equal(60))
			Ω(recorder.TotalFrames).Should(Equal(60))
			Ω(recorder.FTime).Should(Equal(16666666 * time.Nanosecond))
			Ω(recorder.Running).Should(Equal(time.Second))
			last := recorder.Frames[len(recorder.Frames)-1]
			Ω(last.Index).Should(Equal(60))
			Ω(last.Completed).Should(Equal(1.))
			Ω(last.Transitioned).Should(Equal(1.))
			Ω(last.Elapsed).Should(Equal(time.Second))
			//Ω(recorder.Frames).Should(HaveLen(61))
			//Ω(recorder.Frames).Should(Equal([]Frame{}))
			close(done)
		}, 2)
	})
})
