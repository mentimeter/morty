package mortems_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ostenbom/morty/mortems"
)

var _ = Describe("Mortem Content Parsing", func() {
	Describe("title parsing", func() {
		It("can parse correct title", func() {
			title, err := ParseTitle("# It's alive")
			Expect(err).NotTo(HaveOccurred())
			Expect(title).To(Equal("It's alive"))
		})

		It("returns the correct error", func() {
			_, err := ParseTitle("It's dead")
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrNoTitle)).To(BeTrue())
		})
	})
})
