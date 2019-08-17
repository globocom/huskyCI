// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth_test

import (
	"errors"

	. "github.com/globocom/huskyCI/api/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeClient struct {
	expectedPass           string
	expectedGetPassError   error
	expectedHashedPass     string
	expectedGetHashedError error
}

func (fC *FakeClient) GetPassFromDB(username string) (string, error) {
	return fC.expectedPass, fC.expectedGetPassError
}

func (fC *FakeClient) GetHashedPass(password string) (string, error) {
	return fC.expectedHashedPass, fC.expectedGetHashedError
}

var _ = Describe("IsValidUser", func() {
	Context("When GetPassFromDB returns an error", func() {
		It("Should return a false boolean and the same error", func() {
			fakeHandler := FakeClient{
				expectedPass:         "",
				expectedGetPassError: errors.New("Error trying to get password from DB: User not found"),
			}
			mongoClient := MongoBasic{
				ClientHandler: &fakeHandler,
			}
			isValid, err := mongoClient.IsValidUser("husky", "somepass")
			Expect(isValid).To(BeFalse())
			Expect(err).To(BeNil())
		})
	})
	Context("When GetHashedPass returns an error", func() {
		It("Should return a false boolean and the same error", func() {
			fakeHandler := FakeClient{
				expectedPass:           "hashedpass",
				expectedGetPassError:   nil,
				expectedGetHashedError: errors.New("Error trying to generate hash from password"),
				expectedHashedPass:     "",
			}
			mongoClient := MongoBasic{
				ClientHandler: &fakeHandler,
			}
			isValid, err := mongoClient.IsValidUser("husky", "somepass")
			Expect(isValid).To(BeFalse())
			Expect(err).To(Equal(errors.New("Error trying to generate hash from password")))
		})
	})
	Context("When passDB and hashedPass are different", func() {
		It("Should return a false boolean with nil error", func() {
			fakeHandler := FakeClient{
				expectedPass:           "hashedpass",
				expectedGetPassError:   nil,
				expectedGetHashedError: nil,
				expectedHashedPass:     "differenthash",
			}
			mongoClient := MongoBasic{
				ClientHandler: &fakeHandler,
			}
			isValid, err := mongoClient.IsValidUser("husky", "somepass")
			Expect(isValid).To(BeFalse())
			Expect(err).To(BeNil())
		})
	})
	Context("When passDB and hashedPass are equal", func() {
		It("Should return a true boolean with nil error", func() {
			fakeHandler := FakeClient{
				expectedPass:           "hashedpass",
				expectedGetPassError:   nil,
				expectedGetHashedError: nil,
				expectedHashedPass:     "hashedpass",
			}
			mongoClient := MongoBasic{
				ClientHandler: &fakeHandler,
			}
			isValid, err := mongoClient.IsValidUser("husky", "somepass")
			Expect(isValid).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})
})
