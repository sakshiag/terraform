package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCsWorker() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCsWorkerRead,

		Schema: map[string]*schema.Schema{
			"worker_id": {
				Description: "ID of the worker",
				Type:        schema.TypeString,
				Required:    true,
			},
			"state": {
				Description: "State of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "Status of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"private_vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCsWorkerRead(d *schema.ResourceData, meta interface{}) error {
	workerID := d.Get("worker_id").(string)

	workerClient := meta.(ClientSession).ClusterWorkerClient()

	targetEnv := getClusterTargetHeader(d)

	workerFields, err := workerClient.Get(workerID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving worker: %s", err)
	}

	d.SetId(workerFields.ID)
	d.Set("state", workerFields.State)
	d.Set("status", workerFields.Status)
	d.Set("private_vlan", workerFields.PrivateVlan)
	d.Set("public_vlan", workerFields.PublicVlan)
	d.Set("private_ip", workerFields.PrivateIP)
	d.Set("public_ip", workerFields.PublicIP)

	return nil
}
