// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dockers

import (
	"errors"
	"fmt"
	"time"

	"github.com/globocom/huskyCI/api/log"
)

// DockerRun starts a new container and returns its output and an error.
func DockerRun(containerImage, cmd string) (string, string, error) {

	// step 1: create a new docker API client
	d, err := NewDocker()
	if err != nil {
		return "", "", err
	}

	// step 2: pull image if it is not there yet
	if !d.ImageIsLoaded(containerImage) {
		if err := pullImage(d, containerImage); err != nil {
			return "", "", err
		}
	}

	// step 3: create a new container given an image and it's cmd
	CID, err := d.CreateContainer(containerImage, cmd)
	if err != nil {
		return "", "", err
	}
	d.CID = CID

	// step 4: start container
	if err := d.StartContainer(); err != nil {
		return "", "", err
	}
	log.Info("DockerRun", "HUSKYDOCKER", 32, containerImage)

	// step 5: read container's output
	cOutput, err := readOutput(d)
	if err != nil {
		return "", "", err
	}
	log.Info("DockerRun", "HUSKYDOCKER", 34, containerImage)

	return CID, cOutput, nil
}

func pullImage(d *Docker, image string) error {
	canonicalURL := fmt.Sprintf("docker.io/%s", image)
	timeout := time.After(15 * time.Minute)
	retryTick := time.Tick(15 * time.Second)
	for {
		select {
		case <-timeout:
			// log.Error("pullImage", "HUSKYDOCKER", 3013, err)
			return errors.New("time out")
		case <-retryTick:
			log.Info("pullImage", "DOCKERRUN", 31, image)
			if d.ImageIsLoaded(image) {
				// log.Info("pullImage", "HUSKYDOCKER", 3013, err)
				return nil
			}
			if err := d.PullImage(canonicalURL); err != nil {
				// log.Error("pullImage", "HUSKYDOCKER", 3013, err)
				return err
			}
		}
	}
}

func readOutput(d *Docker) (string, error) {
	timeout := time.After(15 * time.Minute)
	retryTick := time.Tick(15 * time.Second)
	for {
		select {
		case <-timeout:
			// log.Error("DockerRun", "HUSKYDOCKER", 3013, err)
			return "", errors.New("time out")
		case <-retryTick:
			cOutput, err := d.ReadOutput()
			if err != nil {
				// log.Error("DockerRun", "HUSKYDOCKER", 3013, err)
				return "", err
			}
			if cOutput != "" {
				// log.Info("DockerRun", "HUSKYDOCKER", 3013, err)
				return cOutput, nil
			}
		}
	}
}
