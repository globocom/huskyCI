// Copyright 2021 Globo authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kubernetes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	goContext "golang.org/x/net/context"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubernetes is the Kubernetes struct
type Kubernetes struct {
	PID              string `json:"Id"`
	client           *kube.Clientset
	Namespace        string
	ProxyAddress     string
	NoProxyAddresses string
}

const logActionNew = "NewKubernetes"
const logInfoAPI = "KUBERNETES"

func NewKubernetes() (*Kubernetes, error) {
	configAPI, err := apiContext.DefaultConf.GetAPIConfig()
	if err != nil {
		log.Error(logActionNew, logInfoAPI, 3026, err)
		return nil, err
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configAPI.KubernetesConfig.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kube.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubernetes := &Kubernetes{
		client:           clientset,
		Namespace:        configAPI.KubernetesConfig.Namespace,
		ProxyAddress:     configAPI.KubernetesConfig.ProxyAddress,
		NoProxyAddresses: configAPI.KubernetesConfig.NoProxyAddresses,
	}

	return kubernetes, nil

}

func (k Kubernetes) CreatePod(image, cmd, podName, securityTestName string) (string, error) {

	ctx := goContext.Background()

	podToCreate := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
			Labels: map[string]string{
				"name":    podName,
				"huskyCI": securityTestName,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            podName,
					Image:           image,
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"/bin/sh",
						"-c",
						cmd,
					},
					Env: []core.EnvVar{
						{
							Name:  "http_proxy",
							Value: k.ProxyAddress,
						},
						{
							Name:  "https_proxy",
							Value: k.ProxyAddress,
						},
						{
							Name:  "no_proxy",
							Value: k.NoProxyAddresses,
						},
					},
				},
			},
			TopologySpreadConstraints: []core.TopologySpreadConstraint{
				{
					MaxSkew:           1,
					TopologyKey:       "kubernetes.io/hostname",
					WhenUnsatisfiable: "ScheduleAnyway",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"huskyCI": securityTestName,
						},
					},
				},
			},
			RestartPolicy: "Never",
		},
	}

	pod, err := k.client.CoreV1().Pods(k.Namespace).Create(ctx, podToCreate, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return string(pod.UID), nil
}

func (k Kubernetes) WaitPod(name string, podSchedulingTimeoutInSeconds, testTimeOutInSeconds int) (string, error) {

	ctx := goContext.Background()

	timeout := func(i int64) *int64 { return &i }(int64(podSchedulingTimeoutInSeconds))
	fmt.Printf("Timeout 1 %s - %+v\n", name, *timeout)
	schedulingTimeout := true

	watch, err := k.client.CoreV1().Pods(k.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector:  fmt.Sprintf("name=%s", name),
		Watch:          true,
		TimeoutSeconds: timeout,
	})
	if err != nil {
		fmt.Printf("Error watch 1 %s\n", name)
		return "", err
	}

	fmt.Printf("Start %s scheduling - %+v\n", name, time.Now())
schedulingLoop:
	for event := range watch.ResultChan() {
		p, ok := event.Object.(*core.Pod)
		if !ok {
			fmt.Printf("Error %s\n", name)
			return "", errors.New("Unexpected Event while waiting for Pod")
		}

		fmt.Printf("Scheduling loop %s - %+v\n", name, p)

		switch p.Status.Phase {
		case "Running":
			schedulingTimeout = false
			break schedulingLoop
		case "Succeeded", "Completed":
			return string(p.Status.Phase), nil
		case "Failed":
			return "", errors.New("Pod execution failed")
		case "Unknown":
			return "", errors.New("Pod terminated with Unknown status")
		}
	}

	fmt.Printf("schedulingTimeout %s - %+v - time %+v\n", name, schedulingTimeout, time.Now())
	if schedulingTimeout {
		err = k.RemovePod(name)
		if err != nil {
			return "", err
		}

		return "", errors.New(fmt.Sprintf("Timed-out waiting for pod scheduling: %s", name))
	}

	timeout_result := func(i int64) *int64 { return &i }(int64(testTimeOutInSeconds))

	fmt.Printf("Timeout 2 %s - %+v\n", name, *timeout_result)

	watch, err = k.client.CoreV1().Pods(k.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector:  fmt.Sprintf("name=%s", name),
		Watch:          true,
		TimeoutSeconds: timeout_result,
	})
	if err != nil {
		fmt.Printf("Error watch 2 %s\n", name)
		return "", err
	}

	fmt.Printf("Watch 2 - %+v\n", watch)
	for event := range watch.ResultChan() {
		p, ok := event.Object.(*core.Pod)
		if !ok {
			fmt.Printf("Error %s\n", name)
			return "", errors.New("Unexpected Event while waiting for Pod")
		}

		fmt.Printf("Waiting result %s - %+v\n", name, p)
		switch p.Status.Phase {
		case "Succeeded", "Completed":
			return string(p.Status.Phase), nil
		case "Failed":
			return "", errors.New("Pod execution failed")
		case "Unknown":
			return "", errors.New("Pod terminated with Unknown status")
		}
	}

	fmt.Printf("waiting result timeout %s - time %+v\n", name, time.Now())
	err = k.RemovePod(name)
	if err != nil {
		fmt.Printf("Error removing pod %s\n", name)
		return "", err
	}

	return "", errors.New(fmt.Sprintf("Timed-out waiting for pod to finish: %s", name))
}

func (k Kubernetes) ReadOutput(name string) (string, error) {
	ctx := goContext.Background()

	req := k.client.CoreV1().Pods(k.Namespace).GetLogs(name, &core.PodLogOptions{})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		errRemovePod := k.RemovePod(name)
		if errRemovePod != nil {
			return "", errRemovePod
		}
		return "", err
	}
	defer podLogs.Close()

	result, err := ioutil.ReadAll(podLogs)
	if err != nil {
		errRemovePod := k.RemovePod(name)
		if errRemovePod != nil {
			return "", errRemovePod
		}
		return "", err
	}

	return string(result), nil
}

func (k Kubernetes) RemovePod(name string) error {
	ctx := goContext.Background()

	return k.client.CoreV1().Pods(k.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
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
