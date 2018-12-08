package analysis

import (
	"os"
	"strings"
	"time"

	"github.com/globocom/glbgelf"
	docker "github.com/globocom/huskyci/dockers"
	"github.com/globocom/huskyci/types"
	"gopkg.in/mgo.v2/bson"
)

// DockerRun starts a new container, runs a given securityTest in it and then updates AnalysisCollection.
func DockerRun(RID string, analysis *types.Analysis, securityTest types.SecurityTest) {

	newContainer := types.Container{SecurityTest: securityTest}
	securityTest.Cmd = handlePrivateSSHKey(securityTest.Cmd)

	d, err := docker.NewDocker()
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error NewDocker():", err)
		return
	}

	// step 0: pull image
	err = dockerPullImage(d, securityTest.Image)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error dockerPullImage():", err)
		return
	}

	// step 1: create a new container.
	err = dockerRunCreateContainer(d, analysis, securityTest, newContainer)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error dockerRunCreateContainer():", err)
		return
	}

	// step 2: start created container.
	err = dockerRunStartContainer(d, analysis)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error dockerRunStartContainer():", err)
		return
	}

	// step 3: wait container finish running.
	err = dockerRunWaitContainer(d, securityTest.TimeOutInSeconds)
	if err != nil {
		// error timeout will enter here!
		if err := dockerRunRegisterError(d, analysis); err != nil {
			glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "DockerRun",
				"info":   "DOCKERRUN"}, "ERROR", "Error dockerRunRegisterError():", err)
			return
		}
		return
	}

	// step 4: read cmd output from container.
	cOutput, err := dockerRunReadOutput(d, analysis)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error dockerRunReadOutput():", err)
		return
	}

	// step 5: send output to the proper analysis result function.
	switch securityTest.Name {
	case "enry":
		EnryStartAnalysis(d.CID, cOutput, analysis.RID)
	case "gosec":
		GosecStartAnalysis(d.CID, cOutput)
	case "bandit":
		BanditStartAnalysis(d.CID, cOutput)
	case "brakeman":
		BrakemanStartAnalysis(d.CID, cOutput)
	default:
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "DockerRun",
			"info":   "DOCKERRUN"}, "ERROR", "Error: Could not find securityTest.Name.")
	}
}

// dockerRunCreateContainer creates a new container, updates the corresponding analysis into MongoDB and returns an error and a CID (container ID).
func dockerRunCreateContainer(d *docker.Docker, analysis *types.Analysis, securityTest types.SecurityTest, newContainer types.Container) error {

	analysisQuery := map[string]interface{}{"RID": analysis.RID}

	// step 1: creating a new container.
	CID, err := d.CreateContainer(*analysis, securityTest.Image, securityTest.Cmd)

	if err != nil {
		// error! update analysis with an error message and quit.
		newContainer.CStatus = "error"
		analysis.Containers = append(analysis.Containers, newContainer)
		err := UpdateOneDBAnalysis(analysisQuery, *analysis)
		if err != nil {
			glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "dockerRunCreateContainer",
				"info":   "DOCKERRUN"}, "ERROR", "Error 1 dockerRunCreateContainer() UpdateOneDBAnalysis():", err)
			return err // implement a maxRetry?
		}
		return err // implement a maxRetry?
	}

	// step 2: update analysis with the container's information.
	d.CID = CID
	newContainer.CID = CID
	newContainer.CStatus = "created"
	analysis.Containers = append(analysis.Containers, newContainer)
	err = UpdateOneDBAnalysis(analysisQuery, *analysis)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "dockerRunCreateContainer",
			"info":   "DOCKERRUN"}, "ERROR", "Error 2 dockerRunCreateContainer() UpdateOneDBAnalysis():", err)
	}
	return err
}

// dockerRunStartContainer starts a container, updates the corresponding analysis into MongoDB and returns an error.
func dockerRunStartContainer(d *docker.Docker, analysis *types.Analysis) error {
	analysisQuery := map[string]interface{}{"containers.CID": d.CID}
	err := d.StartContainer()
	if err != nil {
		// error starting container. maxRetry?
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cStatus": "error",
			},
		}
		err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "dockerRunStartContainer",
				"info":   "DOCKERRUN"}, "ERROR", "Error updating AnalysisCollection (step 2-err):", err)
			return err
		}
		return err
	}
	startedAt := time.Now()
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cStatus":   "running",
			"containers.$.startedAt": startedAt,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		return err
	}
	return err
}

// dockerRunWaitContainer waits a container run its commands.
func dockerRunWaitContainer(d *docker.Docker, timeout int) error {
	err := d.WaitContainer(timeout)
	return err
}

// dockerRunReadOutput reads the output of a container and updates the corresponding analysis into MongoDB.
func dockerRunReadOutput(d *docker.Docker, analysis *types.Analysis) (string, error) {
	analysisQuery := map[string]interface{}{"containers.CID": d.CID}
	cOutput, err := d.ReadOutput()
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "dockerRunReadOutput",
			"info":   "DOCKERRUN"}, "ERROR", "Error reading output from container", d.CID, ":", err)
		return "", err // implement a maxRetry?
	}
	finishedAt := time.Now()
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cStatus":    "finished",
			"containers.$.finishedAt": finishedAt,
			"containers.$.cOutput":    cOutput,
		},
	}
	err = UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "dockerRunReadOutput",
			"info":   "DOCKERRUN"}, "ERROR", "Error updating AnalysisCollection (step 4).", err)
		return "", err
	}
	return cOutput, err
}

// dockerRunRegisterError updates the corresponding analysis into MongoDB with an error status.
func dockerRunRegisterError(d *docker.Docker, analysis *types.Analysis) error {

	analysisQuery := map[string]interface{}{"containers.CID": d.CID}
	finishedAt := time.Now()
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cStatus":    "finished",
			"containers.$.finishedAt": finishedAt,
			"containers.$.cResult":    "failed",
			"containers.$.cOutput":    "Error waiting the container to finish.",
		},
	}
	err := UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		return err
	}
	return nil
}

func handlePrivateSSHKey(rawString string) string {
	cmdReplaced := strings.Replace(rawString, "GIT_PRIVATE_SSH_KEY", os.Getenv("GIT_PRIVATE_SSH_KEY"), -1)
	return cmdReplaced
}

func dockerPullImage(d *docker.Docker, image string) error {

	if d.ImageIsLoaded(image) {
		return nil
	}

	if err := d.PullImage(image); err != nil {
		return err
	}

	// wait for image to be pulled (2 Minutes)
	timeout := time.Now().Add(2 * time.Minute)
	for {
		if d.ImageIsLoaded(image) {
			return nil
		}
		if time.Now().Before(timeout) {
			time.Sleep(5 * time.Second)
			glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "dockerPullImage",
				"info":   "DOCKERRUN"}, "INFO", "Waiting pull image...")
		} else {
			break
		}
	}

	return nil
}
