package mortems_test

import (
	"encoding/json"
	"time"

	. "github.com/mentimeter/morty/mortems"
	"github.com/mentimeter/morty/mortems/mortemsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Mortems", func() {
	var mortems MortemCollector
	var gitService *mortemsfakes.FakeGitService
	var treeEntryFixtures map[string]*RepoFiles

	BeforeEach(func() {
		var err error
		treeEntryFixtures, err = LoadTreeEntryFixtures()
		Expect(err).NotTo(HaveOccurred())

		gitService = new(mortemsfakes.FakeGitService)
		mortems = NewMortemCollector(gitService)
	})

	Context("an empty, uninitialized repository", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["empty-repo"], nil)
		})

		It("creates the post-mortem template in the right place", func() {
			Expect(mortems.Collect()).To(Succeed())

			Expect(gitService.GetFilesCallCount()).To(Equal(1))

			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("post-mortems/template.md")).To(BeFileWithSubstring("<!-- The title of your incident. "))
			Expect(updateFiles.GetFile("post-mortems/template.md")).To(BeFileWithSubstring("Love Lost Globally"))
		})

		It("creates the post-mortem how-to/README", func() {
			Expect(mortems.Collect()).To(Succeed())

			Expect(gitService.GetFilesCallCount()).To(Equal(1))

			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("post-mortems/README.md")).To(BeFileWithSubstring("# Post Mortem Action Plan"))
		})

		It("links to how-to from README", func() {
			Expect(mortems.Collect()).To(Succeed())
			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("#### [How to Post-Mortem](/post-mortems)"))
		})
	})

	Context("old version of mortems README", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["outdated-readme"], nil)
		})

		It("re-creates the post-mortem how-to/README", func() {
			Expect(mortems.Collect()).To(Succeed())

			Expect(gitService.GetFilesCallCount()).To(Equal(1))

			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("post-mortems/README.md")).To(BeFileWithSubstring("# Post Mortem Action Plan"))
		})
	})

	Context("single post mortem", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["basic-single-mortem"], nil)
		})

		It("saves the correct metrics from the mortem", func() {
			Expect(mortems.Collect()).To(Succeed())
			mortemEntries := GetMortemEntries(gitService)

			Expect(mortemEntries).To(ContainElement(FirstMortem()))
		})

		It("creates the README in the root directory", func() {
			Expect(mortems.Collect()).To(Succeed())
			Expect(gitService.GetFilesCallCount()).To(Equal(1))
			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("Love Lost Globally: Jerry Develops Malicious App"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("July 2020"))
		})
	})

	Context("two post mortems in the same month", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["two-close-mortems"], nil)
		})

		It("saves the correct metrics from the mortem", func() {
			Expect(mortems.Collect()).To(Succeed())
			mortemEntries := GetMortemEntries(gitService)

			Expect(mortemEntries).To(ContainElement(FirstMortem()))
			Expect(mortemEntries).To(ContainElement(SecondMortem()))
		})

		It("creates the README in the root directory", func() {
			Expect(mortems.Collect()).To(Succeed())
			Expect(gitService.GetFilesCallCount()).To(Equal(1))
			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("Love Lost Globally: Jerry Develops Malicious App"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("Bad Parenting: Rick Clones Own Daughter"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("July 2020"))
		})
	})

	Context("two post mortems on different months", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["different-month-mortems"], nil)
		})

		It("saves the correct metrics from the mortem", func() {
			Expect(mortems.Collect()).To(Succeed())
			mortemEntries := GetMortemEntries(gitService)

			Expect(mortemEntries).To(ContainElement(FirstMortem()))
			Expect(mortemEntries).To(ContainElement(ThirdMortem()))
		})

		It("creates the README in the root directory", func() {
			Expect(mortems.Collect()).To(Succeed())
			Expect(gitService.GetFilesCallCount()).To(Equal(1))
			Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

			updateFiles := gitService.CommitNewFilesArgsForCall(0)
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("Love Lost Globally: Jerry Develops Malicious App"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("Christmas Lighting Causes Near Death"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("July 2020"))
			Expect(updateFiles.GetFile("README.md")).To(BeFileWithSubstring("August 2020"))
		})
	})

	// TODO: For mortem that has changed its file name
})

func BeFileWithSubstring(contentSubstring string) types.GomegaMatcher {
	return And(
		Not(BeNil()),
		WithTransform(GetMode, Equal("100644")),
		WithTransform(GetType, Equal("blob")),
		WithTransform(GetContent, ContainSubstring(contentSubstring)),
	)
}

func GetMode(e *File) string {
	return e.Mode
}

func GetType(e *File) string {
	return e.Type
}

func GetContent(e *File) string {
	return e.Content
}

func GetMortemEntries(gitService *mortemsfakes.FakeGitService) []MortemData {
	Expect(gitService.GetFilesCallCount()).To(Equal(1))
	Expect(gitService.CommitNewFilesCallCount()).To(Equal(1))

	updateFiles := gitService.CommitNewFilesArgsForCall(0)
	dbFile := updateFiles.GetFile("mortems.json")
	Expect(dbFile).NotTo(BeNil())

	var mortemEntries []MortemData

	dbFileBytes := []byte(dbFile.GetContent())
	Expect(json.Unmarshal(dbFileBytes, &mortemEntries)).To(Succeed())

	return mortemEntries
}

func FirstMortem() MortemData {
	detectTime, err := time.ParseDuration("4m")
	Expect(err).To(BeNil())
	resolveTime, err := time.ParseDuration("6h14m")
	Expect(err).To(BeNil())
	totalDownTime, err := time.ParseDuration("6h28m")
	Expect(err).To(BeNil())

	return MortemData{
		File:     "post-mortems/0001-first-mortem.md",
		Title:    "Love Lost Globally: Jerry Develops Malicious App",
		Owner:    "Morty Smith",
		Date:     time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC),
		Severity: "1",
		Detect:   detectTime,
		Resolve:  resolveTime,
		Downtime: totalDownTime,
	}
}

func SecondMortem() MortemData {
	detectTime, err := time.ParseDuration("26h")
	Expect(err).To(BeNil())
	resolveTime, err := time.ParseDuration("3m")
	Expect(err).To(BeNil())
	totalDownTime, err := time.ParseDuration("26h3m")
	Expect(err).To(BeNil())

	return MortemData{
		File:     "post-mortems/0002-second-mortem.md",
		Title:    "Bad Parenting: Rick Clones Own Daughter",
		Owner:    "Rick Sanchez",
		Date:     time.Date(2020, time.July, 27, 0, 0, 0, 0, time.UTC),
		Severity: "1",
		Detect:   detectTime,
		Resolve:  resolveTime,
		Downtime: totalDownTime,
	}
}

func ThirdMortem() MortemData {
	detectTime, err := time.ParseDuration("2h1m")
	Expect(err).To(BeNil())
	resolveTime, err := time.ParseDuration("3s")
	Expect(err).To(BeNil())
	totalDownTime, err := time.ParseDuration("2h1m3s")
	Expect(err).To(BeNil())

	return MortemData{
		File:     "post-mortems/0003-third-mortem.md",
		Title:    "Christmas Lighting Causes Near Death",
		Owner:    "Jerry Smith",
		Date:     time.Date(2020, time.August, 6, 0, 0, 0, 0, time.UTC),
		Severity: "1",
		Detect:   detectTime,
		Resolve:  resolveTime,
		Downtime: totalDownTime,
	}
}
