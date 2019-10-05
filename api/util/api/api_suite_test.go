package util_test

import (
	"testing"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	TestInitLog(t)
	RunSpecs(t, "Api Suite")
}

func TestInitLog(t *testing.T) {
	apiContext.APIConfiguration = &apiContext.APIConfig{
		GraylogConfig: &apiContext.GraylogConfig{
			DevelopmentEnv: true,
			AppName:        "log_test",
			Tag:            "log_test",
		},
	}

	log.InitLog()

	if log.Logger == nil {
		t.Error("expected logger to be initialized, but it wasn't")
		return
	}
}
