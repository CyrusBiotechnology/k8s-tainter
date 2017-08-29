package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

// Taints nodes that are received by the watcher, if they match the configured
// selectors.
type NodeTainter struct {
	// Watch interface to listen on
	Watch *watch.Interface
	// Selector to match against
	Selector v1.LabelSelector
	// Taints to apply
	Taints []Taint
	// Client to use to apply taings
	Client corev1.NodeInterface
}

// Run the tainter.
func (nt *NodeTainter) Do(stop chan interface{}) {
	c := (*nt.Watch).ResultChan()
	for {
		select {
		case event := <-c:
			// if match label selector /* (watch is shared) */ and not already appropriately tainted :
			//   taint node
			event.Object.(apiv1.Node)
			nt.Client.Update(client.)
		case <-stop:
			return
		}
	}
}
