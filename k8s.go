package main

import (
	"errors"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os/user"
	"path"
)

func initK8s(configPath string) (*kubernetes.Clientset, error) {
	var k8s *kubernetes.Clientset

	// if no kubeconfig is specified, try loading from the user's home dir
	if configPath == "" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		configPath = path.Join(dir, ".kube", "config")
	}

	// try loading the kubeconfig first, then fall back to the in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		var inClusterErr error
		config, inClusterErr = rest.InClusterConfig()
		// uses the current context in kubeconfig
		if err != nil {
			return nil, errors.New(fmt.Sprint("could not get kubernetes config: [", err, ", ", inClusterErr, "]"))
		}
		log.Println("loaded k8s in-cluster config")
	} else {
		log.Println("loaded k8s config from: ", configPath)
	}

	// creates the clientset
	k8s, err = kubernetes.NewForConfig(config)
	return k8s, err
}
