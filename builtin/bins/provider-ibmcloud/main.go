package main

import (
	"log"

	"github.com/hashicorp/terraform/builtin/providers/ibmcloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	log.Println("IBM Cloud Provider version", "tf-v0.9.3-ibm-provider-v0.1")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ibmcloud.Provider,
	})
}
