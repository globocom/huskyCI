// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dockers

import (
	"fmt"

	"github.com/globocom/huskyCI/api/container"
)

// DockerRun starts a new container and returns its CID, output, and an error.
func DockerRun(image, tag, command string, timeOutInSeconds int) (string, string, error) {

	newContainer := container.Container{}

	newContainer.Image.Name = image
	newContainer.Image.Tag = tag
	newContainer.Command = command
	newContainer.Image.CanonicalURL = fmt.Sprintf("docker.io/%s:%s", image, tag)

	if err := newContainer.Run(); err != nil {
		return "", "", err
	}

	return newContainer.CID, newContainer.Output, nil
}
