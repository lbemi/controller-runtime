package main

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func kubernetesConfig() *rest.Config {
	homeDir := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", homeDir+"/.kube/config")
	if err != nil {
		klog.Fatalf("Get kubeconfig error : %s", err.Error())
	}
	return config
}
