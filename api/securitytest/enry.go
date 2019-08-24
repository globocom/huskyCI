// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/globocom/huskyCI/api/log"
)

// EnryOutput is the struct that holds all data from Gosec output.
type EnryOutput struct {
	Codes []Code
}

// Code is the struct that stores all data from code found in a repository.
type Code struct {
	Language string   `bson:"language" json:"language"`
	Files    []string `bson:"files" json:"files"`
}

func analyzeEnry(enryScan *SecTestScanInfo) error {
	// Unmarshall rawOutput into finalOutput, that is a EnryOutput struct.
	if err := json.Unmarshal([]byte(enryScan.Container.COutput), &enryScan.FinalOutput); err != nil {
		log.Error("analyzeEnry", "ENRY", 1002, enryScan.Container.COutput, err)
		enryScan.ErrorFound = err
		return err
	}
	// get all languages and files found based on Enry output
	if err := enryScan.prepareEnryOutput(); err != nil {
		enryScan.ErrorFound = err
		return err
	}
	return nil
}

func (enryScan *SecTestScanInfo) prepareEnryOutput() error {
	repositoryLanguages := []Code{}
	newLanguage := Code{}
	mapLanguages := make(map[string][]interface{})
	err := json.Unmarshal([]byte(enryScan.Container.COutput), &mapLanguages)
	if err != nil {
		log.Error("prepareEnryOutput", "ENRY", 1003, enryScan.Container.COutput, err)
		return err
	}
	for name, files := range mapLanguages {
		fs := []string{}
		for _, f := range files {
			if reflect.TypeOf(f).String() == "string" {
				fs = append(fs, f.(string))
			} else {
				errMsg := errors.New("error mapping languages")
				log.Error("prepareEnryOutput", "ENRY", 1032, errMsg)
				return errMsg
			}
		}
		newLanguage = Code{
			Language: name,
			Files:    fs,
		}
		repositoryLanguages = append(repositoryLanguages, newLanguage)
	}
	enryScan.Codes = repositoryLanguages
	return nil
}
