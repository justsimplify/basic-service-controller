package utils

import (
	"flag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var kubeController *rest.Config

func NewKubeController() *rest.Config {
	if kubeController == nil {
		var kubeconfig string
		var master string
		var err error

		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
		flag.StringVar(&master, "master", "", "master url")
		flag.Parse()

		// creates the connection
		kubeController, err = clientcmd.BuildConfigFromFlags(master, kubeconfig)
		if err != nil {
			klog.Fatal(err)
		}
	}
	return kubeController
}