package main

import (
	"flag"
	"log"

	_ "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// Make GCP great again
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var (
	config     = flag.String("config", "", "path to tainter config file")
	kubeconfig = flag.String("kubeconfig", "", "path to the kubeconfig file")
)

func main() {
	flag.Parse()

	cfgB , err:= ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cfgB))
	var cfg Config
	yaml.Unmarshal(cfgB, &cfg)
	fmt.Println(cfg)

	k8s, err := initK8s(*kubeconfig)
	if err != nil {
		panic(err)
	}

	nodes, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	log.Printf("Successfully connected to kubernetes, %v nodes online\n", len(nodes.Items))

	watch, err := k8s.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	wc := watch.ResultChan()
	for event := range wc {
		node := event.Object.(*v1.Node)
		value, ok := node.Labels["node-type"]
		if !ok {
			continue
		} else {
			fmt.Println(node.Name, value)
		}
	}
}
