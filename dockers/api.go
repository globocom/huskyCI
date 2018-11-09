package dockers

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/globocom/husky/context"
	"github.com/globocom/husky/types"
	goContext "golang.org/x/net/context"
)

const (
	certFile   = "cert.pem"
	keyFile    = "key.pem"
	carootFile = "ca.pem"
)

// Docker is the docker struct
type Docker struct {
	CID string `json:"Id"`
}

// CreateContainerPayload is a struct that represents all data need to create a container.
type CreateContainerPayload struct {
	Image string   `json:"Image"`
	Tty   bool     `json:"Tty,omitempty"`
	Cmd   []string `json:"Cmd"`
}

// NewClient creates http client with certificate authentication
func (d Docker) NewClient() (*client.Client, error) {
	configAPI := context.GetAPIConfig()
	_ = os.Setenv("DOCKER_HOST", "" configAPI.DockerHostsConfig.Host)
	return client.NewEnvClient()
}

// NewClient creates http client with certificate authentication
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

func (d Docker) CreateContainer(analysis types.Analysis, image string, cmd string) (string, error) {
	ctx := goContext.Background()
	cli, err := d.NewClient()
	if err != nil {
		return "", err
	}

	cmd = handleCmd(analysis.URL, analysis.Branch, cmd)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
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
	cli, err := d.NewClient()
	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, d.CID, dockerTypes.ContainerStartOptions{})
}

// WaitContainer returns when container finishes executing cmd.
func (d Docker) WaitContainer(timeOutInSeconds int) error {
	ctx := goContext.Background()
	cli, err := d.NewClient()
	if err != nil {
		return err
	}

	statusCode, err := cli.ContainerWait(ctx, d.CID)
	if statusCode != 0 {
		return fmt.Errorf("Error in POST to wait the container with statusCode %d", statusCode)
	}

	return err
}

// ReadOutput returns the command ouput of a given containerID.
func (d Docker) ReadOutput() (string, error) {
	ctx := goContext.Background()
	cli, err := d.NewClient()
	if err != nil {
		return "", err
	}

	out, err := cli.ContainerLogs(ctx, d.CID, dockerTypes.ContainerLogsOptions{ShowStdout: true})
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
	configAPI := context.GetAPIConfig()
	URL := configAPI.DockerHostsConfig.GetUrlPull(image)

	client, err := d.NewClientTLS()
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}

	resp, err := client.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// ListImages returns the docker images, like docker image ls.
func (d Docker) ListImages() string {
	configAPI := context.GetAPIConfig()
	URL := configAPI.DockerHostsConfig.GetUrlList()

	client, err := d.NewClientTLS()
	if err != nil {
		fmt.Println("Error in GET to get the images list:", err)
	}

	resp, err := client.Get(URL)
	if err != nil {
		fmt.Println("Error in GET to get the images list:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the body response of GET to read the command output:", err)
	}
	return string(body)
}

// HealthCheckDockerAPI returns true if a 200 status code is received from dockerAddress or false otherwise.
func HealthCheckDockerAPI(dockerAddress string) error {
	configAPI := context.GetAPIConfig()
	URL := configAPI.DockerHostsConfig.GetUrlHealthCheck(dockerAddress)
	resp, err := http.Get(URL)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		finalError := fmt.Sprintf("%d status code", resp.StatusCode)
		return errors.New(finalError)
	}
	return nil
}

// handleCmd will extract %GIT_REPO% and %GIT_BRANCH%  from cmd and replace it with the proper repository URL.
func handleCmd(repositoryURL, repositoryBranch, cmd string) string {
	replace1 := strings.Replace(cmd, "%GIT_REPO%", repositoryURL, -1)
	replace2 := strings.Replace(replace1, "%GIT_BRANCH%", repositoryBranch, -1)
	return replace2
}
