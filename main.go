package main

import (
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	// Stop channels for all tainters
	//var tainterStops []chan interface{}

	flag.Parse()

	cfg.Load(*config)

	// Load k8s config
	k8s, err := initK8s(*kubeconfig)
	if err != nil {
		panic(err)
	}

	// Attempt to connect
	nodes, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	log.Printf("Successfully connected to kubernetes, %v nodes online\n", len(nodes.Items))

	// Fire up a watcher
	watch, err := k8s.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	nodeChan := nodeEventFilter(watch.ResultChan(), make(chan struct{}))

	var modified = false
	for {
		node := <-nodeChan
		go func() {
			node, modified = cfg.AddTaints(node)
			if modified {
				fmt.Printf("tainting node %v: %v\n", node.Name, node.Spec.Taints)
				k8s.CoreV1().Nodes().Update(node)
			}
		}()
	}
}
