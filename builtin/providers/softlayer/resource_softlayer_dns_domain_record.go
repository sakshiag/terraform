package softlayer

import (
	"fmt"
	"log"
	"strconv"
//	"time"

//	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
)

func resourceSoftLayerDnsDomainResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerDnsDomainResourceRecordCreate,
		Read: resourceSoftLayerDnsDomainResourceRecordRead,
		Update: resourceSoftLayerDnsDomainResourceRecordUpdate,
		Delete: resourceSoftLayerDnsDomainResourceRecordDelete,
		Schema: get_dns_domain_record_scheme(),
	}
}


//  Creates DNS Domain Resource Record
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/createObject
func resourceSoftLayerDnsDomainResourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainResourceRecordService
	if client == nil {
		return fmt.Errorf("The client was nil.")
	}

	opts := datatypes.SoftLayer_Dns_Domain_Record_Template{
		Data: d.Get("record_data").(string),
		DomainId: d.Get("domain_id").(int),
		Expire: d.Get("expire").(int),
		Host: d.Get("host").(string),
		Minimum: d.Get("minimum_ttl").(int),
		MxPriority: d.Get("mx_priority").(int),
		Refresh: d.Get("refresh").(int),
		ResponsiblePerson: d.Get("contact_email").(string),
		Retry: d.Get("retry").(int),
		Ttl: d.Get("ttl").(int),
		Type: d.Get("record_type").(string),
	}

	log.Printf("[INFO] Creating DNS Resource Record for '%d' dns domain", d.Get("id"))

	record, err := client.CreateObject(opts)

	if err != nil {
		return fmt.Errorf("Error creating DNS Resource Record: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", record.Id))

	log.Printf("[INFO] Dns Resource Record ID: %s", d.Id())

	return resourceSoftLayerDnsDomainResourceRecordRead(d, meta)
}

//  Reads DNS Domain Resource Record from SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/getObject
func resourceSoftLayerDnsDomainResourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainResourceRecordService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	result, err := client.GetObject(id)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Resource Record: %s", err)
	}

	d.Set("data", result.Data)
	d.Set("domainId", result.DomainId)
	d.Set("expire", result.Expire)
	d.Set("host", result.Host)
	d.Set("id", result.Id)
	d.Set("minimum", result.Minimum)
	d.Set("mxPriority", result.MxPriority)
	d.Set("refresh", result.Refresh)
	d.Set("responsiblePerson", result.ResponsiblePerson)
	d.Set("retry", result.Retry)
	d.Set("ttl", result.Ttl)
	d.Set("type", result.Type)

	return nil
}


//  Updates DNS Domain Resource Record in SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/editObject
func resourceSoftLayerDnsDomainResourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainResourceRecordService

	recordId, _ := strconv.Atoi(d.Id())

	record, err := client.GetObject(recordId)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Resource Record: %s", err)
	}

	if data, ok := d.GetOk("record_data"); ok {
		record.Data = data.(string)
	}
	if domain_id, ok := d.GetOk("domain_id"); ok {
		record.DomainId = domain_id.(int)
	}
	if expire, ok := d.GetOk("expire"); ok {
		record.Expire = expire.(int)
	}
	if host, ok := d.GetOk("host"); ok {
		record.Host = host.(string)
	}
	if minimum_ttl, ok := d.GetOk("minimum_ttl"); ok {
		record.Minimum = minimum_ttl.(int)
	}
	if mx_priority, ok := d.GetOk("mx_priority"); ok {
		record.MxPriority = mx_priority.(int)
	}
	if refresh, ok := d.GetOk("refresh"); ok {
		record.Refresh = refresh.(int)
	}
	if contact_email, ok := d.GetOk("contact_email"); ok {
		record.ResponsiblePerson = contact_email.(string)
	}
	if retry, ok := d.GetOk("retry"); ok {
		record.Retry = retry.(int)
	}
	if ttl, ok := d.GetOk("ttl"); ok {
		record.Ttl = ttl.(int)
	}
	if record_type, ok := d.GetOk("record_type"); ok {
		record.Type = record_type.(string)
	}

	_, err = client.EditObject(recordId, record)
	if err != nil {
		return fmt.Errorf("Error editing DNS Resoource Record: %s", err)
	}
	return nil
}

//  Deletes DNS Domain Resource Record in SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/deleteObject
func resourceSoftLayerDnsDomainResourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("DeleteDNS Record!")
	client := meta.(*Client).dnsDomainResourceRecordService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = client.DeleteObject(id)

	if err != nil {
		return fmt.Errorf("Error deleting DNS Resource Record: %s", err)
	}

	return nil
}

// the scheme is used by dns_domain_service,
// so it's creation is extracted to a separate function
func get_dns_domain_record_scheme()  map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

		"contact_email": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"retry": &schema.Schema{
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
	}
}