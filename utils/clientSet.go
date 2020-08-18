package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

func NewKubeClient() *kubernetes.Clientset {
	clientSet, err := kubernetes.NewForConfig(NewKubeController())
	if err != nil {
		klog.Fatal(err)
	}
	return clientSet
}
