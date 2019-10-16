package util_test

import (
	"testing"

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
	log.InitLog(true, "", "", "log_test", "log_test")

	if log.Logger == nil {
		t.Error("expected logger to be initialized, but it wasn't")
		return
	}
}
