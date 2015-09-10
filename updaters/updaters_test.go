package updaters_test

import (
	"image/color"
	"time"

	"github.com/gopackage/tween"
	"github.com/gopackage/tween/curves"
	. "github.com/gopackage/tween/updaters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {
	Describe("Color Tween", func() {
		It("should generate color tween values", func() {
			start := color.RGBA{255, 0, 0, 255}
			end := color.RGBA{0, 128, 255, 0}

			updater := NewColor(start, end)
			engine := tween.NewEngine(time.Second, &curves.Linear{}, updater)
			engine.Start()

			running := true
			colors := []color.RGBA{}
			for running {
				select {
				case color := <-updater.Updates:
					colors = append(colors, color)
				case <-updater.Done:
					running = false
				}
			}
			Ω(colors).Should(HaveLen(61))
		})
	})
})
