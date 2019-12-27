package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/globocom/huskyCI/api/context"
	goContext "golang.org/x/net/context"
)

// Container holds all information regarding a container.
type Container struct {
	dockerClient *client.Client
	CID          string    `bson:"CID" json:"CID"`
	Command      string    `bson:"cmd" json:"cmd"`
	Output       string    `bson:"output" json:"output"`
	Result       string    `bson:"result" json:"result"`
	Image        Image     `bson:"image" json:"image"`
	StartedAt    time.Time `bson:"startedAt" json:"startedAt"`
	FinishedAt   time.Time `bson:"finishedAt" json:"finishedAt"`
}

// Image is the struct that holds all information regarding a container image.
type Image struct {
	Name string `bson:"name" json:"name"`
	Tag  string `bson:"tag" json:"tag"`
}

// NewDockerClient creates a new docker API client and set it to the container struct.
func (c *Container) NewDockerClient() error {

	if err := setDockerClientEnvs(); err != nil {
		return err
	}

	newClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	c.dockerClient = newClient
	return nil
}

// Run runs a container by creating and starting it.
func (c *Container) Run() error {

	if err := c.Create(); err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	return nil
}

// Create creates a new container, set its CID and return an error.
func (c *Container) Create() error {

	ctx := goContext.Background()

	fullImageName := fmt.Sprintf("%s:%s", c.Image.Name, c.Image.Tag)
	containerConfig := &container.Config{
		Image: fullImageName,
		Tty:   true,
		Cmd:   []string{"/bin/sh", "-c", c.Command},
	}

	resp, err := c.dockerClient.ContainerCreate(ctx, containerConfig, nil, nil, "")
	if err != nil {
		return err
	}

	c.CID = resp.ID
	return nil
}

// Start starts a container and returns its error.
func (c *Container) Start() error {

	ctx := goContext.Background()

	return c.dockerClient.ContainerStart(ctx, c.CID, dockerTypes.ContainerStartOptions{})
}

// Stop stops an active container by it's CID.
func (c *Container) Stop() error {

	ctx := goContext.Background()

	return c.dockerClient.ContainerStop(ctx, c.CID, nil)
}

// Remove removes a container by it's CID.
func (c *Container) Remove() error {

	ctx := goContext.Background()

	return c.dockerClient.ContainerRemove(ctx, c.CID, dockerTypes.ContainerRemoveOptions{})
}

// PullImage pulls an image, like docker pull.
func (c *Container) PullImage() error {

	ctx := goContext.Background()

	fullImageName := fmt.Sprintf("%s:%s", c.Image.Name, c.Image.Tag)
	_, err := c.dockerClient.ImagePull(ctx, fullImageName, dockerTypes.ImagePullOptions{})

	return err
}

// ListImages returns docker images, like docker image ls.
func (c *Container) ListImages() ([]dockerTypes.ImageSummary, error) {

	ctx := goContext.Background()

	return c.dockerClient.ImageList(ctx, dockerTypes.ImageListOptions{})
}

// RemoveImage removes an image.
func (c *Container) RemoveImage(imageID string) ([]dockerTypes.ImageDelete, error) {

	ctx := goContext.Background()

	return c.dockerClient.ImageRemove(ctx, imageID, dockerTypes.ImageRemoveOptions{Force: true})
}

// ReadOutput returns the output of a container based on isSTDERR and isSTDOUT bool parameters.
func (c *Container) ReadOutput(isSTDOUT, isSTDERR bool) (string, error) {

	ctx := goContext.Background()
	containerLogOptions := dockerTypes.ContainerLogsOptions{
		ShowStdout: isSTDOUT,
		ShowStderr: isSTDERR,
	}

	cOutput, err := c.dockerClient.ContainerLogs(ctx, c.CID, containerLogOptions)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(cOutput)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ImageIsLoaded returns a bool if a a docker image is loaded in DockerAPI or not.
func (c *Container) ImageIsLoaded() (bool, error) {

	ctx := goContext.Background()

	fullImageName := fmt.Sprintf("%s:%s", c.Image.Name, c.Image.Tag)
	args := filters.NewArgs()
	args.Add("reference", fullImageName)
	options := dockerTypes.ImageListOptions{Filters: args}

	result, err := c.dockerClient.ImageList(ctx, options)
	if err != nil {
		return false, err
	}

	isLoaded := (len(result) != 0)
	return isLoaded, nil
}

// HealthCheckDockerAPI pings DockerAPI to check if it is up and running.
func HealthCheckDockerAPI() error {

	var healthCheckContainer Container

	ctx := goContext.Background()

	err := healthCheckContainer.NewDockerClient()
	if err != nil {
		// log.Error("HealthCheckDockerAPI", logInfoAPI, 3011, err)
		return err
	}

	_, err = healthCheckContainer.dockerClient.Ping(ctx)
	return err
}

// setDockerClientEnvs sets env vars needed by docker/docker library to create a NewEnvClient.
func setDockerClientEnvs() error {

	dockerHost := fmt.Sprintf("https://%s", context.APIConfiguration.DockerHostsConfig.Host)
	pathCertificate := context.APIConfiguration.DockerHostsConfig.PathCertificate
	tlsVerify := strconv.Itoa(context.APIConfiguration.DockerHostsConfig.TLSVerify)

	// env vars needed by docker/docker library to create a NewEnvClient:
	if err := os.Setenv("DOCKER_HOST", dockerHost); err != nil {
		// log.Error(logActionNew, logInfoAPI, 3001, err)
		return err
	}

	if err := os.Setenv("DOCKER_CERT_PATH", pathCertificate); err != nil {
		// log.Error(logActionNew, logInfoAPI, 3019, err)
		return err
	}

	if err := os.Setenv("DOCKER_TLS_VERIFY", tlsVerify); err != nil {
		// log.Error(logActionNew, logInfoAPI, 3020, err)
		return err
	}

	return nil
}
