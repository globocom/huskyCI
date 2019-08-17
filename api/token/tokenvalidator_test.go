package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	. "github.com/globocom/huskyCI/api/token"
	"github.com/globocom/huskyCI/api/types"
)

type FakeVerifier struct {
	expectedValidateError error
	expectedVerifyError   error
}

func (fV *FakeVerifier) GenerateAccessToken(repo types.TokenRequest) (string, error) {
	return "", nil
}

func (fV *FakeVerifier) ValidateToken(token, repositoryURL string) error {
	return fV.expectedValidateError
}

func (fV *FakeVerifier) VerifyRepo(repositoryURL string) error {
	return fV.expectedVerifyError
}

var _ = Describe("Tokenvalidator", func() {
	Describe("HasAuthorization", func() {
		Context("When VerifyRepo returns an error", func() {
			It("Should return a true boolean", func() {
				fakeVerifier := FakeVerifier{
					expectedVerifyError: errors.New("Could not find the repository URL"),
				}
				validator := TokenValidator{
					TokenVerifier: &fakeVerifier,
				}
				Expect(validator.HasAuthorization("MyToken", "MyRepo")).To(BeTrue())
			})
		})
		Context("When ValidateToken returns an error", func() {
			It("Should return a false boolean", func() {
				FakeVerifier := FakeVerifier{
					expectedVerifyError:   nil,
					expectedValidateError: errors.New("Token is not valid"),
				}
				validator := TokenValidator{
					TokenVerifier: &FakeVerifier,
				}
				Expect(validator.HasAuthorization("MyToken", "MyRepo")).To(BeFalse())
			})
		})
		Context("When ValidateToken returns a nil error", func() {
			It("Should return a true boolean", func() {
				FakeVerifier := FakeVerifier{
					expectedVerifyError:   nil,
					expectedValidateError: nil,
				}
				validator := TokenValidator{
					TokenVerifier: &FakeVerifier,
				}
				Expect(validator.HasAuthorization("MyToken", "MyRepo")).To(BeTrue())
			})
		})
	})
})
