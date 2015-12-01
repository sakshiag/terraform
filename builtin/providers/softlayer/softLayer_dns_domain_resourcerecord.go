package softlayer

import (
	"fmt"
	"log"
	"strconv"
//	"time"

//	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
//	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
)

func resourceSoftLayerDnsDomainResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerDnsDomainResourceRecordCreate,
		Read: resourceSoftLayerDnsDomainResourceRecordRead,
		Update: resourceSoftLayerDnsDomainResourceRecordUpdate,
		Delete: resourceSoftLayerDnsDomainResourceRecordDelete,
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
		},
	}
}

/**
  Creates DNS Domain Resource Record
  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/createObject
 */
func resourceSoftLayerDnsDomainResourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainResourceRecord
	if client == nil {
		return fmt.Errorf("The client was nil.")
	}

	log.Println("CreateDNS Record!")
	//	guest, err := client.CreateObject(nil)


	return resourceSoftLayerDnsDomainResourceRecordRead(d, meta)
}

/**
  Creates array of DNS Domain Resource Records
  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/createObjects
 */
func resourceSoftLayerDnsDomainResourceRecordsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).dnsDomainResourceRecord
	if client == nil {
		return fmt.Errorf("The client was nil.")
	}

	log.Println("CreateDNS Record!")
	//	guest, err := client.CreateObject(nil)

	//TODO will return several records
	return resourceSoftLayerDnsDomainResourceRecordRead(d, meta)
}

/**
  Reads DNS Domain Resource Record from SL system
  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/getObject
 */
func resourceSoftLayerDnsDomainResourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("ReadDNS Record!")
	client := meta.(*Client).dnsDomainResourceRecord
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	result, err := client.GetObject(id)
	if err != nil {
		return fmt.Errorf("Error retrieving dns domain resource record: %s", err)
	}

	d.Set("name", result.Host)
	return nil
}

/**
  Updates DNS Domain Resource Record in SL system
  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/editObject
 */
func resourceSoftLayerDnsDomainResourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	//TODO
	return nil
}

/**

 */
func resourceSoftLayerDnsDomainResourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("DeleteDNS Record!")
	client := meta.(*Client).dnsDomainResourceRecord
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = client.DeleteObject(id)

	if err != nil {
		return fmt.Errorf("Error deleting dns domain resource record: %s", err)
	}

	return nil
}

