package main

import (
	"controller/signals"
	"controller/utils"
	"k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	kubeClientSet := utils.NewKubeClient()
	serviceWatch := createServiceWatcher(kubeClientSet)
	wq := createWorkQueue()
	indexer, informer := createIndexerAndInformer(serviceWatch, wq)

	controller := utils.NewController(wq, indexer, informer)

	stop := signals.SetupSignalHandler()

	go controller.Run(1, stop)
	// Wait forever
	select {}
}

func createServiceWatcher(clientSet *kubernetes.Clientset) *cache.ListWatch {
	return cache.NewListWatchFromClient(clientSet.CoreV1().RESTClient(), "services", v1.NamespaceDefault, fields.Everything())
}

func createWorkQueue() workqueue.RateLimitingInterface {
	return workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
}

func createIndexerAndInformer(watcher *cache.ListWatch, queue workqueue.RateLimitingInterface) (cache.Indexer, cache.Controller) {
	return cache.NewIndexerInformer(watcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
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
}
