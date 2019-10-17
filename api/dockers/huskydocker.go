// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dockers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/globocom/huskyCI/api/log"
)

// DockerRun starts a new container and returns its output and an error.
func DockerRun(fullContainerImage, cmd string, timeOutInSeconds int) (string, string, error) {

	// step 1: create a new docker API client
	d, err := NewDocker()
	if err != nil {
		return "", "", err
	}

	// step 2: pull image if it is not there yet
	if !d.ImageIsLoaded(fullContainerImage) {
		if err := PullImage(d, fullContainerImage); err != nil {
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
		log.Error("DockerRun", "HUSKYDOCKER", 3015, err)
		return "", "", err
	}
	log.Info("DockerRun", "HUSKYDOCKER", 32, fullContainerImage, d.CID)

	// step 5: wait container finish
	if err := d.WaitContainer(timeOutInSeconds); err != nil {
		log.Error("DockerRun", "HUSKYDOCKER", 3016, err)
		return "", "", err
	}

	// step 6: read container's output when it finishes
	cOutput, err := d.ReadOutput()
	if err != nil {
		return "", "", err
	}
	log.Info("DockerRun", "HUSKYDOCKER", 34, fullContainerImage, d.CID)

	// step 7: remove container from docker API
	if err := d.RemoveContainer(); err != nil {
		log.Error("DockerRun", "HUSKYDOCKER", 3027, err)
		return "", "", err
	}

	return CID, cOutput, nil
}

// PullImage pulls docker images. If there's an error it retries each 3 seconds only 3 times
func PullImage(docker *Docker, canonicalImg string) error {
	canonicalURL := fmt.Sprintf("docker.io/%s", canonicalImg)
	retryTick := time.NewTicker(3 * time.Second)
	retries := 3

	splitted := strings.Split(canonicalImg, "/")
	image := splitted[len(splitted)-1]

	for {
		select {
		case <-retryTick.C:
			log.Info("pullImage", "DOCKERRUN", 31, image)
			if err := docker.PullImage(canonicalURL); err != nil {
				log.Error("pullImage", "HUSKYDOCKER", 3013, err)
				return err
			}
			retries--
			if retries == 0 {
				err := errors.New("no left retries")
				log.Error("pullImage", "HUSKYDOCKER", 3013, err)
				return err
			}
		default:
			if docker.ImageIsLoaded(image) || docker.ImageIsLoaded(canonicalImg) {
				log.Info("pullImage", "HUSKYDOCKER", 35, image)
				return nil
			}
		}
	}
}
