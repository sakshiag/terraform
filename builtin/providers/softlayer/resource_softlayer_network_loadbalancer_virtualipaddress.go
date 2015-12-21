package softlayer
import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressCreate,
		Read: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressRead,
		Update: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressUpdate,
		Delete: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressDelete,
		Exists: resourceSoftLayerNetworkLoadBalancerVirtualIpAddressExists,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"nad_controller_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"load_balancing_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"load_labancing_method_name": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"modify_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"virtual_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

	template := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template {
		ConnectionLimit: d.Get("connection_limit").(int),
		LoadBalancingMethod: d.Get("load_balancing_method").(string),
		Name: d.Get("name").(string),
		Notes: d.Get("notes").(string),
		SourcePort: d.Get("source_port").(int),
		Type: d.Get("type").(string),
		VirtualIpAddress: d.Get("virtual_ip_address").(string),
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

	d.SetId(fmt.Sprintf("%d", vip.Id))
	d.Set("nad_controller_id", nadcId)
	d.Set("connection_limit", vip.ConnectionLimit)
	d.Set("load_balancing_method", vip.LoadBalancingMethod)
	d.Set("load_labancing_method_name", vip.LoadBalancingMethodFullName)
	d.Set("modify_date", vip.ModifyDate)
	d.Set("name", vip.Name)
	d.Set("notes", vip.Notes)
	d.Set("source_port", vip.SourcePort)
	d.Set("type", vip.Type)
	d.Set("virtual_ip_address", vip.VirtualIpAddress)

	return nil
}

func resourceSoftLayerNetworkLoadBalancerVirtualIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	//	client := meta.(*Client).networkApplicationDeliveryControllerService
	return fmt.Errorf("Update is not supported yet")
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

