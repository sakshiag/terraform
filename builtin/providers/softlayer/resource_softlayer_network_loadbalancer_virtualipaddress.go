package softlayer

import (
	"fmt"
	"log"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressCreate,
		Read:   resourceSoftLayerNetworkLoadBalancerVirtualIpAddressRead,
		Delete: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressDelete,
		Exists: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressExists,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"nad_controller_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"load_balancing_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"load_labancing_method_name": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"modify_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// name field is actually used as an ID in SoftLayer
			// http://sldn.softlayer.com/reference/services/SoftLayer_Network_Application_Delivery_Controller/updateLiveLoadBalancer
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_certificate_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"source_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	nadcId := d.Get("nad_controller_id").(int)

	template := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
		LoadBalancingMethod:   d.Get("load_balancing_method").(string),
		Name:                  d.Get("name").(string),
		Notes:                 d.Get("notes").(string),
		SourcePort:            d.Get("source_port").(int),
		Type:                  d.Get("type").(string),
		VirtualIpAddress:      d.Get("virtual_ip_address").(string),
		SecurityCertificateId: d.Get("security_certificate_id").(int),
	}

	log.Printf("[INFO] Creating Virtual Ip Address %s", template.VirtualIpAddress)

	successFlag, err := client.CreateVirtualIpAddress(nadcId, template)

	if err != nil {
		return fmt.Errorf("Error creating Virtual Ip Address: %s", err)
	}

	if !successFlag {
		return fmt.Errorf("Error creating Virtual Ip Address")
	}

	return resourceSoftLayerNetworkLoadBalancerVirtualIpAddressRead(d, meta)
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	nadcId := d.Get("nad_controller_id").(int)
	vipName := d.Get("name").(string)

	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	vip, err := client.GetVirtualIpAddress(nadcId, vipName)
	if err != nil {
		return fmt.Errorf("Error getting Virtual Ip Address: %s", err)
	}

	d.SetId(vip.Name)
	d.Set("nad_controller_id", nadcId)
	d.Set("load_balancing_method", vip.LoadBalancingMethod)
	d.Set("load_labancing_method_name", vip.LoadBalancingMethodFullName)
	d.Set("modify_date", vip.ModifyDate)
	d.Set("name", vip.Name)
	d.Set("notes", vip.Notes)
	d.Set("security_certificate_id", vip.SecurityCertificateId)
	d.Set("source_port", vip.SourcePort)
	d.Set("type", vip.Type)
	d.Set("virtual_ip_address", vip.VirtualIpAddress)

	return nil
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	nadcId := d.Get("nad_controller_id").(int)
	template := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
		Name: d.Get("name").(string),
	}

	if d.HasChange("load_balancing_method") {
		template.LoadBalancingMethod = d.Get("load_balancing_method").(string)
	}

	if d.HasChange("notes") {
		template.Notes = d.Get("notes").(string)
	}

	if d.HasChange("security_certificate_id") {
		template.SecurityCertificateId = d.Get("security_certificate_id").(int)
	}

	if d.HasChange("source_port") {
		template.SourcePort = d.Get("source_port").(int)
	}

	if d.HasChange("type") {
		template.Type = d.Get("type").(string)
	}

	if d.HasChange("virtual_ip_address") {
		template.VirtualIpAddress = d.Get("virtual_ip_address").(string)
	}

	_, err := client.EditVirtualIpAddress(nadcId, template)

	if err != nil {
		return fmt.Errorf("Error updating Virtual Ip Address: %s", err)
	}

	return nil
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	nadcId := d.Get("nad_controller_id").(int)
	vipName := d.Get("name").(string)

	_, err := client.DeleteVirtualIpAddress(nadcId, vipName)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Ip Address %s: %s", vipName, err)
	}

	return nil
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return false, fmt.Errorf("The client is nil.")
	}

	vipName := d.Get("name").(string)
	nadcId := d.Get("nad_controller_id").(int)

	vip, err := client.GetVirtualIpAddress(nadcId, vipName)

	if err != nil {
		return false, fmt.Errorf("Error fetching Virtual Ip Address: %s", err)
	}

	return vip.Name == vipName && err == nil, nil
}
