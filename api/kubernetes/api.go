// Copyright 2021 Globo authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kubernetes

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/globocom/huskyCI/api/log"
	goContext "golang.org/x/net/context"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubernetes is the Kubernetes struct
type Kubernetes struct {
	PID    string `json:"Id"`
	client *kube.Clientset
}

// CreatePodPayload is a struct that represents all data needed to create a kubernetes Pod.
type CreatePodPayload struct {
	Image     string   `json:"Image"`
	Tty       bool     `json:"Tty,omitempty"`
	Cmd       []string `json:"Cmd"`
	Name      string   `json:"Name"`
	Namespace string   `json:"Namspace"`
}

const logActionNew = "NewKubernetes"
const logInfoAPI = "KUBERNETES"

func NewKubernetes() (*Kubernetes, error) {
	// configAPI, err := apiContext.DefaultConf.GetAPIConfig()
	// if err != nil {
	// 	log.Error(logActionNew, logInfoAPI, 3026, err)
	// 	return nil, err
	// }

	// kubeClusterAddress := "kubernetes.docker.internal"

	// kubeconfig := "~/.kube/config"
	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// 	} else {
	// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// 	}
	// kubeconfig = flag.String("kubeconfig", "/go/src/github.com/globocom/huskyCI/kubeconfig", "(optional) absolute path to the kubeconfig file")
	// flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/go/src/github.com/globocom/huskyCI/kubeconfig")
	if err != nil {
		return nil, err
	}
	// kubeConfig := &rest.Config{
	// 	Host: "kubernetes.docker.internal"
	// }

	// create the clientset
	clientset, err := kube.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubernetes := &Kubernetes{
		client: clientset,
	}

	return kubernetes, nil

}

func (k Kubernetes) CreatePod(id, image, cmd, name string) (string, error) {

	ctx := goContext.Background()

	podToCreate := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				id: id,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            name,
					Image:           image,
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"/bin/sh",
						"-c",
						cmd,
					},
				},
			},
			RestartPolicy: "Never",
		},
	}

	pod, err := k.client.CoreV1().Pods("default").Create(ctx, podToCreate, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return string(pod.UID), nil
}

func (k Kubernetes) WaitPod(id, name string, timeOutInSeconds int) (string, error) {

	ctx := goContext.Background()

	timeout := func(i int64) *int64 { return &i }(int64(timeOutInSeconds))

	watch, err := k.client.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{
		LabelSelector:  id,
		Watch:          true,
		TimeoutSeconds: timeout,
	})
	if err != nil {
		return "", err
	}

	for event := range watch.ResultChan() {
		p, ok := event.Object.(*core.Pod)
		if !ok {
			return "", errors.New("Unexpected Event Type while waiting for Pod")
		}
		fmt.Println(p.Status.Phase)
		switch p.Status.Phase {
		case "Succeeded":
			return string(p.Status.Phase), nil
		case "Failed":
			return "", errors.New("Pod execution failed")
		case "Unknown":
			return "", errors.New("Pod terminated with Unknown status")
		}
	}

	err = k.RemovePod(name)
	if err != nil {
		return "", err
	}

	return "", errors.New("Timed-out waiting for pod!")
}

func (k Kubernetes) ReadOutput(name string) (string, error) {
	ctx := goContext.Background()

	req := k.client.CoreV1().Pods("default").GetLogs(name, &core.PodLogOptions{})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		errRemovePod := k.RemovePod(name)
		if errRemovePod != nil {
			return "", errRemovePod
		}
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		errRemovePod := k.RemovePod(name)
		if errRemovePod != nil {
			return "", errRemovePod
		}
		return "", nil
	}
	return buf.String(), nil

}

func (k Kubernetes) RemovePod(name string) error {
	ctx := goContext.Background()

	return k.client.CoreV1().Pods("default").Delete(ctx, name, metav1.DeleteOptions{})
}

// HealthCheckKubernetesAPI returns true if a 200 status code is received from kubernetes or false otherwise.
func HealthCheckKubernetesAPI() error {
	k, err := NewKubernetes()
	if err != nil {
		log.Error("HealthCheckKubernetesAPI", logInfoAPI, 3011, err)
		return err
	}

	ctx := goContext.Background()
	_, err = k.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error("HealthCheckKubernetesAPI", logInfoAPI, 3011, err)
		return err
	}
	return nil
}
