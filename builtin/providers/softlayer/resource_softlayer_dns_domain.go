package softlayer

import (
	"fmt"
	"strconv"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
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
				Computed:	true,
			},

			"update_date": &schema.Schema{
				Type:		schema.TypeString,
				Computed:	true,
			},

			"records": &schema.Schema{
				Type:		schema.TypeList,
				Computed:	true,
				Elem:		&schema.Resource{
					Schema: get_dns_domain_record_scheme(),
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

	// create Dns_Domain object
	response, err := client.CreateObject(opts)
	if err != nil {
		return fmt.Errorf("Error creating Dns Domain: %s", err)
	}

	// populate id
	id := response.Id
	d.SetId(strconv.Itoa(id))
	log.Printf("[INFO] Created Dns Domain: %d", id)

	// read remote state
	return resourceSoftLayerDnsDomainRead(d, meta)
}

func resourceSoftLayerDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainService

	dnsId, _ := strconv.Atoi(d.Id())

	// retrieve remote object state
	dns_domain, err := client.GetObject(dnsId)
	if err != nil {
		return fmt.Errorf("Error retrieving Dns Domain %d: %s", dnsId, err)
	}

	// populate fields
	d.Set("id", dns_domain.Id)
	d.Set("name", dns_domain.Name)
	d.Set("serial", dns_domain.Serial)
	d.Set("update_date", dns_domain.UpdateDate)
	d.Set("records", read_resource_records(dns_domain.ResourceRecords))

	return nil
}

func read_resource_records(list []datatypes.SoftLayer_Dns_Domain_Record) []map[string]interface{} {
	records := make([]map[string]interface{}, 0, len(list))
	for _,record := range list {
		r := make(map[string]interface{})
		r["record_data"] =	record.Data
		r["domain_id"] =	record.DomainId
		r["expire"] = 		record.Expire
		r["host"] = 		record.Host
		r["minimum_ttl"] = 	record.Minimum
		r["mx_priority"] =	record.MxPriority
		r["refresh"] = 		record.Refresh
		r["retry"] = 		record.Retry
		r["ttl"] = 			record.Ttl
		r["record_type"] =	record.Type
		records = append(records, r)
	}
	return records
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

	if !result {
		return fmt.Errorf("Error deleting Dns Domain")
	}

	d.SetId("")
	return nil
}