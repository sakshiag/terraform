package brooklyn

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/brooklyncentral/brooklyn-cli/net"
)

const (
	// LOCAL_BROOKLYN_URL is the default endpoint for locally deployed Apache Brooklyn
	LOCAL_BROOKLYN_URL = "http://localhost:8081"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROOKLYN_URL", LOCAL_BROOKLYN_URL),
				Description: "A Brooklyn application URL.",
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROOKLYN_USERNAME", nil),
				Description: "A Brooklyn username.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROOKLYN_PASSWORD", nil),
				Description: "The Brooklyn password.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"brooklyn_application": resourceApplication(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	network := &net.Network{
		BrooklynUrl: d.Get("api_url").(string),
		BrooklynUser: d.Get("username").(string),
		BrooklynPass: d.Get("password").(string),
	}
	return network, nil
}
