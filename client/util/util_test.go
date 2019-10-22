// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/globocom/huskyCI/client/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("CreateFile", func() {
		Context("When analyzing an analysis with vulnerabilities", func() {
			outputFileName := "sonarqube_test.json"
			outputFilePath := fmt.Sprintf("./huskyCI/%s", outputFileName)
			filePath := "testdata/analysis/sonarqube_test_example.json"
			fileString, err := ioutil.ReadFile(filePath)
			if err != nil {
				Fail(fmt.Sprintf("error trying to read fixture file: %s", err.Error()))
			}
			bytesInput := []byte(fileString)
			err = util.CreateFile(bytesInput, outputFileName)
			if err != nil {
				Fail(fmt.Sprintf("eror trying to execute util.CreateFile: %s", err.Error()))
			}
			It("should not return error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("Should create a directory and file", func() {
				_, err := os.Stat(outputFilePath)
				Expect(!os.IsNotExist(err)).To(Equal(true))
			})
			It("File content should match the input string", func() {
				outputString, err := ioutil.ReadFile(outputFilePath)
				if err != nil {
					Fail(fmt.Sprintf("error trying to read test output file: %s", err.Error()))
				}
				Expect(outputString).To(Equal(fileString))
			})
		})
	})
})
