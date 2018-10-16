package dockers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/globocom/husky/context"
	"github.com/globocom/husky/types"
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

// clientAPI is the HTTP client being used to query Docker API.
var clientAPI = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS11,
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       false,
		},
	},
}

// CreateContainer creates a container and returns its ID.
func (d Docker) CreateContainer(analysis types.Analysis, image string, cmd string) (string, error) {

	configAPI := context.GetAPIConfig()
	URL := fmt.Sprintf("http://%s:%d/v1.24/containers/create", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort)
	cmd = handleCmd(analysis.URL, analysis.Branch, cmd)

	createContainerPayload := CreateContainerPayload{
		Image: image,
		Tty:   true,
		Cmd:   []string{"/bin/sh", "-c", cmd},
	}

	jsonPayload, err := json.Marshal(createContainerPayload)
	if err != nil {
		fmt.Println("Error in JSON Marshal.")
		return "", err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in POST to create a container:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the body response of POST to create the container:", err)
		return "", err
	}
	err = json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println("Error reading container ID:", err)
		return "", err
	}

	return d.CID, err
}

// StartContainer starts a container and returns its error.
func (d Docker) StartContainer() error {
	configAPI := context.GetAPIConfig()
	dockerHost := fmt.Sprintf("%s:%d", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort)
	URL := fmt.Sprintf("http://%s/v1.24/containers/%s/start", dockerHost, d.CID)
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// WaitContainer returns when container finishes executing cmd.
func (d Docker) WaitContainer(timeOutInSeconds int) error {

	configAPI := context.GetAPIConfig()
	URL := fmt.Sprintf("http://%s:%d/v1.24/containers/%s/wait", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort, d.CID)
	jsonPayload := []byte{}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	reqTimeOut := time.Duration(timeOutInSeconds) * time.Duration(time.Second)
	clientAPI.Timeout = reqTimeOut

	resp, err := clientAPI.Do(req)
	if err != nil {
		// timeout will enter here!
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ReadOutput returns the command ouput of a given containerID.
func (d Docker) ReadOutput() (string, error) {
	configAPI := context.GetAPIConfig()
	dockerHost := fmt.Sprintf("%s:%d", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort)
	URL := fmt.Sprintf("http://%s/v1.24/containers/%s/logs?stdout=1", dockerHost, d.CID)
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// PullImage pulls an image, like docker pull.
func (d Docker) PullImage(image string) error {
	configAPI := context.GetAPIConfig()
	dockerHost := fmt.Sprintf("%s:%d", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort)
	URL := fmt.Sprintf("http://%s/v1.24/images/create?fromImage=%s", dockerHost, image)
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// ListImages returns the docker images, like docker image ls.
func (d Docker) ListImages() string {
	configAPI := context.GetAPIConfig()
	dockerHost := fmt.Sprintf("%s:%d", configAPI.DockerHostsConfig.Addresses[0], configAPI.DockerHostsConfig.DockerAPIPort)
	URL := fmt.Sprintf("http://%s/v1.24/images/json", dockerHost)
	resp, err := http.Get(URL)
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
	URL := fmt.Sprintf("http://%s/v1.24/version", dockerAddress)
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
