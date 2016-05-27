package softlayer

import (
	"fmt"
	"log"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/TheWeatherCompany/softlayer-go/services"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
)

func resourceSoftLayerNetworkLoadBalancerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerNetworkLoadBalancerServiceCreate,
		Read:   resourceSoftLayerNetworkLoadBalancerServiceRead,
		Delete: resourceSoftLayerNetworkLoadBalancerServiceDelete,
		Exists: resourceSoftLayerNetworkLoadBalancerServiceExists,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"vip_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func parseVipUniqueId(vipUniqueId string) (vipId string, nacdId int, err error) {
	nacdId, err = strconv.Atoi(strings.Split(vipUniqueId, services.ID_DELIMITER)[1])
	vipId = strings.Split(vipUniqueId, services.ID_DELIMITER)[0]

	if err != nil {
		return "", -1, fmt.Errorf("Error parsing vip id: %s", err)
	}

	return vipId, nacdId, nil
}

func resourceSoftLayerNetworkLoadBalancerServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService

	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	vipUniqueId := d.Get("vip_id").(string)

	vipId, nacdId, err := parseVipUniqueId(vipUniqueId)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	template := []datatypes.SoftLayer_Network_LoadBalancer_Service_Template{datatypes.SoftLayer_Network_LoadBalancer_Service_Template{
		Name:                 d.Get("name").(string),
		DestinationIpAddress: d.Get("destination_ip_address").(string),
		DestinationPort:      d.Get("destination_port").(int),
		Weight:               d.Get("weight").(int),
	}}

	log.Printf("[INFO] Creating LoadBalancer Service %s", template[0].Name)

	successFlag, err := client.CreateLoadBalancerService(vipId, nacdId, template)

	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service: %s", err)
	}

	if !successFlag {
		return fmt.Errorf("Error creating LoadBalancer Service")
	}

	return resourceSoftLayerNetworkLoadBalancerServiceRead(d, meta)
}

func resourceSoftLayerNetworkLoadBalancerServiceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	vipUniqueId := d.Get("vip_id").(string)

	vipId, nadcId, err := parseVipUniqueId(vipUniqueId)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	service, err := client.GetLoadBalancerService(nadcId, vipId, d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Unable to get LoadBalancerService: %s", err)
	}

	d.SetId(service.Name)
	d.Set("name", service.Name)
	d.Set("destination_ip_address", service.DestinationIpAddress)
	d.Set("destination_port", service.DestinationPort)
	d.Set("weight", service.Weight)

	return nil
}

func resourceSoftLayerNetworkLoadBalancerServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	vipUniqueId := d.Get("vip_id").(string)

	vipId, nadcId, err := parseVipUniqueId(vipUniqueId)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	serviceId := d.Get("name").(string)

	_, err = client.DeleteLoadBalancerService(nadcId, vipId, serviceId)
	if err != nil {
		return fmt.Errorf("Error deleting Load Balancer Service %s: %s", serviceId, err)
	}

	return nil
}

func resourceSoftLayerNetworkLoadBalancerServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return false, fmt.Errorf("The client is nil.")
	}

	vipUniqueId := d.Get("vip_id").(string)

	vipId, nadcId, err := parseVipUniqueId(vipUniqueId)

	if err != nil {
		return false, fmt.Errorf("Error parsing vip id: %s", err)
	}

	serviceId := d.Get("name").(string)

	service, err := client.GetLoadBalancerService(nadcId, vipId, serviceId)

	if err != nil {
		return false, fmt.Errorf("Error fetching Load Balancer Service: %s", err)
	}

	return service.Name == serviceId && err == nil, nil
}
