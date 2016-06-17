package softlayer

import (
	"fmt"
	"log"

	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
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
			"service_group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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

	//d.SetId(fmt.Sprintf("%d", loadBalancerServiceGroup.ServiceGroups[0].Id))

	log.Printf("[INFO] Load Balancer Service Group ID: %s", d.Id())

	return resourceSoftLayerLoadBalancerRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceSoftLayerLoadBalancerRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSoftLayerLoadBalancerServiceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSoftLayerLoadBalancerServiceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}
