package tween_test

import (
	. "../tween"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Color Manager", func() {
		It("should convert colors to PWM values", func() {
			// c, err := controller.NewController()
			// Ω(err).ShouldNot(HaveOccurred())
			c, colors, done := NewColorTween(600, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 128, 255})

			running := true
			var color *color.RGBA
			for running {
				select {
				case c = <-colors:
					fmt.Printf("New color: %+v\n", c)
				case <-done:
					running = false
				}
			}

			Ω(c).ShouldNot(BeNil())
		})
	})
})
