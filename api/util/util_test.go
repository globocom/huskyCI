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

})
