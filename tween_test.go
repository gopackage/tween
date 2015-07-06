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
			c := &Tween{}

			Ω(c).ShouldNot(BeNil())
		})
	})
})
