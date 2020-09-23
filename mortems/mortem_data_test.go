package mortems_test

import (
	"errors"
	"time"

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

	Describe("owner parsing", func() {
		It("can parse the owner", func() {
			owner, err := ParseOwner("Owner: Oliver Stenbom")
			Expect(err).NotTo(HaveOccurred())
			Expect(owner).To(Equal("Oliver Stenbom"))
		})
	})

	Describe("date parsing", func() {
		It("can parse the date", func() {
			date, err := ParseDate("Date: June 2, 2020")
			Expect(err).NotTo(HaveOccurred())
			Expect(date).To(Equal(time.Date(2020, time.June, 2, 0, 0, 0, 0, time.UTC)))
		})
	})

	Describe("severity parsing", func() {
		It("can parse the severity", func() {
			sev, err := ParseSeverity("| Severity | 1 |")
			Expect(err).NotTo(HaveOccurred())
			Expect(sev).To(Equal("1"))
		})
	})
})
