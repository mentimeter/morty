package mortems_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ostenbom/morty/mortems"
	"github.com/ostenbom/morty/mortems/mortemsfakes"
)

var _ = Describe("Mortems", func() {
	var mortems MortemCollector
	var gitService *mortemsfakes.FakeGitService

	BeforeEach(func() {
		gitService = new(mortemsfakes.FakeGitService)
		mortems = NewMortemCollector(gitService)
	})

	It("Gets the TreeEntries from Git", func() {
		Expect(mortems.Collect()).To(Succeed())
		Expect(gitService.GetTreeEntriesCallCount()).To(Equal(1))
	})
})
