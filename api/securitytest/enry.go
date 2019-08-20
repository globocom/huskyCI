// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/db"
	huskydocker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
)

// EnryScan holds all information needed for a gosec scan.
type EnryScan struct {
	RID         string
	CID         string
	URL         string
	Branch      string
	Image       string
	Command     string
	RawOutput   string
	FinalOutput EnryOutput
}

// EnryOutput is the struct that holds all data from Gosec output.
type EnryOutput struct {
	Codes []Code
}

// Code is the struct that stores all data from code found in a repository.
type Code struct {
	Language string   `bson:"language" json:"language"`
	Files    []string `bson:"files" json:"files"`
}

func newScanEnry(URL, branch, command string) EnryScan {
	return EnryScan{
		Image:   "huskyci/enry",
		URL:     URL,
		Branch:  branch,
		Command: util.HandleCmd(URL, branch, "", command),
	}
}

func newContainerEnry() (types.Container, error) {
	enryContainer := types.Container{}
	enryQuery := map[string]interface{}{"name": "enry"}
	enrySecurityTest, err := db.FindOneDBSecurityTest(enryQuery)
	if err != nil {
		return enryContainer, err
	}
	return types.Container{
		SecurityTest: enrySecurityTest,
		StartedAt:    time.Now(),
	}, nil
}

// RunScanEnry runs enry as the initial step of huskyCI and returns an error
func RunScanEnry(repository types.Repository) (EnryScan, error) {

	enryScan := EnryScan{}
	enryContainer, err := newContainerEnry()
	if err != nil {
		return enryScan, err
	}

	enryScan = newScanEnry(repository.URL, repository.Branch, enryContainer.SecurityTest.Cmd)
	if err := enryScan.startEnry(); err != nil {
		return enryScan, err
	}

	return enryScan, nil
}

// StartEnry starts a new enryScan and returns an error.
func (enryScan *EnryScan) startEnry() error {
	if err := enryScan.dockerRunEnry(); err != nil {
		return err
	}
	if err := enryScan.analyzeEnry(); err != nil {
		return err
	}
	// log.Info("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
	return nil
}

func (enryScan *EnryScan) dockerRunEnry() error {
	CID, cOutput, err := huskydocker.DockerRun(enryScan.Image, enryScan.Command)
	if err != nil {
		// log.Error("DockerRun", "DOCKERRUN", 3013, err)
		return err
	}
	enryScan.CID = CID
	enryScan.RawOutput = cOutput
	return nil
}

func (enryScan *EnryScan) analyzeEnry() error {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(enryScan.RawOutput, "ERROR_CLONING")
	if errorClonning {
		// log.Error("GosecStartAnalysis", "GOSEC", 1002, cOutput, err)
		return errors.New("error clonning")
	}

	// step 2: nil cOutput states that no Issues were found.
	if enryScan.RawOutput == "" {
		return errors.New("empty enry results")
	}

	// step 3: Get all languages and files found based on Enry output
	codesFound, err := prepareEnryOutput(enryScan.RawOutput)
	if err != nil {
		return err
	}

	enryScan.FinalOutput.Codes = codesFound
	return nil
}

func prepareEnryOutput(cOutput string) ([]Code, error) {
	repositoryLanguages := []Code{}
	newLanguage := Code{}
	mapLanguages := make(map[string][]interface{})
	err := json.Unmarshal([]byte(cOutput), &mapLanguages)
	if err != nil {
		// log.Error("EnryStartAnalysis", "ENRY", 1003, cOutput, err)
		return repositoryLanguages, err
	}
	for name, files := range mapLanguages {
		fs := []string{}
		for _, f := range files {
			if reflect.TypeOf(f).String() == "string" {
				fs = append(fs, f.(string))
			} else {
				// log.Error("getLanguagesFromEnryOutput", "ENRY", 1004, err)
				return repositoryLanguages, errors.New("error mapping languages")
			}
		}
		newLanguage = Code{
			Language: name,
			Files:    fs,
		}
		repositoryLanguages = append(repositoryLanguages, newLanguage)
	}
	return repositoryLanguages, nil
}
