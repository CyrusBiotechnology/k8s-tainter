package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type NodeTainter struct {
	Watch    *watch.Interface
	Selector v1.LabelSelector
	Taints   []Taints
}

func (nt *NodeTainter) Do(stop chan bool) {
	c := (*nt.Watch).ResultChan()
	for {
		select {
		case <-c:
			// if match label selector (watch is shared):
			//   handle node event
		case <-stop:
			return
		}
	}
}
