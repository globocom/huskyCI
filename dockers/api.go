package dockers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

func handleCmd(analysis types.Analysis, cmd string) string {
	cmdReplaced := strings.Replace(cmd, "%GIT_REPO%", analysis.URL, -1)
	return cmdReplaced
}

// CreateContainer creates a container and returns its ID
// use docker as a parameter?
func (d Docker) CreateContainer(analysis types.Analysis, image string, cmd string) (string, error) {

	dockerHost := os.Getenv("DOCKER_HOST")
	cmd = handleCmd(analysis, cmd)

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
	req, err := http.NewRequest("POST", "http://"+dockerHost+"/v1.24/containers/create", bytes.NewBuffer(jsonPayload))
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

// StartContainer starts a container and returns its error
func (d Docker) StartContainer() error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + d.CID + "/start"
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// WaitContainer returns when container finishes executing cmd
func (d Docker) WaitContainer() error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + d.CID + "/wait"
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST /wait:", err)
	}
	defer resp.Body.Close()
	return err
}

// ReadOutput returns the command ouput of a given containerID
func (d Docker) ReadOutput() (string, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + d.CID + "/logs?stdout=1"
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

// PullImage pulls an image, like docker pull
func (d Docker) PullImage(image string) error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/images/create?fromImage=" + image
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// ListImages returns the docker images, like docker image ls
func (d Docker) ListImages() string {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/images/json"
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
