package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	// Make GCP great again
	"k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	config     = flag.String("config", "", "path to tainter config file")
	kubeconfig = flag.String("kubeconfig", "", "path to the kubeconfig file")
)

func main() {
	defer glog.Flush()

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
	glog.V(3).Infoln("Successfully connected to kubernetes, %v nodes online\n", len(nodes.Items))

	// create the node watcher
	nodeListWatcher := cache.NewListWatchFromClient(k8s.CoreV1().RESTClient(), "nodes", v1.NamespaceAll, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the node key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Node than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(nodeListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(2, stop)

	// Wait forever
	select {}

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
