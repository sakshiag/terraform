package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfServiceKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfServiceKeyCreate,
		Read:     resourceIBMCloudCfServiceKeyRead,
		Delete:   resourceIBMCloudCfServiceKeyDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " The name of the service key ",
			},

			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the service instance for which to create service key",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary parameters to pass along to the service broker. Must be a JSON object",
			},
			"credentials": {
				Description: "Credentials asociated with the key",
				Type:        schema.TypeMap,
				Computed:    true,
			},
		},
	}
}

func resourceIBMCloudCfServiceKeyCreate(d *schema.ResourceData, meta interface{}) error {
	serviceRepo := meta.(ClientSession).CloudFoundryServiceKeyClient()
	name := d.Get("name").(string)
	serviceInstanceGUID := d.Get("service_instance_guid").(string)
	var parameters map[string]interface{}

	if parameters, ok := d.GetOk("parameters"); ok {
		parameters = parameters.(map[string]interface{})
	}

	serviceKey, err := serviceRepo.Create(serviceInstanceGUID, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating service key: %s", err)
	}

	d.SetId(serviceKey.Metadata.GUID)

	return resourceIBMCloudCfServiceKeyRead(d, meta)
}

func resourceIBMCloudCfServiceKeyRead(d *schema.ResourceData, meta interface{}) error {
	serviceRepo := meta.(ClientSession).CloudFoundryServiceKeyClient()
	serviceKeyGUID := d.Id()

	serviceKey, err := serviceRepo.Get(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service key: %s", err)
	}
	d.Set("credentials", serviceKey.Entity.Credentials)

	return nil
}

func resourceIBMCloudCfServiceKeyDelete(d *schema.ResourceData, meta interface{}) error {
	serviceRepo := meta.(ClientSession).CloudFoundryServiceKeyClient()

	serviceKeyGUID := d.Id()

	err := serviceRepo.Delete(serviceKeyGUID)
	if err != nil {
		return fmt.Errorf("Error deleting service key: %s", err)
	}

	d.SetId("")

	return nil
}
