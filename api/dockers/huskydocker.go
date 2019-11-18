// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dockers

import (
	"errors"
	"fmt"
	"time"

	"regexp"

	"github.com/globocom/huskyCI/api/log"
)

const logActionRun = "DockerRun"
const logInfoHuskyDocker = "HUSKYDOCKER"
const logActionPull = "pullImage"

const urlRegexp = `([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`

func configureImagePath(image, tag string) (string, string) {
	fullContainerImage := fmt.Sprintf("%s:%s", image, tag)
	regex := regexp.MustCompile(urlRegexp)
	canonicalURL := image
	if !regex.MatchString(canonicalURL) {
		canonicalURL = fmt.Sprintf("docker.io/%s", fullContainerImage)
	} else {
		canonicalURL = fullContainerImage
	}

	return canonicalURL, fullContainerImage
}

// DockerRun starts a new container and returns its output and an error.
func DockerRun(image, imageTag, cmd string, timeOutInSeconds int) (string, string, error) {

	// step 1: create a new docker API client
	d, err := NewDocker()
	if err != nil {
		return "", "", err
	}

	canonicalURL, fullContainerImage := configureImagePath(image, imageTag)
	// step 2: pull image if it is not there yet
	if !d.ImageIsLoaded(fullContainerImage) {
		if err := pullImage(d, canonicalURL, fullContainerImage); err != nil {
			return "", "", err
		}
	}

	// step 3: create a new container given an image and it's cmd
	CID, err := d.CreateContainer(fullContainerImage, cmd)
	if err != nil {
		return "", "", err
	}
	d.CID = CID

	// step 4: start container
	if err := d.StartContainer(); err != nil {
		log.Error(logActionRun, logInfoHuskyDocker, 3015, err)
		return "", "", err
	}
	log.Info(logActionRun, logInfoHuskyDocker, 32, fullContainerImage, d.CID)

	// step 5: wait container finish
	if err := d.WaitContainer(timeOutInSeconds); err != nil {
		log.Error(logActionRun, logInfoHuskyDocker, 3016, err)
		return "", "", err
	}

	// step 6: read container's output when it finishes
	cOutput, err := d.ReadOutput()
	if err != nil {
		return "", "", err
	}
	log.Info(logActionRun, logInfoHuskyDocker, 34, fullContainerImage, d.CID)

	// step 7: remove container from docker API
	if err := d.RemoveContainer(); err != nil {
		log.Error(logActionRun, logInfoHuskyDocker, 3027, err)
		return "", "", err
	}

	return CID, cOutput, nil
}

func pullImage(d *Docker, canonicalURL, image string) error {
	timeout := time.After(15 * time.Minute)
	retryTick := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-timeout:
			timeOutErr := errors.New("timeout")
			log.Error(logActionPull, logInfoHuskyDocker, 3013, timeOutErr)
			return timeOutErr
		case <-retryTick.C:
			log.Info(logActionPull, logInfoHuskyDocker, 31, image)
			if d.ImageIsLoaded(image) {
				log.Info(logActionPull, logInfoHuskyDocker, 35, image)
				return nil
			}
			if err := d.PullImage(canonicalURL); err != nil {
				log.Error(logActionPull, logInfoHuskyDocker, 3013, err)
				return err
			}
		}
	}
}
