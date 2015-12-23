package main

import (
	"github.com/hashicorp/terraform/builtin/providers/brooklyn"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: brooklyn.Provider,
	})
}

