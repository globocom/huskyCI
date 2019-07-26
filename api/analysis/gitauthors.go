package analysis

import (
	"encoding/json"
	"strings"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"gopkg.in/mgo.v2/bson"
)

// GitAuthorsOutput is the struct that holds all commit authors from a branch.
type GitAuthorsOutput struct {
	Authors []string `json:"authors"`
}

// GitAuthorsCheckOutputFlow analyses the output from Gitauthors and sets a cResult based on it.
func GitAuthorsCheckOutputFlow(CID string, cOutput string, RID string) {

	analysisQuery := map[string]interface{}{"containers.CID": CID}

	// check if there were errors when clonning repository
	if strings.Contains(cOutput, "ERROR_CLONING") {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "error",
				"containers.$.cInfo":   "Error clonning repository.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 2007, "Step 0.2 ", err)
		}
		return
	}

	// Unmarshall cOutput into GitAuthorsOutput struct.
	gitAuthorsOutput := GitAuthorsOutput{}
	err := json.Unmarshal([]byte(cOutput), &gitAuthorsOutput)
	if err != nil {
		log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 1002, cOutput, err)
		return
	}

	// check if authors is empty (master branch was probably sent)
	if len(gitAuthorsOutput.Authors) == 0 {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cResult": "warning",
				"containers.$.cInfo":   "Could not get authors. Probably master branch is being analyzed.",
			},
		}
		err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 2007, "Step 0.2 ", err)
		}
		return
	}

	// get updated analysis based on its RID
	analysisQuery = map[string]interface{}{"RID": RID}
	analysis, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 2008, CID, err)
		return
	}

	analysis.CommitAuthors = gitAuthorsOutput.Authors
	err = db.UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		log.Error("GitAuthorsStartAnalysis", "GITAUTHORS", 2007, err)
		return
	}

}
