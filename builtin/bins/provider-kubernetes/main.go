package main

import (
	"log"

	"github.com/hashicorp/terraform/builtin/providers/kubernetes"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	log.Println("IBM Cloud Provider Kubernetes version", "tf-v0.9.3-ibm-k8s-v0.1")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubernetes.Provider,
	})
}
