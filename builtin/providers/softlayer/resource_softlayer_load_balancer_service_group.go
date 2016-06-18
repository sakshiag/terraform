package softlayer

import (
	"fmt"
	"log"

	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
)

func resourceSoftLayerLoadBalancerServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerLoadBalancerServiceGroupCreate,
		Read:   resourceSoftLayerLoadBalancerServiceGroupRead,
		Update: resourceSoftLayerLoadBalancerServiceGroupUpdate,
		Delete: resourceSoftLayerLoadBalancerServiceGroupDelete,
		Exists: resourceSoftLayerLoadBalancerServiceGroupExists,

		Schema: map[string]*schema.Schema{
			"virtual_server_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"allocation": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"routing_method_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"routing_type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceSoftLayerLoadBalancerServiceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	loadBalancer, err := client.GetObject(d.Get("load_balancer_id").(int))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	opts := softlayer.SoftLayer_Load_Balancer_Service_Group_CreateOptions{
		Allocation:      d.Get("allocation").(int),
		Port:            d.Get("port").(int),
		RoutingMethodId: d.Get("routing_method_id").(int),
		RoutingTypeId:   d.Get("routing_type_id").(int),
	}

	log.Printf("[INFO] Creating load balancer service group")

	success, err := client.CreateLoadBalancerVirtualServer(loadBalancer.Id, &opts)

	if err != nil {
		return fmt.Errorf("Error creating load balancer service group: %s", err)
	}

	if !success {
		return fmt.Errorf("Error creating load balancer service group")
	}

	loadBalancer, err = client.GetObject(d.Get("load_balancer_id").(int))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	for _, virtualServer := range loadBalancer.VirtualServers {
		if virtualServer.Port == d.Get("port").(int) {
			d.SetId(fmt.Sprintf("%d|%d", loadBalancer.Id, virtualServer.ServiceGroups[0].Id))
		}
	}

	log.Printf("[INFO] Load Balancer Service Group ID: %s", d.Id())

	return resourceSoftLayerLoadBalancerServiceGroupRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSoftLayerLoadBalancerServiceGroupRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	id, err := strconv.Atoi(strings.Split(d.Id(), "|")[1])
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	loadBalancer, err := client.GetObject(d.Get("load_balancer_id").(int))
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	for _, virtualServer := range loadBalancer.VirtualServers {
		serviceGroup := virtualServer.ServiceGroups[0]
		if serviceGroup.Id == id {
			d.Set("virtual_server_id", virtualServer.Id)
			d.Set("allocation", virtualServer.Allocation)
			d.Set("port", virtualServer.Port)
			d.Set("routing_method_id", serviceGroup.RoutingMethodId)
			d.Set("routing_type_id", serviceGroup.RoutingTypeId)
		}
	}

	return nil
}

func resourceSoftLayerLoadBalancerServiceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	success, err := client.DeleteLoadBalancerVirtualServer(d.Get("virtual_server_id").(int))

	if err != nil {
		return fmt.Errorf("Error deleting service group: %s", err)
	}

	if !success {
		return fmt.Errorf("Error deleting service group")
	}

	return nil
}

func resourceSoftLayerLoadBalancerServiceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func getLbId(id string) int {
	lbId, err := strconv.Atoi(strings.Split(id, "|")[0])

	if err != nil {
		return -1
	}

	return lbId
}

func getServiceGroupId(id string) int {
	lbId, err := strconv.Atoi(strings.Split(id, "|")[1])

	if err != nil {
		return -1
	}

	return lbId
}
