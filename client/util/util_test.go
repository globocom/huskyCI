// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/globocom/huskyCI/client/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("CreateFile", func() {
		outputFileName := "sonarqube_test.json"
		outputFilePath := testOutputFilesPath + outputFileName
		testDataPath := "./testdata/sonarqube/sonarqube_test_example.json"
		fileString, err := ioutil.ReadFile(testDataPath)
		if err != nil {
			Fail(fmt.Sprintf("error trying to read fixture file: %s", err.Error()))
		}
		err = util.CreateFile(fileString, testOutputFilesPath, outputFileName)
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

	Describe("GetAllLinesButLast", func() {
		rawString := strings.Join([]string{"A", "B", "C", "D"}, "\n")
		expected := []string{"A", "B", "C"}

		Context("When rawString is not empty", func() {
			It("Should return the slice of strings except the last line", func() {
				Expect(util.GetAllLinesButLast(rawString)).To(Equal(expected))
			})
		})
	})

	Describe("GetLastLine", func() {
		rawString := strings.Join([]string{"A", "B", "C", "D"}, "\n")
		expected := "D"

		Context("When rawString is not empty", func() {
			It("Should return the string that is in the last position", func() {
				Expect(util.GetLastLine(rawString)).To(Equal(expected))
			})
		})
	})
})
