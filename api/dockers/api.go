// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dockers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	goContext "golang.org/x/net/context"
)

// Docker is the docker struct
type Docker struct {
	CID    string `json:"Id"`
	client *client.Client
}

// CreateContainerPayload is a struct that represents all data needed to create a container.
type CreateContainerPayload struct {
	Image string   `json:"Image"`
	Tty   bool     `json:"Tty,omitempty"`
	Cmd   []string `json:"Cmd"`
}

const logActionNew = "NewDocker"
const logInfoAPI = "DOCKERAPI"

// NewDocker returns a new docker.
func NewDocker() (*Docker, error) {
	configAPI, err := apiContext.DefaultConf.GetAPIConfig()
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3026, err)
		return nil, err
	}
	dockerHost := fmt.Sprintf("https://%s", configAPI.DockerHostsConfig.Host)

	// env vars needed by docker/docker library to create a NewEnvClient:
	err = os.Setenv("DOCKER_HOST", dockerHost)
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3001, err)
		return nil, err
	}

	err = os.Setenv("DOCKER_CERT_PATH", configAPI.DockerHostsConfig.PathCertificate)
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3019, err)
		return nil, err
	}

	tlsVerify := strconv.Itoa(configAPI.DockerHostsConfig.TLSVerify)
	err = os.Setenv("DOCKER_TLS_VERIFY", tlsVerify)
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3020, err)
		return nil, err
	}

	client, err := client.NewEnvClient()
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3002, err)
		return nil, err
	}
	docker := &Docker{
		client: client,
	}
	return docker, nil
}

// CreateContainer creates a new container and return its CID and an error
func (d Docker) CreateContainer(image, cmd string) (string, error) {
	ctx := goContext.Background()
	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   true,
		Cmd:   []string{"/bin/sh", "-c", cmd},
	}, nil, nil, "")

	if err != nil {
		log.Error("CreateContainer", logInfoAPI, 3005, err)
		return "", err
	}
	return resp.ID, nil
}

// StartContainer starts a container and returns its error.
func (d Docker) StartContainer() error {
	ctx := goContext.Background()
	return d.client.ContainerStart(ctx, d.CID, dockerTypes.ContainerStartOptions{})
}

// WaitContainer returns when container finishes executing cmd.
func (d Docker) WaitContainer(timeOutInSeconds int) error {
	ctx := goContext.Background()
	statusCode, err := d.client.ContainerWait(ctx, d.CID)

	if statusCode != 0 {
		return fmt.Errorf("Error in POST to wait the container with statusCode %d", statusCode)
	}

	return err
}

// StopContainer stops an active container by it's CID
func (d Docker) StopContainer() error {
	ctx := goContext.Background()
	err := d.client.ContainerStop(ctx, d.CID, nil)
	if err != nil {
		log.Error("StopContainer", logInfoAPI, 3022, err)
	}
	return err
}

// RemoveContainer removes a container by it's CID
func (d Docker) RemoveContainer() error {
	ctx := goContext.Background()
	err := d.client.ContainerRemove(ctx, d.CID, dockerTypes.ContainerRemoveOptions{})
	if err != nil {
		log.Error("RemoveContainer", logInfoAPI, 3023, err)
	}
	return err
}

// ListStoppedContainers returns a Docker type list with CIDs of stopped containers
func (d Docker) ListStoppedContainers() ([]Docker, error) {

	ctx := goContext.Background()
	dockerFilters := filters.NewArgs()
	dockerFilters.Add("status", "exited")
	options := dockerTypes.ContainerListOptions{
		Quiet:   true,
		All:     true,
		Filters: dockerFilters,
	}

	containerList, err := d.client.ContainerList(ctx, options)
	if err != nil {
		log.Error("ListContainer", logInfoAPI, 3021, err)
		return nil, err
	}

	var dockerList []Docker
	for _, c := range containerList {
		docker := Docker{
			CID:    c.ID,
			client: d.client,
		}
		dockerList = append(dockerList, docker)
	}

	return dockerList, nil
}

// DieContainers stops and removes all containers
func (d Docker) DieContainers() error {
	containerList, err := d.ListStoppedContainers()
	if err != nil {
		return err
	}
	for _, c := range containerList {
		err := c.StopContainer()
		if err != nil {
			return err
		}
	}
	for _, c := range containerList {
		err := c.RemoveContainer()
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadOutput returns STDOUT of a given containerID.
func (d Docker) ReadOutput() (string, error) {
	ctx := goContext.Background()
	out, err := d.client.ContainerLogs(ctx, d.CID, dockerTypes.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Error("ReadOutput", logInfoAPI, 3006, err)
		return "", nil
	}

	body, err := ioutil.ReadAll(out)
	if err != nil {
		log.Error("ReadOutput", logInfoAPI, 3007, err)
		return "", err
	}
	return string(body), err
}

// ReadOutputStderr returns STDERR of a given containerID.
func (d Docker) ReadOutputStderr() (string, error) {
	ctx := goContext.Background()
	out, err := d.client.ContainerLogs(ctx, d.CID, dockerTypes.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		log.Error("ReadOutputStderr", logInfoAPI, 3006, err)
		return "", nil
	}

	body, err := ioutil.ReadAll(out)
	if err != nil {
		log.Error("ReadOutputStderr", logInfoAPI, 3008, err)
		return "", err
	}
	return string(body), err
}

// PullImage pulls an image, like docker pull.
func (d Docker) PullImage(image string) error {
	ctx := goContext.Background()
	_, err := d.client.ImagePull(ctx, image, dockerTypes.ImagePullOptions{})
	if err != nil {
		log.Error("PullImage", logInfoAPI, 3009, err)
	}
	return err
}

// ImageIsLoaded returns a bool if a a docker image is loaded or not.
func (d Docker) ImageIsLoaded(image string) bool {
	args := filters.NewArgs()
	args.Add("reference", image)
	options := dockerTypes.ImageListOptions{Filters: args}

	ctx := goContext.Background()
	result, err := d.client.ImageList(ctx, options)
	if err != nil {
		log.Error("ImageIsLoaded", logInfoAPI, 3010, err)
		panic(err)
	}

	return len(result) != 0
}

// ListImages returns docker images, like docker image ls.
func (d Docker) ListImages() ([]dockerTypes.ImageSummary, error) {
	ctx := goContext.Background()
	return d.client.ImageList(ctx, dockerTypes.ImageListOptions{})
}

// RemoveImage removes an image.
func (d Docker) RemoveImage(imageID string) ([]dockerTypes.ImageDelete, error) {
	ctx := goContext.Background()
	return d.client.ImageRemove(ctx, imageID, dockerTypes.ImageRemoveOptions{Force: true})
}

// HealthCheckDockerAPI returns true if a 200 status code is received from dockerAddress or false otherwise.
func HealthCheckDockerAPI() error {
	d, err := NewDocker()
	if err != nil {
		log.Error("HealthCheckDockerAPI", logInfoAPI, 3011, err)
		return err
	}

	ctx := goContext.Background()
	_, err = d.client.Ping(ctx)
	return err
}
