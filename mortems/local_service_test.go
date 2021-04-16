package mortems_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/mentimeter/morty/mortems"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mortems", func() {
	var mortems MortemCollector
	var localFileService RepoFileService
	var tmpDir string

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "testlocalmortem")
		Expect(err).NotTo(HaveOccurred())

		_, err = copyDir("testdata/different-month-mortems/", tmpDir)
		Expect(err).NotTo(HaveOccurred())

		localFileService = NewLocalFileService(tmpDir)
		mortems = NewMortemCollector(localFileService)
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	It("saves the correct metrics from the mortem", func() {
		Expect(mortems.Collect()).To(Succeed())

		mortemFiles, err := localFileService.GetFiles()
		Expect(err).NotTo(HaveOccurred())

		dbFile := mortemFiles.GetFile("mortems.json")
		Expect(dbFile).NotTo(BeNil())

		var mortemEntries []MortemData

		dbFileBytes := []byte(dbFile.GetContent())
		Expect(json.Unmarshal(dbFileBytes, &mortemEntries)).To(Succeed())

		Expect(mortemEntries).To(ContainElement(FirstMortem()))
		Expect(mortemEntries).To(ContainElement(ThirdMortem()))
	})

	It("creates the README in the root directory", func() {
		Expect(mortems.Collect()).To(Succeed())

		mortemFiles, err := localFileService.GetFiles()
		Expect(err).NotTo(HaveOccurred())

		fullReadmePath := "README.md"
		Expect(mortemFiles.GetFile(fullReadmePath)).To(BeFileWithSubstring("Love Lost Globally: Jerry Develops Malicious App"))
		Expect(mortemFiles.GetFile(fullReadmePath)).To(BeFileWithSubstring("Christmas Lighting Causes Near Death"))
		Expect(mortemFiles.GetFile(fullReadmePath)).To(BeFileWithSubstring("July 2020"))
		Expect(mortemFiles.GetFile(fullReadmePath)).To(BeFileWithSubstring("August 2020"))
	})

	It("creates the local morty install script", func() {
		Expect(mortems.Collect()).To(Succeed())

		mortemFiles, err := localFileService.GetFiles()
		Expect(err).NotTo(HaveOccurred())

		Expect(mortemFiles.GetFile("install_morty").GetContent()).To(ContainSubstring("wget"))
		Expect(mortemFiles.GetFile("install_morty").GetMode()).To(Equal("100755"))
	})
})

func copyDir(source, dest string) (string, error) {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -R %s %s", source, dest)).Output()

	return string(output), err
}
