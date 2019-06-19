package util_test

import (
	"github.com/globocom/huskyCI/api/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {

	Describe("CreateContainerName", func() {
		inputURL := "https://github.com/globocom/secDevLabs.git"
		inputBranch := "myBranch"
		inputImage := "secdevLabs/bandit"
		expected := "globocom_secDevLabs_myBranch_bandit"

		Context("When inputURL, imputBranch and inputImage is not empty", func() {
			It("Should return a container name based on these params", func() {
				Expect(util.CreateContainerName(inputURL, inputBranch, inputImage)).To(Equal(expected))
			})
		})
		Context("When inputURL is empty", func() {
			It("Should return empty string and docker will generate a default name", func() {
				Expect(util.CreateContainerName("", inputBranch, inputImage)).To(Equal(""))
			})
		})
		Context("When inputBranch is empty", func() {
			It("Should return empty string and docker will generate a default name", func() {
				Expect(util.CreateContainerName(inputURL, "", inputImage)).To(Equal(""))
			})
		})
		Context("When inputImage is empty", func() {
			It("Should return empty string and docker will generate a default name", func() {
				Expect(util.CreateContainerName(inputURL, inputBranch, "")).To(Equal(""))
			})
		})
	})

	Describe("HandleCmd", func() {
		inputRepositoryURL := "https://github.com/globocom/secDevLabs.git"
		inputRepositoryBranch := "myBranch"
		inputCMD := "git clone -b %GIT_BRANCH% --single-branch %GIT_REPO% code --quiet 2> /tmp/errorGitCloneRetirejs"
		expected := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitCloneRetirejs"

		Context("When inputRepositoryURL, inputRepositoryBranch and inputCMD is not empty", func() {
			It("Should return a string based on these params", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, inputCMD)).To(Equal(expected))
			})
		})
		Context("When inputRepositoryURL is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.CreateContainerName("", inputRepositoryBranch, inputCMD)).To(Equal(""))
			})
		})
		Context("When inputRepositoryBranch is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.CreateContainerName(inputRepositoryURL, "", inputCMD)).To(Equal(""))
			})
		})
		Context("When inputCMD is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.CreateContainerName(inputRepositoryURL, inputRepositoryBranch, "")).To(Equal(""))
			})
		})
	})

})
