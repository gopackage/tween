package updaters

import (
	"image/color"
	"time"

	"github.com/gopackage/tween"
)

// NewColor creates a new color updater with the provided colors and
// initializes unbuffered channels for Updates and Done signal.
func NewColor(from, to color.Color) *Color {
	r, g, b, a := from.RGBA()
	f := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	r, g, b, a = to.RGBA()
	t := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	return &Color{
		From:    f,
		To:      t,
		Updates: make(chan color.RGBA),
		Done:    make(chan int),
	}
}

// Color provides tween support for colors.
type Color struct {
	From    color.RGBA      // From the color we transition from
	To      color.RGBA      // To the color we transition to
	Updates chan color.RGBA // A channel that receives color updates
	Done    chan int        // A channel to receive a done signal

	from color.RGBA // from is the starting color snapshot
	to   color.RGBA // to is the ending color snapshot
	r    float64    // r is the total red transition
	g    float64    // r is the total green transition
	b    float64    // b is the total blue transition
	a    float64    // b is the total alpha transition
}

// Start begins the color update.
func (c *Color) Start(framerate, frames int, frameTime, runningTime time.Duration) {
	// Snapshot the color values - just in case someone tries to change it
	c.from = c.From
	c.to = c.To
	// Calculate how much each color changes during the tween
	c.r = float64(int(c.to.R) - int(c.from.R))
	c.g = float64(int(c.to.G) - int(c.from.G))
	c.b = float64(int(c.to.B) - int(c.from.B))
	c.a = float64(int(c.to.A) - int(c.from.A))
}

// Update interpolates the color between start and end.
func (c *Color) Update(frame tween.Frame) {
	c.Updates <- color.RGBA{
		R: c.from.R + uint8(c.r*frame.Transitioned),
		G: c.from.G + uint8(c.g*frame.Transitioned),
		B: c.from.B + uint8(c.b*frame.Transitioned),
		A: c.from.A + uint8(c.a*frame.Transitioned),
	}
}

// End terminates the color updates.
func (c *Color) End() {
	close(c.Done)
}
