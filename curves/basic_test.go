package curves_test

import (
	. "github.com/gopackage/tween/curves"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basic Curves", func() {
	Describe("Linear", func() {
		It("should generate a linear curve", func() {
			linear := Linear{}
			Ω(linear.Transition(.0)).Should(Equal(.0))
			Ω(linear.Transition(.1)).Should(Equal(.1))
			Ω(linear.Transition(.2)).Should(Equal(.2))
			Ω(linear.Transition(.3)).Should(Equal(.3))
			Ω(linear.Transition(.4)).Should(Equal(.4))
			Ω(linear.Transition(.5)).Should(Equal(.5))
			Ω(linear.Transition(.6)).Should(Equal(.6))
			Ω(linear.Transition(.7)).Should(Equal(.7))
			Ω(linear.Transition(.8)).Should(Equal(.8))
			Ω(linear.Transition(.9)).Should(Equal(.9))
			Ω(linear.Transition(1.)).Should(Equal(1.))
		})
	})
})
