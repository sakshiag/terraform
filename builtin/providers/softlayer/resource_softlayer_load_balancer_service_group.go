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
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"routing_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_type": &schema.Schema{
				Type:     schema.TypeString,
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
		Allocation:    d.Get("allocation").(int),
		Port:          d.Get("port").(int),
		RoutingMethod: d.Get("routing_method").(string),
		RoutingType:   d.Get("routing_type").(string),
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
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	loadBalancer, err := client.GetObject(GetLbId(d.Id()))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	opts := softlayer.SoftLayer_Load_Balancer_Service_Group_CreateOptions{
		Allocation:    d.Get("allocation").(int),
		Port:          d.Get("port").(int),
		RoutingMethod: d.Get("routing_method").(string),
		RoutingType:   d.Get("routing_type").(string),
	}

	log.Printf("[INFO] Updating load balancer service group")

	success, err := client.UpdateLoadBalancerVirtualServer(loadBalancer.Id, GetServiceGroupId(d.Id()), &opts)

	if err != nil {
		return fmt.Errorf("Error updating load balancer service group: %s", err)
	}

	if !success {
		return fmt.Errorf("Error updating load balancer service group")
	}

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
			d.Set("routing_method", serviceGroup.RoutingMethod)
			d.Set("routing_type", serviceGroup.RoutingType)
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
	client := meta.(*Client).loadBalancerService
	lb, err := client.GetObject(GetLbId(d.Id()))
	if err != nil {
		return false, err
	}

	for _, virtualServer := range lb.VirtualServers {
		if virtualServer.ServiceGroups[0].Id == GetServiceGroupId(d.Id()) {
			return true, nil
		}
	}

	return false, nil
}

func GetLbId(id string) int {
	lbId, err := strconv.Atoi(strings.Split(id, "|")[0])

	if err != nil {
		return -1
	}

	return lbId
}

func GetServiceGroupId(id string) int {
	lbId, err := strconv.Atoi(strings.Split(id, "|")[1])

	if err != nil {
		return -1
	}

	return lbId
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
