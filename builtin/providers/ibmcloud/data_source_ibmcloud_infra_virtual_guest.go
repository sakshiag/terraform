package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func dataSourceIBMCloudInfraVirtualGuest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudInfraVirtualGuestRead,

		Schema: map[string]*schema.Schema{

			"hostname": &schema.Schema{
				Description: "The hostname of the virtual guest",
				Type:        schema.TypeString,
				Required:    true,
			},

			"domain": &schema.Schema{
				Description: "The domain of the virtual guest",
				Type:        schema.TypeString,
				Required:    true,
			},

			"datacenter": &schema.Schema{
				Description: "Datacenter in which the virtual guest is deployed",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"cores": &schema.Schema{
				Description: "Number of cpu cores",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"status": &schema.Schema{
				Description: "The power status of the virtual guest",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMCloudInfraVirtualGuestRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	hostname := d.Get("hostname").(string)
	domain := d.Get("domain").(string)

	vgs, err := service.
		Filter(filter.Build(filter.Path("virtualGuests.hostname").Eq(hostname),
			filter.Path("virtualGuests.domain").Eq(domain))).Mask(
		"hostname,domain,startCpus,datacenter[id,name,longName],statusId,status,id,powerState",
	).GetVirtualGuests()

	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest details for host %s: %s", hostname, err)
	}
	if len(vgs) == 0 {
		return fmt.Errorf("No virtual guest with hostname %s and domain  %s", hostname, domain)
	}
	var vg datatypes.Virtual_Guest

	vg = vgs[0]
	d.Set("hostname", *vg.Hostname)
	d.Set("domain", *vg.Domain)

	if vg.Datacenter != nil {
		d.Set("datacenter", *vg.Datacenter.Name)
	}

	d.Set("cores", *vg.StartCpus)
	d.Set("status", *vg.PowerState.KeyName)
	d.SetId(fmt.Sprintf("%d", *vg.Id))
	return nil
}
