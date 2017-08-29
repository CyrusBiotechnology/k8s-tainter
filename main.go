package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api/v1"

	// Make GCP great again
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	config     = flag.String("config", "", "path to tainter config file")
	kubeconfig = flag.String("kubeconfig", "", "path to the kubeconfig file")
)

func main() {
	// The parsed config file
	var cfg Config

	flag.Parse()

	// Load the config file from disk
	cfgB, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(fmt.Sprint("unable to read config: ", err))
	}
	err = yaml.Unmarshal(cfgB, &cfg)
	if err != nil {
		panic(fmt.Sprint("unable to parse config: ", err))
	}

	// Load k8s config
	k8s, err := initK8s(*kubeconfig)
	if err != nil {
		panic(err)
	}

	// Attempt to connect
	nodes, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	log.Printf("Successfully connected to kubernetes, %v nodes online\n", len(nodes.Items))

	// Start a tainter for each object

	for _, taint := range cfg.Taints {
		watch, err := k8s.CoreV1().Nodes().Watch(metav1.ListOptions{LabelSelector: taint.Selector})
		taint.Selector
	}
	watch, err := k8s.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	wc := watch.ResultChan()
	for event := range wc {
		node := event.Object.(*v1.Node)
		//value, ok := node.Labels["node-type"]
		//if !ok {
		//	continue
		//} else {
		//	fmt.Println(node.Name, value)
		//}
	}
}
