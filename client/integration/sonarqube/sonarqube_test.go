// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sonarqube_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/globocom/huskyCI/client/integration/sonarqube"
	"github.com/globocom/huskyCI/client/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("SonarQube", func() {
	Describe("GenerateOutputFile", func() {
		DescribeTable("Analysis containing vulnerabilities",
			func(inputAnalysisFile, outputPath, outputFileName, expectedOutputFile string) {
				analysisFilePath := analysisTestDataPath + inputAnalysisFile
				analysisFileString, err := ioutil.ReadFile(analysisFilePath)
				if err != nil {
					Fail(fmt.Sprintf("error trying to read fixture file %s: %s", analysisFilePath, err.Error()))
				}

				analysis := types.Analysis{}
				err = json.Unmarshal([]byte(analysisFileString), &analysis)
				if err != nil {
					Fail(fmt.Sprintf("error trying to unmarshal fixture file %s: %s", analysisFilePath, err.Error()))
				}

				err = sonarqube.GenerateOutputFile(analysis, outputPath, outputFileName)
				Expect(err).NotTo(HaveOccurred())

				expectedOutputFilePath := sonarqubeTestDataPath + expectedOutputFile
				expectedFileString, err := ioutil.ReadFile(expectedOutputFilePath)
				if err != nil {
					Fail(fmt.Sprintf("error trying to read fixture file %s: %s", expectedOutputFilePath, err.Error()))
				}

				testOutputFilePath := outputPath + outputFileName
				testFileString, err := ioutil.ReadFile(testOutputFilePath)
				Expect(err).NotTo(HaveOccurred())

				Expect(testFileString).To(Equal(expectedFileString))
			},
			Entry("Vulnerable go project", "vulnerable_go_project.json", testOutputFilesPath, "sonarqube_go_test.json", "vulnerable_go_output.json"),
			Entry("Vulnerable python project", "vulnerable_python_project.json", testOutputFilesPath, "sonarqube_python_test.json", "vulnerable_python_output.json"),
			Entry("Vulnerable ruby project", "vulnerable_ruby_project.json", testOutputFilesPath, "sonarqube_ruby_test.json", "vulnerable_ruby_output.json"),
			Entry("Vulnerable js project", "vulnerable_js_project.json", testOutputFilesPath, "sonarqube_js_test.json", "vulnerable_js_output.json"),
			Entry("Not Vulnerable project", "not_vulnerable_project.json", testOutputFilesPath, "sonarqube_not_vulnerable_test.json", "not_vulnerable_output.json"),
		)
	})
})
