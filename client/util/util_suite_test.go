// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Suite")
}

var _ = AfterSuite(func() {
	cleanTestOutputFiles()
})

const testOutputFilesPath = "./huskyCITest/"

func cleanTestOutputFiles() {
	os.RemoveAll(testOutputFilesPath)
}
