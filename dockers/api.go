package dockers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/globocom/glbgelf"
	"github.com/globocom/huskyci/context"
	"github.com/globocom/huskyci/types"
	goContext "golang.org/x/net/context"
)

const (
	certFile   = "cert.pem"
	keyFile    = "key.pem"
	carootFile = "ca.pem"
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

// NewDocker returns a new docker.
func NewDocker() (*Docker, error) {
	configAPI := context.GetAPIConfig()
	dockerHost := fmt.Sprintf("http://%s", configAPI.DockerHostsConfig.Host)
	glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "NewDocker",
		"info":   "API"}, "INFO", "dockerHost:", dockerHost)

	err := os.Setenv("DOCKER_HOST", dockerHost)
	if err != nil {
		return nil, err
	}

	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	docker := &Docker{
		client: client,
	}
	return docker, nil
}

// NewClientTLS returns an http client with certificate authentication.
func (d Docker) NewClientTLS() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(carootFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:               tls.VersionTLS11,
				MaxVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				InsecureSkipVerify:       false,
				Certificates:             []tls.Certificate{cert},
				RootCAs:                  caCertPool,
			},
		},
	}
	return client, nil
}

// CreateContainer creates a new container
func (d Docker) CreateContainer(analysis types.Analysis, image string, cmd string) (string, error) {
	cmd = handleCmd(analysis.URL, analysis.Branch, cmd)
	ctx := goContext.Background()
	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   true,
		Cmd:   []string{"/bin/sh", "-c", cmd},
	}, nil, nil, "")

	if err != nil {
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

// ReadOutput returns the command ouput of a given containerID.
func (d Docker) ReadOutput() (string, error) {
	ctx := goContext.Background()
	out, err := d.client.ContainerLogs(ctx, d.CID, dockerTypes.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", nil
	}

	body, err := ioutil.ReadAll(out)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// PullImage pulls an image, like docker pull.
func (d Docker) PullImage(image string) error {
	ctx := goContext.Background()
	_, err := d.client.ImagePull(ctx, image, dockerTypes.ImagePullOptions{})
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
		panic(err)
	}

	return len(result) != 0
}

// ListImages returns the docker images, like docker image ls.
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
		glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "HealthCheckDockerAPI",
			"info":   "API"}, "ERROR", "Error HealthCheckDockerAPI():", err)
		return err
	}

	ctx := goContext.Background()
	_, err = d.client.Ping(ctx)
	return err
}

// handleCmd will extract %GIT_REPO% and %GIT_BRANCH%  from cmd and replace it with the proper repository URL.
func handleCmd(repositoryURL, repositoryBranch, cmd string) string {
	replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
	replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
	return replace2
}
