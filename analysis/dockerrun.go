package analysis

import (
	"fmt"
	"regexp"
	"time"

	docker "github.com/globocom/husky/dockers"
	"github.com/globocom/husky/types"
	"gopkg.in/mgo.v2/bson"
)

// DockerRun starts a new container, runs a given securityTest in it and then updates AnalysisCollection.
func DockerRun(RID string, analysis *types.Analysis, securityTest types.SecurityTest) {

	// step 0: add a new container to the analysis.
	newContainer := types.Container{SecurityTest: securityTest}
	analysisQuery := map[string]interface{}{"RID": RID}
	startedAt := time.Now()

	// step 1: create a new container.
	d := docker.Docker{}
	CID, err := d.CreateContainer(*analysis, securityTest.Image, securityTest.Cmd)
	if err != nil {
		// error creating container. maxRetry?
		newContainer.CStatus = "error"
		analysis.Containers = append(analysis.Containers, newContainer)
		err := UpdateOneDBAnalysis(analysisQuery, *analysis)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (step 1-err):", err)
		}
	} else {
		newContainer.CID = CID
		newContainer.CStatus = "created"
		analysis.Containers = append(analysis.Containers, newContainer)
		err := UpdateOneDBAnalysis(analysisQuery, *analysis)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (step 1):", err)
		}
		analysisQuery = map[string]interface{}{"containers.CID": CID}
	}

	// step 2: start created container.
	err = d.StartContainer(CID)
	if err != nil {
		// error starting container. maxRetry?
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cStatus": "error",
			},
		}
		err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (step 2-err):", err)
		}
	} else {
		startedAt = time.Now()
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cStatus":   "running",
				"containers.$.startedAt": startedAt,
			},
		}
		err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (step 2):", err)
		}
	}

	// step 3: wait container finish running.
	err = d.WaitContainer(CID)
	if err != nil {
		fmt.Println("Error waiting the container", CID, ":", err)
	}

	// step 4: read cmd output from container.
	var cleanedOutput string
	newContainer.COuput, err = d.ReadOutput(CID)
	if err != nil {
		// error reading container's output. maxRetry?
		fmt.Println("Error reading output from container", CID, ":", err)
	} else {
		finishedAt := time.Now()
		// cleaning json output from dockerfile logs.
		reg, err := regexp.Compile(`[{\[]{1}([,:{}\[\]0-9.\-+Eaeflnr-u \n\r\t]|".*?")+[}\]]{1}`)
		if err != nil {
			fmt.Println("Error regexp:", err)
			return
		}
		cleanedOutput = reg.FindString(newContainer.COuput)
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cStatus":    "finished",
				"containers.$.finishedAt": finishedAt,
				"containers.$.cOutput":    cleanedOutput,
			},
		}
		err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			fmt.Println("Error updating AnalysisCollection (step 4).", err)
		}
	}

	// step 5: send output to the proper analysis result function.
	switch securityTest.Name {
	case "enry":
		EnryStartAnalysis(CID, cleanedOutput, analysis.RID)
	case "gas":
		GasStartAnalysis(CID, cleanedOutput)
	default:
		fmt.Println("Error: Could not find securityTest.Name.")
	}
}
