package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func NewKubeClient(config *rest.Config) *kubernetes.Clientset {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	return clientSet
}
