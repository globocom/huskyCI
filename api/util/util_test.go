// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {

	Describe("HandleCmd", func() {
		inputRepositoryURL := "https://github.com/globocom/secDevLabs.git"
		inputRepositoryBranch := "myBranch"
		inputCMD := "git clone -b %GIT_BRANCH% --single-branch %GIT_REPO% code --quiet 2> /tmp/errorGitClone -- "
		expected := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitClone -- "

		Context("When inputRepositoryURL, inputRepositoryBranch and inputCMD are not empty", func() {
			It("Should return a string based on these params", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, inputCMD)).To(Equal(expected))
			})
		})
		Context("When inputRepositoryURL is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd("", inputRepositoryBranch, inputCMD)).To(Equal(""))
			})
		})
		Context("When inputRepositoryBranch is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd(inputRepositoryURL, "", inputCMD)).To(Equal(""))
			})
		})
		Context("When inputCMD is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.HandleCmd(inputRepositoryURL, inputRepositoryBranch, "")).To(Equal(""))
			})
		})
	})

	Describe("HandleGitURLSubstitution", func() {

		rawString := "git config --global url.\"%GIT_SSH_URL%:\".insteadOf \"%GIT_URL_TO_SUBSTITUTE%\""
		expectedURLToSubstituteNotEmpty := "git config --global url.\"nil:\".insteadOf \"nil\""
		expectedSSHURLNotEmpty := "git config --global url.\"nil:\".insteadOf \"nil\""
		expectedBothVarsEmpty := "git config --global url.\"nil:\".insteadOf \"nil\""
		expectedNotEmpty := "git config --global url.\"gitlab@gitlab.example.com:\".insteadOf \"https://gitlab.example.com/\""

		Context("When rawString is not empty, HUSKYCI_API_GIT_SSH_URL is empty, but HUSKYCI_API_GIT_URL_TO_SUBSTITUTE is not empty", func() {
			It("Should return a string based on these params", func() {
				os.Setenv("HUSKYCI_API_GIT_SSH_URL", "")
				os.Setenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE", "https://gitlab.example.com/")
				Expect(util.HandleGitURLSubstitution(rawString)).To(Equal(expectedURLToSubstituteNotEmpty))
			})
		})
		Context("When rawString is not empty, HUSKYCI_API_GIT_SSH_URL is not empty and HUSKYCI_API_GIT_URL_TO_SUBSTITUTE is empty", func() {
			It("Should return a string based on these params", func() {
				os.Setenv("HUSKYCI_API_GIT_SSH_URL", "gitlab@gitlab.example.com")
				os.Setenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE", "")
				Expect(util.HandleGitURLSubstitution(rawString)).To(Equal(expectedSSHURLNotEmpty))
			})
		})
		Context("When rawString is not empty, HUSKYCI_API_GIT_SSH_URL is empty and HUSKYCI_API_GIT_URL_TO_SUBSTITUTE is empty", func() {
			It("Should return a string based on these params", func() {
				os.Setenv("HUSKYCI_API_GIT_SSH_URL", "")
				os.Setenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE", "")
				Expect(util.HandleGitURLSubstitution(rawString)).To(Equal(expectedBothVarsEmpty))
			})
		})
		Context("When rawString is not empty, HUSKYCI_API_GIT_SSH_URL is not empty and HUSKYCI_API_GIT_URL_TO_SUBSTITUTE is not empty", func() {
			It("Should return a string based on these params", func() {
				os.Setenv("HUSKYCI_API_GIT_SSH_URL", "gitlab@gitlab.example.com")
				os.Setenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE", "https://gitlab.example.com/")
				Expect(util.HandleGitURLSubstitution(rawString)).To(Equal(expectedNotEmpty))
			})
		})
		Context("When rawString is empty", func() {
			It("Should return an empty string", func() {
				os.Setenv("HUSKYCI_API_GIT_SSH_URL", "gitlab@gitlab.example.com")
				os.Setenv("HUSKYCI_API_GIT_URL_TO_SUBSTITUTE", "https://gitlab.example.com/")
				Expect(util.HandleGitURLSubstitution("")).To(Equal(""))
			})
		})
	})

	Describe("HandlePrivateSSHKey", func() {

		rawString := "echo '%GIT_PRIVATE_SSH_KEY%' > ~/.ssh/huskyci_id_rsa &&"
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
			It("Should return slice of non-duplicate elements", func() {
				Expect(util.RemoveDuplicates(rawSliceString)).To(Equal(expected))
			})
		})
		Context("When rawSliceString is empty", func() {
			It("Should return an empty slice of string.", func() {
				Expect(util.GetAllLinesButLast("")).To(Equal([]string{}))
			})
		})
	})

	Describe("CheckMaliciousRID", func() {
		e := echo.New()

		Context("When RID is valid", func() {
			It("Should pass with no error", func() {
				w := httptest.NewRecorder()
				c := e.NewContext(httptest.NewRequest(http.MethodGet, "/foo", nil), w)
				Expect(util.CheckMaliciousRID("*", c)).To(BeNil())
			})
		})
		Context("When RID is invalid", func() {
			It("Should response with invalid RID", func() {
				w := httptest.NewRecorder()
				c := e.NewContext(httptest.NewRequest(http.MethodGet, "/foo", nil), w)
				Expect(util.CheckMaliciousRID("*", c)).To(BeNil())

				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(ioutil.ReadAll(resp.Body)).To(
					MatchJSON(`{"success": false, "error": "invalid RID"}`),
				)
			})
		})
	})

	Describe("CheckValidInput", func() {
		e := echo.New()
		log.InitLog(true, "", "", "log_test", "log_test")

		Context("When URL is already ok", func() {
			repository := types.Repository{
				URL:    "https://github.com/globocom/secDevLabs.git",
				Branch: "branch",
			}

			c := e.NewContext(nil, nil)
			It("Should return the same URL", func() {
				Expect(util.CheckValidInput(repository, c)).To(Equal(repository.URL))
			})
		})

		Context("When URL is invalid", func() {
			It("Should response with invalid repository URL", func() {
				repository := types.Repository{
					URL:    "http://globo.com",
					Branch: "branch",
				}

				w := httptest.NewRecorder()
				c := e.NewContext(httptest.NewRequest(http.MethodGet, "/foo", nil), w)

				Expect(util.CheckValidInput(repository, c)).To(HaveLen(0))

				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(ioutil.ReadAll(resp.Body)).To(
					MatchJSON(`{"success": false, "error": "invalid repository URL"}`),
				)
			})

			It("Should response with invalid branch", func() {
				repository := types.Repository{
					URL:    "https://github.com/globocom/secDevLabs.git",
					Branch: " [bra nch] ",
				}

				w := httptest.NewRecorder()
				c := e.NewContext(httptest.NewRequest(http.MethodGet, "/foo", nil), w)

				Expect(util.CheckValidInput(repository, c)).To(Equal(repository.URL))

				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(ioutil.ReadAll(resp.Body)).To(
					MatchJSON(`{"success": false, "error": "invalid repository branch"}`),
				)
			})
		})
	})

	Describe("SanitizeSafetyJSON", func() {

		rawSliceString := `{"result":"This vulnerability was found \\ and should be replaced.}`
		expected := `{"result":"This vulnerability was found \\\\ and should be replaced.}`

		Context("When rawSliceString is not empty", func() {
			It("Should return the string expected.", func() {
				Expect(util.SanitizeSafetyJSON(rawSliceString)).To(Equal(expected))
			})
		})
		Context("When rawSliceString is empty", func() {
			It("Should return an empty string.", func() {
				Expect(util.SanitizeSafetyJSON("")).To(Equal(""))
			})
		})
	})

	Describe("CountDigits", func() {

		rawSliceInteger := []int{-1, 0, 10}
		expected := []int{0, 1, 2}

		Context("When rawSliceInteger is greater than zero", func() {
			It("Should return the expected integer.", func() {
				Expect(util.CountDigits(rawSliceInteger[2])).To(Equal(expected[2]))
			})
		})
		Context("When rawSliceInteger is less than zero", func() {
			It("Should return the expected integer.", func() {
				Expect(util.CountDigits(rawSliceInteger[0])).To(Equal(expected[1]))
			})
		})
		Context("When rawSliceInteger is zero", func() {
			It("Should return the expected integer.", func() {
				Expect(util.CountDigits(rawSliceInteger[1])).To(Equal(expected[0]))
			})
		})
	})

	Describe("VerifyNoHusky", func() {

		rawBanditCodeSliceString := []string{"1 secret = 'password123!'\n2 password = 'thisisnotapassword' #nohusky"}
		rawLineNumberSliceInteger := []int{1, 2}
		rawSecurityToolSliceString := []string{"Bandit"}

		Context("Bandit: When line number doesn't match the one in the code string", func() {
			It("Should return false.", func() {
				Expect(util.VerifyNoHusky(rawBanditCodeSliceString[0], rawLineNumberSliceInteger[0], rawSecurityToolSliceString[0])).To(BeFalse())
			})
		})

		Context("Bandit: When line number matches the one in the code string", func() {
			It("Should return true.", func() {
				Expect(util.VerifyNoHusky(rawBanditCodeSliceString[0], rawLineNumberSliceInteger[1], rawSecurityToolSliceString[0])).To(BeTrue())
			})
		})

		Context("Bandit: When line number doesn't match the one in the code string", func() {
			It("Should return false.", func() {
				Expect(util.VerifyNoHusky(rawBanditCodeSliceString[0], rawLineNumberSliceInteger[0], rawSecurityToolSliceString[0])).To(BeFalse())
			})
		})
	})
})
