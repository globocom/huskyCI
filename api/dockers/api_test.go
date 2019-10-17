package dockers_test

import (
	"os"
	"testing"

	apiContext "github.com/globocom/huskyCI/api/context"
	. "github.com/globocom/huskyCI/api/dockers"
	"github.com/globocom/huskyCI/api/log"
)

func TestIntegration(t *testing.T) {
	docker, err := newDockerTest()
	if err != nil {
		t.Error(err)
		return
	}

	err = HealthCheckDockerAPI()
	if err != nil {
		t.Error(err)
		return
	}

	testcases := []struct {
		name          string
		docker        *Docker
		image         string
		cmd           string
		timeoutInSecs int
		waitFor       int
		wantOutput    string
		wantOutputErr string
		wantErr       string
	}{
		{
			name:          "TestStartContainer: starting docker container",
			docker:        docker,
			image:         "library/alpine:latest",
			cmd:           "echo 'huskyCI rules'",
			timeoutInSecs: 5,
			wantOutput:    "huskyCI rules\r\n",
		},
		{
			name:    "Test: repository name no canonical",
			docker:  docker,
			image:   "1234567890",
			cmd:     "ls",
			wantErr: "repository name must be canonical",
		},
		{
			name:    "Test: image not found",
			docker:  docker,
			image:   "1234567890/1234567890:1234567890",
			cmd:     "ls",
			wantErr: "Error response from daemon: pull access denied for 1234567890/1234567890, repository does not exist or may require 'docker login': denied: requested access to the resource is denied",
		},
		{
			name:    "Test: command not found",
			docker:  docker,
			image:   "library/alpine:latest",
			cmd:     "1234567890",
			wantErr: "Error in POST to wait the container with statusCode 127",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			_, output, err := DockerRun(testcase.image, testcase.cmd, testcase.timeoutInSecs)
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in docker run, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			} else if output != testcase.wantOutput {
				t.Errorf("in docker run, we expected %x; but got %x", testcase.wantOutput, output)
				return
			}

			if !testcase.docker.ImageIsLoaded(testcase.image) {
				err := testcase.docker.PullImage("docker.io/" + testcase.image)
				if err != nil {
					if err.Error() != testcase.wantErr {
						t.Errorf("in pulling image, we expected %s; but got %s", testcase.wantErr, err.Error())
					}
				}
				err = PullImage(testcase.docker, testcase.image)
				if err != nil {
					t.Error(err)
				}
			}

			images, err := testcase.docker.ListImages()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in listing images, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			} else if len(images) == 0 {
				t.Error("in listing images, we expected at least 1 image; but got 0")
				return
			}

			defer func() {
				_, err = testcase.docker.RemoveImage(testcase.image)
				if err != nil {
					if err.Error() != testcase.wantErr {
						t.Errorf("in removing the image, we expected %s; but got %s", testcase.wantErr, err.Error())
					}
					return
				}
			}()

			cid, err := docker.CreateContainer(testcase.image, testcase.cmd)
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in creating a container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			}
			testcase.docker.CID = cid

			defer func() {
				err = testcase.docker.DieContainers()
				if err != nil {
					if err.Error() != testcase.wantErr {
						t.Errorf("in stopping & removing all containers, we expected %s; but got %s", testcase.wantErr, err.Error())
					}
					return
				}
			}()

			err = testcase.docker.StartContainer()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in starting the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			}

			err = testcase.docker.WaitContainer(testcase.waitFor)
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in waiting for the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			}

			output, err = testcase.docker.ReadOutput()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in reading output the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			} else if output != testcase.wantOutput {
				t.Errorf("in reading output the container, we expected %x; but got %x", testcase.wantOutput, output)
				return
			}

			outputErr, err := testcase.docker.ReadOutputStderr()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in reading error output the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			} else if outputErr != testcase.wantOutputErr {
				t.Errorf("in reading output the container, we expected %s; but got %s", testcase.wantOutputErr, outputErr)
				return
			}

			err = testcase.docker.StopContainer()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in stopping the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			}

			stopped, err := testcase.docker.ListStoppedContainers()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in listing stopped containers, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			} else if len(stopped) == 0 {
				t.Error("in listing stopped containers, we expected at least 1 stopped container; but got 0")
				return
			}

			err = testcase.docker.RemoveContainer()
			if err != nil {
				if err.Error() != testcase.wantErr {
					t.Errorf("in removing the container, we expected %s; but got %s", testcase.wantErr, err.Error())
				}
				return
			}

		})
	}
}

func newDockerTest() (*Docker, error) {
	// fix for getting the config.yml file
	os.Chdir("../..")
	apiContext.APIConfiguration = &apiContext.APIConfig{
		GraylogConfig: &apiContext.GraylogConfig{
			DevelopmentEnv: true,
			AppName:        "docker_test",
			Tag:            "docker_test",
		},
	}
	docker, err := NewDocker()
	log.InitLog()
	return docker, err
}

