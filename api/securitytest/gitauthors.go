// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package securitytest

import (
	"encoding/json"

	"github.com/globocom/huskyCI/api/log"
)

// GitAuthorsOutput is the struct that holds all commit authors from a branch.
type GitAuthorsOutput struct {
	Authors []string `json:"authors"`
}

func analyzeGitAuthors(gitAuthorsScan *SecTestScanInfo) error {

	gitAuthorsOutput := GitAuthorsOutput{}
	gitAuthorsScan.FinalOutput = gitAuthorsOutput

	// Unmarshall rawOutput into finalOutput, that is a GitAuthors struct.
	if err := json.Unmarshal([]byte(gitAuthorsScan.Container.COutput), &gitAuthorsOutput); err != nil {
		log.Error("analyzeGitAuthors", "GITAUTHORS", 1035, gitAuthorsScan.Container.COutput, err)
		gitAuthorsScan.ErrorFound = err
		gitAuthorsScan.prepareContainerAfterScan()
		return err
	}
	gitAuthorsScan.FinalOutput = gitAuthorsOutput

	// check if authors is empty (master branch was probably sent)
	if len(gitAuthorsOutput.Authors) == 0 {
		gitAuthorsScan.CommitAuthorsNotFound = true
		gitAuthorsScan.prepareContainerAfterScan()
	}

	gitAuthorsScan.CommitAuthors = gitAuthorsOutput
	gitAuthorsScan.prepareContainerAfterScan()
	return nil
}
