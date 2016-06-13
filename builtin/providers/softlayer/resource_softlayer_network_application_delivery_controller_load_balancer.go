package softlayer

import (
	"fmt"
	"log"
	"strconv"

	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	LB_LARGE_150000_CONNECTIONS = 150000
	LB_SMALL_15000_CONNECTIONS  = 15000
)

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerCreate,
		Read:   resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerRead,
		Delete: resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerDelete,
		Exists: resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerExists,

		Schema: map[string]*schema.Schema{
			"connections": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerLoadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	opts := softlayer.SoftLayer_Network_Application_Delivery_Controller_Load_Balancer_Service_CreateOptions{
		Connections: d.Get("connections").(int),
		Location:    d.Get("location").(string),
	}

	log.Printf("[INFO] Creating load balancer")

	loadBalancer, err := client.CreateLoadBalancer(&opts)

	if err != nil {
		return fmt.Errorf("Error creating load balancer: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", loadBalancer.Id))
	d.Set("connections", getConnectionLimit(loadBalancer.ConnectionLimit))
	d.Set("location", loadBalancer.SoftlayerHardware[0].Datacenter.Name)

	log.Printf("[INFO] Load Balancer ID: %s", d.Id())

	return resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerRead(d, meta)
}

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerLoadBalancerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	getObjectResult, err := client.GetObject(id)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	d.SetId(strconv.Itoa(getObjectResult.Id))
	d.Set("connections", getConnectionLimit(getObjectResult.ConnectionLimit))
	d.Set("location", getObjectResult.SoftlayerHardware[0].Datacenter.Name)

	return nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerLoadBalancerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = client.DeleteObject(id)

	if err != nil {
		return fmt.Errorf("Error deleting network application delivery controller load balancer: %s", err)
	}

	return nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func getConnectionLimit(connectionLimit int) int {
	if connectionLimit >= LB_LARGE_150000_CONNECTIONS {
		return LB_LARGE_150000_CONNECTIONS
	} else if connectionLimit >= LB_SMALL_15000_CONNECTIONS &&
		connectionLimit < LB_LARGE_150000_CONNECTIONS {
		return LB_SMALL_15000_CONNECTIONS
	} else {
		return 0
	}
}
