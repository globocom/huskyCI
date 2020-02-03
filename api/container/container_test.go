// package container_test

// import "os"

// var _ = Describe("Container", func() {

// 	Describe("HandleCmd", func() {
// 		inputRepositoryURL := "https://github.com/globocom/secDevLabs.git"
// 		inputRepositoryBranch := "myBranch"
// 		inputCMD := "git clone -b %GIT_BRANCH% --single-branch %GIT_REPO% code --quiet 2> /tmp/errorGitClone -- "
// 		expected := "git clone -b myBranch --single-branch https://github.com/globocom/secDevLabs.git code --quiet 2> /tmp/errorGitClone -- "

// 		Context("When inputRepositoryURL, inputRepositoryBranch and inputCMD are not empty", func() {
// 			It("Should return a string based on these params", func() {
// 				Expect(HandleCmd(inputRepositoryURL, inputRepositoryBranch, inputCMD)).To(Equal(expected))
// 			})
// 		})
// 		Context("When inputRepositoryURL is empty", func() {
// 			It("Should return an empty string.", func() {
// 				Expect(HandleCmd("", inputRepositoryBranch, inputCMD)).To(Equal(""))
// 			})
// 		})
// 		Context("When inputRepositoryBranch is empty", func() {
// 			It("Should return an empty string.", func() {
// 				Expect(HandleCmd(inputRepositoryURL, "", inputCMD)).To(Equal(""))
// 			})
// 		})
// 		Context("When inputCMD is empty", func() {
// 			It("Should return an empty string.", func() {
// 				Expect(HandleCmd(inputRepositoryURL, inputRepositoryBranch, "")).To(Equal(""))
// 			})
// 		})
// 	})

// 	Describe("HandlePrivateSSHKey", func() {

// 		rawString := "echo 'GIT_PRIVATE_SSH_KEY' > ~/.ssh/huskyci_id_rsa &&"
// 		expectedNotEmpty := "echo 'PRIVKEYTEST' > ~/.ssh/huskyci_id_rsa &&"
// 		expectedEmpty := "echo '' > ~/.ssh/huskyci_id_rsa &&"

// 		Context("When rawString and HUSKYCI_API_GIT_PRIVATE_SSH_KEY are not empty", func() {
// 			It("Should return a string based on these params", func() {
// 				os.Setenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY", "PRIVKEYTEST")
// 				Expect(HandlePrivateSSHKey(rawString)).To(Equal(expectedNotEmpty))
// 			})
// 		})
// 		Context("When rawString is empty and HUSKYCI_API_GIT_PRIVATE_SSH_KEY is not empty", func() {
// 			It("Should return an empty string.", func() {
// 				Expect(HandlePrivateSSHKey("")).To(Equal(""))
// 			})
// 		})
// 		Context("When rawString is not empty and HUSKYCI_API_GIT_PRIVATE_SSH_KEY is empty", func() {
// 			It("Should return a string based on these params.", func() {
// 				os.Unsetenv("HUSKYCI_API_GIT_PRIVATE_SSH_KEY")
// 				Expect(HandlePrivateSSHKey(rawString)).To(Equal(expectedEmpty))
// 			})
// 		})
// 		Context("When rawString and HUSKYCI_API_GIT_PRIVATE_SSH_KEY are empty", func() {
// 			It("Should return an empty string.", func() {
// 				Expect(HandlePrivateSSHKey("")).To(Equal(""))
// 			})
// 		})
// 	})
// })
