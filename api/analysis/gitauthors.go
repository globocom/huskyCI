package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
)

// GitAuthorsOutput is the struct that holds all commit authors from a branch.
type GitAuthorsOutput struct {
	Authors []string `json:"authors"`
}

// GitAuthorsCheckOutputFlow analyses the output from Gitauthors and sets a cResult based on it.
func GitAuthorsCheckOutputFlow(CID string, cOutput string, RID string) {

	// step 1: check for any errors when clonning repo
	errorClonning := strings.Contains(cOutput, "ERROR_CLONING")
	if errorClonning {
		if err := updateInfoAndResultBasedOnCID("Error clonning repository", "error", CID); err != nil {
			return
		}
		return
	}

	// step 2: Unmarshall cOutput into GitAuthorsOutput struct.
	gitAuthorsOutput := GitAuthorsOutput{}
	err := json.Unmarshal([]byte(cOutput), &gitAuthorsOutput)
	if err != nil {
		log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 1002, cOutput, err)
		return
	}

	// step 3: check if authors is empty (master branch was probably sent)
	if len(gitAuthorsOutput.Authors) == 0 {
		if err := updateInfoAndResultBasedOnCID("Could not get authors. Probably master branch is being analyzed.", "warning", CID); err != nil {
			return
		}
		return
	}

	// step 4: update analysis with the commit authors found
	if err := updateCommitAuthorsBasedOnRID(gitAuthorsOutput.Authors, RID); err != nil {
		return
	}
}

func updateCommitAuthorsBasedOnRID(commitAuthors []string, RID string) error {

	analysisQuery := map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("updateCommitAuthorsBasedOnRID", "GITAUTHORS", 2008, RID, err)
		return err
	}

	analysis.CommitAuthors = commitAuthors
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("updateCommitAuthorsBasedOnRID", "GITAUTHORS", 2007, err)
		return err
	}

	return nil
}
