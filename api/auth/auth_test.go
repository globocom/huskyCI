package auth_test

import (
	"errors"
	"github.com/globocom/huskyCI/api/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateUser", func() {
	Context("When IsValidUser returns an error", func() {
		It("Should return the same error with false bool", func() {
			fakeClient := auth.FakeClient{
				ExpectedIsValidBool:  false,
				ExpectedIsValidError: errors.New("Failed to get user's info"),
			}
			userAuth := auth.BasicAuthentication{
				BasicClient: &fakeClient,
			}
			isValid, err := userAuth.ValidateUser("husky", "dumbpass", nil)
			Expect(isValid).To(BeFalse())
			Expect(err).To(Equal(errors.New("Failed to get user's info")))
		})
	})
	Context("When IsValidUser return an invalid authentication", func() {
		It("Should return a false boolean with a nil error", func() {
			fakeClient := auth.FakeClient{
				ExpectedIsValidBool:  false,
				ExpectedIsValidError: nil,
			}
			userAuth := auth.BasicAuthentication{
				BasicClient: &fakeClient,
			}
			isValid, err := userAuth.ValidateUser("husky", "dumbpass", nil)
			Expect(isValid).To(BeFalse())
			Expect(err).To(BeNil())
		})
	})
	Context("When IsValidUser return a valid authentication", func() {
		It("Should return a true boolean with a nil error", func() {
			fakeClient := auth.FakeClient{
				ExpectedIsValidBool:  true,
				ExpectedIsValidError: nil,
			}
			userAuth := auth.BasicAuthentication{
				BasicClient: &fakeClient,
			}
			isValid, err := userAuth.ValidateUser("husky", "dumbpass", nil)
			Expect(isValid).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})

})
