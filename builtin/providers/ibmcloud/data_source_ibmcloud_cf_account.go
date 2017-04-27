package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfAccountRead,

		Schema: map[string]*schema.Schema{
			"org_guid": {
				Description: "The guid of the org",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfAccountRead(d *schema.ResourceData, meta interface{}) error {
	or := meta.(ClientSession).BluemixAcccountClient()

	orgGUID := d.Get("org_guid").(string)

	account, err := or.FindByOrg(orgGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving organisation: %s", err)
	}

	d.SetId(account.GUID)

	return nil
}
