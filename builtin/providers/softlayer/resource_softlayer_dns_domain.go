package softlayer

import (
	"github.com/hashicorp/terraform/helper/schema"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"fmt"
	"strconv"
	"log"
)

func resourceSoftLayerDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerDnsDomainCreate,
		Read: resourceSoftLayerDnsDomainRead,
		Update: resourceSoftLayerDnsDomainUpdate,
		Delete: resourceSoftLayerDnsDomainDelete,
		Schema: map[string]*schema.Schema {
			"id": &schema.Schema{
				Type: 		schema.TypeInt,
				Computed: 	true,
			},

			"name": &schema.Schema{
				Type: 		schema.TypeString,
				Required: 	true,
			},

			"serial": &schema.Schema{
				Type: 		schema.TypeInt,
				Optional:	true,
				Computed:	true,
			},

			"update_date": &schema.Schema{
				Type:		schema.TypeString,
				Computed:	true,
			},

			"records": &schema.Schema{
				Type:		schema.TypeList,
				Optional:	true,
				Elem:		&schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_data": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"domain_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"expire": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"host": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"minimum_ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"mx_priority": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"refresh": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"contact_email ": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"retry ": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"record_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSoftLayerDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainService

	// prepare creation parameters
	opts := datatypes.SoftLayer_Dns_Domain_Template{
		Name: d.Get("name").(string),
	}

	if serial, ok := d.GetOk("serial"); ok {
		opts.Serial = serial.(int)
	}

	response, err := client.CreateObject(opts)
	if err != nil {
		return fmt.Errorf("Error creating Dns Domain: %s", err)
	}

	id := response.Id
	d.SetId(strconv.Itoa(id))
	log.Printf("[INFO] Created Dns Domain: %d", id)

	return resourceSoftLayerDnsDomainRead(d, meta)
}

func resourceSoftLayerDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainService

	dnsId, _ := strconv.Atoi(d.Id())

	dns_domain, err := client.GetObject(dnsId)
	if err != nil {
		return fmt.Errorf("Error retrieving Dns Domain %d: %s", dnsId, err)
	}

	d.Set("id", dns_domain.Id)
	d.Set("name", dns_domain.Name)
	d.Set("serial", dns_domain.Serial)
	d.Set("update_date", dns_domain.UpdateDate)

	return nil
}

func resourceSoftLayerDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO - update is not supported - implement delete-create?
	return fmt.Errorf("Not implemented. Update Dns Domain is currently unsupported")
}

func resourceSoftLayerDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainService

	dnsId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Dns Domain: %s", err)
	}

	log.Printf("[INFO] Deleting Dns Domain: %d", dnsId)
	result, err := client.DeleteObject(dnsId)
	if err != nil {
		return fmt.Errorf("Error deleting Dns Domain: %s", err)
	}

	if result {
		return fmt.Errorf("Error deleting Dns Domain")
	}

	d.SetId("")
	return nil
}