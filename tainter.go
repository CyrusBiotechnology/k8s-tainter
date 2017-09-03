package main

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api/v1"
)

// Filter events to select the ones we care about
func nodeEventFilter(eventChan <-chan watch.Event, done chan struct{}) chan *v1.Node {
	nodeChan := make(chan *v1.Node)
	go func() {
		for {
			select {
			case event := <-eventChan:
				if !(event.Type == watch.Added || event.Type == watch.Modified) {
					return
				}
				node := event.Object.(*v1.Node)
				nodeChan <- node
			case <-done:
				// all done!
				return
			}
		}
	}()
	return nodeChan
}
