package softlayer

import (
	"fmt"
	"github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSoftLayerLoadBalancerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerLoadBalancerServiceCreate,
		Read:   resourceSoftLayerLoadBalancerServiceRead,
		Update: resourceSoftLayerLoadBalancerServiceUpdate,
		Delete: resourceSoftLayerLoadBalancerServiceDelete,
		Exists: resourceSoftLayerLoadBalancerServiceExists,

		Schema: map[string]*schema.Schema{
			"service_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"health_check_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceSoftLayerLoadBalancerServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	loadBalancer, err := client.GetObject(GetLbId(d.Get("service_group_id").(string)))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	opts := softlayer.SoftLayer_Load_Balancer_Service_CreateOptions{
		ServiceGroupId:  GetServiceGroupId(d.Get("service_group_id").(string)),
		Enabled:         1,
		Port:            d.Get("port").(int),
		IpAddressId:     d.Get("ip_address_id").(int),
		HealthCheckType: d.Get("health_check_type").(string),
		Weight:          d.Get("weight").(int),
	}

	log.Printf("[INFO] Creating load balancer service")

	success, err := client.CreateLoadBalancerService(loadBalancer.Id, &opts)

	if err != nil {
		return fmt.Errorf("Error creating load balancer service: %s", err)
	}

	if !success {
		return fmt.Errorf("Error creating load balancer service")
	}

	loadBalancer, err = client.GetObject(GetLbId(d.Get("service_group_id").(string)))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	for _, virtualServer := range loadBalancer.VirtualServers {
		if virtualServer.Id == virtualServer.Id {
			for _, service := range virtualServer.ServiceGroups[0].Services {
				if service.IpAddressId == d.Get("ip_address_id").(int) &&
					service.Port == d.Get("port").(int) {
					d.SetId(fmt.Sprintf("%d|%d|%d", loadBalancer.Id, virtualServer.ServiceGroups[0].Id, service.Id))
				}
			}
		}
	}

	log.Printf("[INFO] Load Balancer Service ID: %s", d.Id())

	return resourceSoftLayerLoadBalancerServiceRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	loadBalancer, err := client.GetObject(GetLbId(d.Id()))

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	opts := softlayer.SoftLayer_Load_Balancer_Service_CreateOptions{
		ServiceGroupId:  GetServiceGroupId(d.Id()),
		Enabled:         1,
		Port:            d.Get("port").(int),
		IpAddressId:     d.Get("ip_address_id").(int),
		HealthCheckType: d.Get("health_check_type").(string),
		Weight:          d.Get("weight").(int),
	}

	log.Printf("[INFO] Updating load balancer service")

	success, err := client.UpdateLoadBalancerService(loadBalancer.Id, GetServiceGroupId(d.Id()), getServiceId(d.Id()), &opts)

	if err != nil {
		return fmt.Errorf("Error updating load balancer service: %s", err)
	}

	if !success {
		return fmt.Errorf("Error updating load balancer service")
	}

	return resourceSoftLayerLoadBalancerServiceRead(d, meta)
}

func resourceSoftLayerLoadBalancerServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	loadBalancer, err := client.GetObject(GetLbId(d.Id()))
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	for _, virtualServer := range loadBalancer.VirtualServers {
		serviceGroup := virtualServer.ServiceGroups[0]
		if serviceGroup.Id == GetServiceGroupId(d.Id()) {
			for _, service := range serviceGroup.Services {
				if service.Id == getServiceId(d.Id()) {
					d.Set("ip_address_id", service.IpAddressId)
					d.Set("port", service.Port)
					d.Set("health_check_type", service.HealthChecks[0].HealthCheckType)
					d.Set("weight", service.GroupReferences[0].Weight)
				}
			}
		}
	}

	return nil
}

func resourceSoftLayerLoadBalancerServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).loadBalancerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	success, err := client.DeleteLoadBalancerService(getServiceId(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting service group: %s", err)
	}

	if !success {
		return fmt.Errorf("Error deleting service group")
	}

	return nil
}

func resourceSoftLayerLoadBalancerServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).loadBalancerService
	lb, err := client.GetObject(GetLbId(d.Id()))
	if err != nil {
		return false, err
	}

	for _, virtualServer := range lb.VirtualServers {
		if virtualServer.ServiceGroups[0].Id == GetServiceGroupId(d.Id()) {
			for _, service := range virtualServer.ServiceGroups[0].Services {
				if service.Id == getServiceId(d.Id()) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func getServiceId(id string) int {
	lbId, err := strconv.Atoi(strings.Split(id, "|")[2])

	if err != nil {
		return -1
	}

	return lbId
}
