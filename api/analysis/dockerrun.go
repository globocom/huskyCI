// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package analysis

import (
	"fmt"
	"time"

	"github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/db"
	docker "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"gopkg.in/mgo.v2/bson"
)

// maxRetryStartContainer is maximum number of retries allowed
var maxRetryStartContainer = 5

// listedContainers is number of containers listed
var listedContainers = 0

// DockerRun starts a new container, runs a given securityTest in it and then updates AnalysisCollection.
func DockerRun(RID string, analysis *types.Analysis, securityTest types.SecurityTest) {

	newContainer := types.Container{SecurityTest: securityTest}
	securityTest.Cmd = util.HandlePrivateSSHKey(securityTest.Cmd)

	d, err := docker.NewDocker()
	if err != nil {
		log.Error("DockerRun", "DOCKERRUN", 3012, err)
		return
	}

	// step 0: pull image
	err = dockerPullImage(d, securityTest.Image)
	if err != nil {
		log.Error("DockerRun", "DOCKERRUN", 3013, err)
		return
	}

	// step 1: create a new container.
	err = dockerRunCreateContainer(d, analysis, securityTest, newContainer)
	if err != nil {
		return
	}

	// step 2: start created container.
	err = dockerRunStartContainer(d, analysis)
	if err != nil {
		log.Error("DockerRun", "DOCKERRUN", 3015, err)
		return
	}

	// step 3: wait container finish running.
	err = dockerRunWaitContainer(d, securityTest.TimeOutInSeconds)
	if err != nil {
		// error timeout will enter here!
		log.Error("DockerRun", "DOCKERRUN", 3016, err)
		if err := dockerRunRegisterError(d, analysis); err != nil {
			log.Error("DockerRun", "DOCKERRUN", 3017, err)
			return
		}
		return
	}

	// step 4: read cmd output from container.
	cOutput, err := dockerRunReadOutput(d, analysis)
	if err != nil {
		log.Error("DockerRun", "DOCKERRUN", 3018, err)
		return
	}

	// step 5: send output to the proper analysis result function.
	switch securityTest.Name {
	case "enry":
		EnryStartAnalysis(d.CID, cOutput, analysis.RID)
	case "gosec":
		GosecStartAnalysis(d.CID, cOutput, analysis.RID)
	case "bandit":
		BanditStartAnalysis(d.CID, cOutput)
	case "brakeman":
		BrakemanStartAnalysis(d.CID, cOutput)
	case "retirejs":
		RetirejsStartAnalysis(d.CID, cOutput)
	case "safety":
		SafetyStartAnalysis(d.CID, cOutput)
	case "npmaudit":
		NpmAuditStartAnalysis(d.CID, cOutput)
	default:
		log.Error("DockerRun", "DOCKERRUN", 3018, err)
	}

	// step 6: update active containers counter and check if it's time to kill all containers
	err = updateAndCheckContainerList(d)
	if err != nil {
		log.Error("DockerRun", "DOCKERRUN", 3025, err)
		return
	}
}

// dockerRunCreateContainer creates a new container, updates the corresponding analysis into MongoDB and returns an error and a CID (container ID).
func dockerRunCreateContainer(d *docker.Docker, analysis *types.Analysis, securityTest types.SecurityTest, newContainer types.Container) error {

	analysisQuery := map[string]interface{}{"RID": analysis.RID}

	// step 0: Verifies if the container has already been created
	for _, container := range analysis.Containers {
		if container.SecurityTest.Name == securityTest.Name {
			return nil
		}
	}

	// step 1: creating a new container.
	CID, err := d.CreateContainer(*analysis, securityTest.Image, securityTest.Cmd)

	if err != nil {
		// error! update analysis with an error message and quit.
		log.Error("dockerRunCreateContainer", "DOCKERRUN", 3014, err)
		newContainer.CStatus = "error"
		analysis.Containers = append(analysis.Containers, newContainer)
		err := db.UpdateOneDBAnalysis(analysisQuery, *analysis)
		if err != nil {
			log.Error("dockerRunCreateContainer", "DOCKERRUN", 2007, err)
			return err // implement a maxRetry?
		}
		return err // implement a maxRetry?
	}

	// step 2: update analysis with the container's information.
	d.CID = CID
	newContainer.CID = CID
	newContainer.CStatus = "created"
	analysis.Containers = append(analysis.Containers, newContainer)
	err = db.UpdateOneDBAnalysis(analysisQuery, *analysis)
	if err != nil {
		log.Error("dockerRunCreateContainer", "DOCKERRUN", 2007, err)
	}
	return err
}

// dockerRunStartContainer starts a container, updates the corresponding analysis into MongoDB and returns an error.
func dockerRunStartContainer(d *docker.Docker, analysis *types.Analysis) error {
	analysisQuery := map[string]interface{}{"containers.CID": d.CID}
	var err error
	// Tries to start a container maxRetryStartContainer times
	for i := 0; i < maxRetryStartContainer; i++ {
		err := d.StartContainer()
		if err == nil {
			break
		}
		log.Warning("dockerRunStartContainer", "DOCKERRUN", 3015, err)
	}
	if err != nil {
		updateContainerAnalysisQuery := bson.M{
			"$set": bson.M{
				"containers.$.cStatus": "error",
			},
		}
		err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
		if err != nil {
			log.Error("dockerRunStartContainer", "DOCKERRUN", 2007, err)
			return err
		}
		return err
	}
	// err is nil, container started successfully
	log.Info("dockerRunStartContainer", "DOCKERRUN", 32, d.CID)
	startedAt := time.Now()
	updateContainerAnalysisQuery := bson.M{
		"$set": bson.M{
			"containers.$.cStatus":   "running",
			"containers.$.startedAt": startedAt,
		},
	}
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
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
		log.Error("dockerRunReadOutput", "DOCKERRUN", 3017, err)
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
	err = db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		log.Error("dockerRunReadOutput", "DOCKERRUN", 2007, err)
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
			"containers.$.cInfo":      "Error waiting the container to finish.",
		},
	}
	err := db.UpdateOneDBAnalysisContainer(analysisQuery, updateContainerAnalysisQuery)
	if err != nil {
		return err
	}
	return nil
}

func dockerPullImage(d *docker.Docker, image string) error {

	canonicalURL := fmt.Sprintf("docker.io/%s", image)

	if d.ImageIsLoaded(image) {
		return nil
	}

	if err := d.PullImage(canonicalURL); err != nil {
		return err
	}

	// wait for image to be pulled (3 Minutes)
	timeout := time.Now().Add(3 * time.Minute)
	for {
		if d.ImageIsLoaded(image) {
			return nil
		}
		if time.Now().Before(timeout) {
			time.Sleep(30 * time.Second)
			log.Info("dockerPullImage", "DOCKERRUN", 31, image)
		} else {
			break
		}
	}

	return nil
}

func updateAndCheckContainerList(d *docker.Docker) error {
	configAPI := context.GetAPIConfig()
	maxContainersAllowed := configAPI.DockerHostsConfig.MaxContainersAllowed
	listedContainers++
	if listedContainers >= maxContainersAllowed {
		messageLog := fmt.Sprintf("Maximum allowed: %d -- Listed containers: %d", maxContainersAllowed, listedContainers)
		log.Info("updateAndCheckContainerList", "DOCKERRUN", 33, messageLog)
		err := d.DieContainers()
		if err != nil {
			log.Error("updateAndCheckContainerList", "DOCKERRUN", 3024, err)
			return err
		}
		listedContainers = 0
	}
	return nil
}
