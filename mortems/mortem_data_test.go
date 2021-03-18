package mortems_test

import (
	"errors"
	"time"

	. "github.com/mentimeter/morty/mortems"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		It("doesn't matter about the heading", func() {
			owner, err := ParseOwner("### Owner: Oliver Stenbom")
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
		It("can parse the date in short month format", func() {
			date, err := ParseDate("Date: Jun 2, 2020")
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

		It("can parse the severity case-insensitive", func() {
			sev, err := ParseSeverity("| severity | 1 |")
			Expect(err).NotTo(HaveOccurred())
			Expect(sev).To(Equal("1"))
		})
	})

	Describe("detection parsing", func() {
		It("can parse the detection", func() {
			detect, err := ParseDetect("| Time to Detect | 1 day, 3 hours, 6 minutes, 5 seconds |")
			Expect(err).NotTo(HaveOccurred())
			dur, err := time.ParseDuration("27h6m5s")
			Expect(err).NotTo(HaveOccurred())
			Expect(detect).To(Equal(dur))
		})

		It("can parse the detection regardless of case", func() {
			detect, err := ParseDetect("| Time to Detect | 1 Day, 3 Hours, 6 Minutes, 5 SeCoNds |")
			Expect(err).NotTo(HaveOccurred())
			dur, err := time.ParseDuration("27h6m5s")
			Expect(err).NotTo(HaveOccurred())
			Expect(detect).To(Equal(dur))
		})

		It("can parse the detection regardless of case", func() {
			detect, err := ParseDetect("| Time To Detect  | 30 minutes          |")
			Expect(err).NotTo(HaveOccurred())
			dur, err := time.ParseDuration("30m")
			Expect(err).NotTo(HaveOccurred())
			Expect(detect).To(Equal(dur))
		})
	})

	Describe("resolve parsing", func() {
		It("can parse the detection", func() {
			resolve, err := ParseResolve("| Time to Resolve | 1 day, 3 hours, 6 minutes, 5 seconds |")
			Expect(err).NotTo(HaveOccurred())
			dur, err := time.ParseDuration("27h6m5s")
			Expect(err).NotTo(HaveOccurred())
			Expect(resolve).To(Equal(dur))
		})
	})

	Describe("downtime parsing", func() {
		It("can parse the detection", func() {
			down, err := ParseDowntime("| Total Downtime | 1 day, 3 hours, 6 minutes, 5 seconds |")
			Expect(err).NotTo(HaveOccurred())
			dur, err := time.ParseDuration("27h6m5s")
			Expect(err).NotTo(HaveOccurred())
			Expect(down).To(Equal(dur))
		})
	})
})
