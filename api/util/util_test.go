package util_test

import (
	"os"

	"github.com/globocom/huskyCI/api/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {

	Describe("HandleCmd", func() {
		inputRepositoryURL := "https://github.com/globocom/secDevLabs.git"
		inputRepositoryBranch := "myBranch"
		inputCMD := "git clone -b %GIT_BRANCH% --single-branch %GIT_REPO% code --quiet 2> /tmp/errorGitCloneRetirejs"
		expected := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitCloneRetirejs"

		Context("When inputRepositoryURL, inputRepositoryBranch and inputCMD are not empty", func() {
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

	Describe("HandlePrivateSSHKey", func() {

		rawString := "echo 'GIT_PRIVATE_SSH_KEY' > ~/.ssh/huskyci_id_rsa &&"
		expectedNotEmpty := "echo 'PRIVKEYTEST' > ~/.ssh/huskyci_id_rsa &&"
		expectedEmpty := "echo '' > ~/.ssh/huskyci_id_rsa &&"

		Context("When rawString and HUSKYCI_API_GIT_PRIVATE_SSH_KEY are not empty", func() {
			It("Should return a string based on these params", func() {
				os.Setenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY", "PRIVKEYTEST")
				Expect(util.HandlePrivateSSHKey(rawString)).To(Equal(expectedNotEmpty))
			})
		})
		Context("When rawString is empty and HUSKYCI_API_GIT_PRIVATE_SSH_KEY is not empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandlePrivateSSHKey("")).To(Equal(""))
			})
		})
		Context("When rawString is not empty and HUSKYCI_API_GIT_PRIVATE_SSH_KEY is empty", func() {
			It("Should return a string based on these params.", func() {
				os.Unsetenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
				Expect(util.HandlePrivateSSHKey(rawString)).To(Equal(expectedEmpty))
			})
		})
		Context("When rawString and HUSKYCI_API_GIT_PRIVATE_SSH_KEY are empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandlePrivateSSHKey("")).To(Equal(""))
			})
		})
	})

	Describe("GetLastLine", func() {

		rawString := `Warning: unpinned requirement
{"name":"enry", "vulnerability":"low"}`
		expected := `{"name":"enry", "vulnerability":"low"}`

		Context("When rawString is not empty", func() {
			It("Should return the string that is in the last position", func() {
				Expect(util.GetLastLine(rawString)).To(Equal(expected))
			})
		})
		Context("When rawString is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.GetLastLine("")).To(Equal(""))
			})
		})
	})

	Describe("GetAllLinesButLast", func() {

		rawString := `Line1
Line2
Line3
Line4`
		expected := []string{"Line1", "Line2", "Line3"}

		Context("When rawString is not empty", func() {
			It("Should return the slice of strings except the last line", func() {
				Expect(util.GetAllLinesButLast(rawString)).To(Equal(expected))
			})
		})
		Context("When rawString is empty", func() {
			It("Should return an empty slice of string.", func() {
				Expect(util.GetAllLinesButLast("")).To(Equal([]string{}))
			})
		})
	})

	Describe("RemoveDuplicates", func() {

		rawSliceString := []string{"item1", "item2", "item3", "item1", "item2"}
		expected := []string{"item1", "item2", "item3"}

		Context("When rawSliceString is not empty", func() {
			It("Should return the slice of strings except the last line", func() {
				Expect(util.RemoveDuplicates(rawSliceString)).To(Equal(expected))
			})
		})
		Context("When rawSliceString is empty", func() {
			It("Should return an empty slice of string.", func() {
				Expect(util.GetAllLinesButLast("")).To(Equal([]string{}))
			})
		})
	})

	Describe("CreateContainerName", func() {
		inputURL := "https://github.com/globocom/secDevLabs.git"
		inputBranch := "myBranch"
		inputImage := "secdevLabs/bandit"
		expected := "globocom_secDevLabs_myBranch_bandit"

		Context("When inputURL, imputBranch and inputImage are not empty", func() {
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
