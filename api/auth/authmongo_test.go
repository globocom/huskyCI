package auth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	. "github.com/globocom/huskyCI/api/auth"
	"hash"
	"os"
)

type FakeGen struct {
	expectedHash string
}

func (fG *FakeGen) GenHashValue(value, salt []byte, iter, keyLen int, h hash.Hash) []byte {
	return []byte(fG.expectedHash)
}

var _ = Describe("Authmongo", func() {
	Context("When env var is nil", func() {
		It("Should return an empty string with the expected error", func() {
			fakeGen := FakeGen{
				expectedHash: "nothing",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Salt:         "mystrongsalt",
				Iterations:   1,
				KeyLen:       12,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal(""))
			Expect(err).To(Equal(errors.New("Failed to generate a hash! It doesn't meet all criteria")))
		})
	})
	Context("When hash algorithm chosen is not valid", func() {
		It("Should return the expected error and a nil string", func() {
			os.Setenv("HUSKY_PBKDF2_VALID_HASHES", "sha512")
			fakeGen := FakeGen{
				expectedHash: "nothing",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Salt:         "mystrongsalt",
				Iterations:   1,
				KeyLen:       12,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal(""))
			Expect(err).To(Equal(errors.New("Failed to generate a hash! It doesn't meet all criteria")))
		})
	})
	Context("When hash algorithm chosen is valid", func() {
		It("Should return an nil error and the expected string", func() {
			os.Setenv("HUSKY_PBKDF2_VALID_HASHES", "SHA512,SHA256,MD5")
			fakeGen := FakeGen{
				expectedHash: "MyHashedString",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Salt:         "mystrongsalt",
				Iterations:   1,
				KeyLen:       12,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal("MyHashedString"))
			Expect(err).To(BeNil())
		})
	})
	Context("When one of the required fields for PBKDF2 is not valid", func() {
		It("Should return an the expected error and an empty hash for an empty salt", func() {
			fakeGen := FakeGen{
				expectedHash: "MyHashedString",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Iterations:   1,
				KeyLen:       12,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal(""))
			Expect(err).To(Equal(errors.New("Failed to generate a hash! It doesn't meet all criteria")))
		})
		It("Should return an the expected error and an empty hash for a 0 iteration", func() {
			fakeGen := FakeGen{
				expectedHash: "MyHashedString",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Salt:         "ValidSalt",
				KeyLen:       12,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal(""))
			Expect(err).To(Equal(errors.New("Failed to generate a hash! It doesn't meet all criteria")))
		})
		It("Should return an the expected error and an empty hash for a 0 keyLength", func() {
			fakeGen := FakeGen{
				expectedHash: "MyHashedString",
			}
			pbkdf2Client := ClientPbkdf2{
				HashGen:      &fakeGen,
				Salt:         "ValidSalt",
				Iterations:   1,
				KeyLen:       0,
				HashFunction: "sha256",
			}
			hashVal, err := pbkdf2Client.GetHashedPass("mypass")
			Expect(hashVal).To(Equal(""))
			Expect(err).To(Equal(errors.New("Failed to generate a hash! It doesn't meet all criteria")))
		})
	})
})
