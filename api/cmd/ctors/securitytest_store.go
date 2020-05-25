package ctors

import (
	"github.com/globocom/huskyCI/api/securitytest"
	"github.com/spf13/viper"
)

// NewSecurityTestStore create a new security test store containing all
// data from the security tests availables
func NewSecurityTestStore(setttings *viper.Viper) (securitytest.Store, error) {

	var newSecurityTestStore securitytest.Store

	if err := setttings.Unmarshal(&newSecurityTestStore); err != nil {
		return securitytest.Store{}, err
	}

	return newSecurityTestStore, nil
}
