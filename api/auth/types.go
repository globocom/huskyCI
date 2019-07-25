package auth

type BasicVerification interface {
	IsValidUser(username, password string) (bool, error)
}

type BasicAuthentication struct {
	BasicClient BasicVerification
}

type MongoBasic struct{}

type FakeClient struct {
	ExpectedIsValidBool  bool
	ExpectedIsValidError error
}

func (fC *FakeClient) IsValidUser(username, password string) (bool, error) {
	return fC.ExpectedIsValidBool, fC.ExpectedIsValidError
}
