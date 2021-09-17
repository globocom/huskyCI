// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kubernetes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/globocom/huskyCI/api/log"
)

const logActionRun = "KubernetesRun"
const logInfoHuskyKube = "HUSKYKUBE"

const urlRegexp = `([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`

func configureImagePath(image, tag string) (string, string) {
	fullContainerImage := fmt.Sprintf("%s:%s", image, tag)
	regex := regexp.MustCompile(urlRegexp)
	canonicalURL := image
	if !regex.MatchString(canonicalURL) {
		// canonicalURL = fmt.Sprintf("docker.io/%s", fullContainerImage)
		canonicalURL = fmt.Sprintf("%s", fullContainerImage)
	} else {
		canonicalURL = fullContainerImage
	}

	return canonicalURL, fullContainerImage
}

// KubeRun starts a new pod and returns its output and an error.
func KubeRun(image, imageTag, cmd, securityTestName, id string, timeOutInSeconds int) (string, string, error) {

	// step 1: create a new Kubernetes API client
	k, err := NewKubernetes()
	if err != nil {
		log.Error(logActionRun, logInfoHuskyKube, 5001, k.PID, err.Error())
		return "", "", err
	}
	log.Info(logActionRun, logInfoHuskyKube, 41, k.PID)

	_, fullContainerImage := configureImagePath(image, imageTag)
	podName := fmt.Sprintf("%s-%s", strings.ToLower(id), securityTestName)

	// step 3: create a new container given an image and it's cmd
	podUID, err := k.CreatePod(fullContainerImage, cmd, podName, securityTestName)
	if err != nil {
		log.Error(logActionRun, logInfoHuskyKube, 5002, fullContainerImage, k.PID, err.Error())
		return "", "", err
	}
	k.PID = podUID

	log.Info(logActionRun, logInfoHuskyKube, 42, fullContainerImage, k.PID)

	// step 5: wait container finish
	_, err = k.WaitPod(podName, timeOutInSeconds)
	if err != nil {
		log.Error(logActionRun, logInfoHuskyKube, 5003, fullContainerImage, k.PID, err.Error())
		return "", "", err
	}

	log.Info(logActionRun, logInfoHuskyKube, 43, fullContainerImage, k.PID)

	// step 6: read container's output when it finishes
	cOutput, err := k.ReadOutput(podName)
	if err != nil {
		log.Error(logActionRun, logInfoHuskyKube, 5004, fullContainerImage, k.PID, err.Error())
		return "", "", err
	}

	log.Info(logActionRun, logInfoHuskyKube, 44, fullContainerImage, k.PID)

	// step 7: remove container from docker API
	if err := k.RemovePod(podName); err != nil {
		log.Error(logActionRun, logInfoHuskyKube, 5005, fullContainerImage, k.PID, err.Error())
		return "", "", err
	}

	log.Info(logActionRun, logInfoHuskyKube, 45, fullContainerImage, k.PID)

	return podUID, cOutput, nil
}
