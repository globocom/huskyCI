package util_test

import (
	"github.com/globocom/huskyCI/api/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {

	Describe("HandleCmd", func() {
		inputRepositoryURL := "https://github.com/globocom/secDevLabs.git"
		inputRepositoryBranch := "myBranch"
		internalDepURL := "https://myinternalurl.com"
		inputCMD := "git clone -b %GIT_BRANCH% --single-branch %GIT_REPO% code --quiet 2> /tmp/errorGitCloneRetirejs -- %INTERNAL_DEP_URL%"
		expected := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitCloneRetirejs -- https://myinternalurl.com"
		expectedEmptyDepURL := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitCloneRetirejs -- "

		Context("When inputRepositoryURL, inputRepositoryBranch, internalDepURL and inputCMD are not empty", func() {
			It("Should return a string based on these params", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, internalDepURL, inputCMD)).To(Equal(expected))
			})
		})
		Context("When inputRepositoryURL is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd("", inputRepositoryBranch, internalDepURL, inputCMD)).To(Equal(""))
			})
		})
		Context("When inputRepositoryBranch is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd(inputRepositoryURL, "", internalDepURL, inputCMD)).To(Equal(""))
			})
		})
		Context("When inputCMD is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, internalDepURL, "")).To(Equal(""))
			})
		})
		Context("When internalDepURL is empty", func() {
			It("Should return expectedEmptyDepURL", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, "", inputCMD)).To(Equal(expectedEmptyDepURL))
			})
		})
	})
})
