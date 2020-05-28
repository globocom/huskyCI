package runner

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	goContext "golang.org/x/net/context"
)

// DockerAPI implements the runner interface
type DockerAPI struct {
	GitSSHKey   string
	Artifactory string
	Session     *client.Client
}

// NewDockerAPISession starts a new DockerAPI session.
func NewDockerAPISession(lc fx.Lifecycle, settings *viper.Viper) (*DockerAPI, error) {

	dockerAPIHost := fmt.Sprintf("https://%s", settings.GetString("HUSKYCI_DOCKERAPI_ADDR"))
	dockerAPICertPath := settings.GetString("HUSKYCI_DOCKERAPI_CERT_PATH")
	dockerAPIVerifyTLS := settings.GetString("HUSKYCI_DOCKERAPI_TLS_VERIFY")

	// env vars needed by docker/docker library to create a new client:
	if err := os.Setenv("DOCKER_HOST", dockerAPIHost); err != nil {
		return &DockerAPI{}, err
	}
	if err := os.Setenv("DOCKER_CERT_PATH", dockerAPICertPath); err != nil {
		return &DockerAPI{}, err
	}
	if err := os.Setenv("DOCKER_TLS_VERIFY", dockerAPIVerifyTLS); err != nil {
		return &DockerAPI{}, err
	}

	client, err := client.NewEnvClient()
	if err != nil {
		return &DockerAPI{}, err
	}

	dockerAPISession := &DockerAPI{
		Artifactory: settings.GetString("HUSKYCI_ARTIFACTORY_URL"),
		GitSSHKey:   settings.GetString("HUSKYCI_API_GIT_PRIVATE_SSH_KEY"),
		Session:     client,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dockerAPISession.Ping()
		},
		OnStop: func(ctx context.Context) error {
			return dockerAPISession.Close()
		},
	})

	return dockerAPISession, nil
}

// Ping checks the DockerAPI session
func (d *DockerAPI) Ping() error {
	ctx := goContext.Background()
	fmt.Println("Checking DockerAPI Session...")
	_, err := d.Session.Ping(ctx)
	return err
}

// Close closes the DockerAPI session
func (d *DockerAPI) Close() error {
	fmt.Println("Closing DockerAPI Session...")
	return d.Session.Close()
}

// Run runs a container given its image plus a command and returns an error
func (d *DockerAPI) Run(containerImage string, command string) (string, error) {

	// step 1: pull image if it is not there yet
	if !d.ImageIsLoaded(containerImage) {
		if err := d.PullImageWorker(containerImage); err != nil {
			return "", err
		}
	}

	// step 2: create a new container given an image and it's cmd
	CID, err := d.CreateContainer(containerImage, command)
	if err != nil {
		return "", err
	}

	// step 3: start the container
	if err := d.StartContainer(CID); err != nil {
		return "", err
	}

	// step 4: wait the container to finish
	if err := d.WaitContainer(CID); err != nil {
		return "", err
	}

	// step 5: read the container's output
	cOutput, err := d.ReadOutput(CID)
	if err != nil {
		return "", err
	}

	// step 6: remove the container from docker API
	if err := d.RemoveContainer(CID); err != nil {
		return "", err
	}

	return cOutput, nil
}

// ImageIsLoaded returns a bool if a given image is loaded or not at Docker API
func (d *DockerAPI) ImageIsLoaded(image string) bool {
	args := filters.NewArgs()
	args.Add("reference", image)
	options := dockerTypes.ImageListOptions{Filters: args}
	ctx := goContext.Background()
	result, err := d.Session.ImageList(ctx, options)
	if err != nil {
		fmt.Println(err) //log
	}
	return len(result) != 0
}

// PullImage pulls an image, like docker pull.
func (d *DockerAPI) PullImage(image string) error {
	canonicalURL := fmt.Sprintf("%s/%s", d.Artifactory, image)
	ctx := goContext.Background()
	_, err := d.Session.ImagePull(ctx, canonicalURL, dockerTypes.ImagePullOptions{})
	if err != nil {
		// log
	}
	return err
}

// PullImageWorker tries to pull image every 15 seconds
func (d *DockerAPI) PullImageWorker(image string) error {
	timeout := time.After(15 * time.Minute)
	retryTick := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-timeout:
			timeOutErr := errors.New("timeout")
			// log
			return timeOutErr
		case <-retryTick.C:
			// log
			if d.ImageIsLoaded(image) {
				// log
				return nil
			}
			if err := d.PullImage(image); err != nil {
				// log
				return err
			}
		}
	}
}

// CreateContainer creates a new container and return its CID and an error
func (d *DockerAPI) CreateContainer(image, cmd string) (string, error) {
	ctx := goContext.Background()
	resp, err := d.Session.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   true,
		Cmd:   []string{"/bin/sh", "-c", cmd},
	}, nil, nil, "")
	if err != nil {
		return "", err // log
	}
	return resp.ID, nil
}

// StartContainer starts a container and returns its error.
func (d *DockerAPI) StartContainer(CID string) error {
	ctx := goContext.Background()
	return d.Session.ContainerStart(ctx, CID, dockerTypes.ContainerStartOptions{})
}

// WaitContainer returns when container finishes executing cmd.
func (d *DockerAPI) WaitContainer(CID string) error {
	ctx := goContext.Background()
	statusCode, err := d.Session.ContainerWait(ctx, CID)
	if statusCode != 0 {
		return fmt.Errorf("Error in POST to wait the container with statusCode %d", statusCode)
	}
	return err
}

// StopContainer stops an active container by it's CID
func (d *DockerAPI) StopContainer(CID string) error {
	ctx := goContext.Background()
	err := d.Session.ContainerStop(ctx, CID, nil)
	if err != nil {
		// log
	}
	return err
}

// RemoveContainer removes a container by it's CID
func (d *DockerAPI) RemoveContainer(CID string) error {
	ctx := goContext.Background()
	err := d.Session.ContainerRemove(ctx, CID, dockerTypes.ContainerRemoveOptions{})
	if err != nil {
		// log
	}
	return err
}

// ReadOutput returns STDOUT of a given containerID.
func (d *DockerAPI) ReadOutput(CID string) (string, error) {
	ctx := goContext.Background()
	out, err := d.Session.ContainerLogs(ctx, CID, dockerTypes.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		// log
		return "", err
	}
	body, err := ioutil.ReadAll(out)
	if err != nil {
		// log
		return "", err
	}
	return string(body), nil
}
