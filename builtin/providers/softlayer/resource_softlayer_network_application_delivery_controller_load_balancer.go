package softlayer

import (
	"fmt"
	"log"
	"strconv"

	"bytes"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/hashcode"
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
		Update: resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerUpdate,
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
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"virtual_server": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceVirtualServerHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocation": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"service_group": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Set:      resourceServiceGroupHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"routing_method_id": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"routing_type_id": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"service": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Set:      resourceServiceHash,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
												"ip_address_id": &schema.Schema{
													Type:     schema.TypeInt,
													Required: true,
												},
												"port": &schema.Schema{
													Type:     schema.TypeInt,
													Required: true,
												},
												"health_check_type_id": &schema.Schema{
													Type:     schema.TypeInt,
													Required: true,
												},
												"weight": &schema.Schema{
													Type:     schema.TypeInt,
													Required: true,
												},
												"enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerLoadBalancerService

	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	virtualServers, err := expandVirtualServers(d.Get("virtual_server").(*schema.Set).List())

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer info: %s", err)
	}

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer info: %s", err)
	}

	success, err := client.UpdateLoadBalancer(id, virtualServers)

	if err != nil {
		return fmt.Errorf("Error updating load balancer info: %s", err)
	}

	if !success {
		return fmt.Errorf("Error updating load balancer info")
	}

	return resourceSoftLayerNetworkApplicationDeliveryControllerLoadBalancerRead(d, meta)
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
	d.Set("ip_address", loadBalancer.IpAddress.IpAddress)
	d.Set("subnet_id", loadBalancer.IpAddress.SubnetId)

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
	d.Set("ip_address", getObjectResult.IpAddress.IpAddress)
	d.Set("subnet_id", getObjectResult.IpAddress.SubnetId)
	d.Set("virtual_server", flattenVirtualServers(getObjectResult.VirtualServers))

	log.Println("[TRACE] VIRTUAL_SERVER_RESULTS")
	log.Println("[TRACE] %s", flattenVirtualServers(getObjectResult.VirtualServers))

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

/* When requesting 15000 SL creates between 15000 and 150000. When requesting 150000 SL creates >= 150000 */
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

func expandVirtualServers(configured []interface{}) ([]*datatypes.Softlayer_Load_Balancer_Virtual_Server, error) {
	virtualServers := make([]*datatypes.Softlayer_Load_Balancer_Virtual_Server, 0, len(configured))

	for _, rawVirtualServer := range configured {
		dataVirtualServer := rawVirtualServer.(map[string]interface{})

		serviceGroups := make([]*datatypes.Softlayer_Service_Group, 0, len(dataVirtualServer["service_group"].(*schema.Set).List()))

		for _, rawServiceGroup := range dataVirtualServer["service_group"].(*schema.Set).List() {
			dataServiceGroup := rawServiceGroup.(map[string]interface{})

			services := make([]*datatypes.Softlayer_Service, 0, len(dataServiceGroup["service"].(*schema.Set).List()))

			for _, rawService := range dataServiceGroup["service"].(*schema.Set).List() {
				dataService := rawService.(map[string]interface{})

				service := &datatypes.Softlayer_Service{
					Enabled:         btoi(dataService["enabled"].(bool)),
					Port:            dataService["port"].(int),
					IpAddressId:     dataService["ip_address_id"].(int),
					HealthChecks:    []*datatypes.Softlayer_Health_Check{{HealthCheckTypeId: dataService["health_check_type_id"].(int)}},
					GroupReferences: []*datatypes.Softlayer_Group_Reference{{Weight: dataService["weight"].(int)}},
				}

				services = append(services, service)
			}

			serviceGroup := &datatypes.Softlayer_Service_Group{
				RoutingMethodId: dataServiceGroup["routing_method_id"].(int),
				RoutingTypeId:   dataServiceGroup["routing_type_id"].(int),
				Services:        services,
			}

			serviceGroups = append(serviceGroups, serviceGroup)
		}

		virtualServerId, found := dataVirtualServer["name"].(int)

		if found {
			virtualServer := &datatypes.Softlayer_Load_Balancer_Virtual_Server{
				Allocation:    dataVirtualServer["allocation"].(int),
				Port:          dataVirtualServer["port"].(int),
				Id:            virtualServerId,
				ServiceGroups: serviceGroups,
			}

			virtualServers = append(virtualServers, virtualServer)
		} else {
			virtualServer := &datatypes.Softlayer_Load_Balancer_Virtual_Server{
				Allocation:    dataVirtualServer["allocation"].(int),
				Port:          dataVirtualServer["port"].(int),
				ServiceGroups: serviceGroups,
			}

			virtualServers = append(virtualServers, virtualServer)
		}
	}

	return virtualServers, nil
}

func flattenVirtualServers(list []*datatypes.Softlayer_Load_Balancer_Virtual_Server) []map[string]interface{} {
	finalResult := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		serviceGroupsResult := make([]map[string]interface{}, 0, len(i.ServiceGroups))
		for _, j := range i.ServiceGroups {
			servicesResults := make([]map[string]interface{}, 0, len(j.Services))
			for _, k := range j.Services {
				service := map[string]interface{}{
					"name":                 k.Id,
					"enabled":              k.Enabled,
					"port":                 k.Port,
					"ip_address_id":        k.IpAddressId,
					"health_check_type_id": k.HealthChecks[0].HealthCheckTypeId,
					"weight":               k.GroupReferences[0].Weight,
				}

				servicesResults = append(servicesResults, service)
			}

			serviceGroup := map[string]interface{}{
				"name":              j.Id,
				"routing_method_id": j.RoutingMethodId,
				"routing_type_id":   j.RoutingTypeId,
				"service":           servicesResults,
			}

			serviceGroupsResult = append(serviceGroupsResult, serviceGroup)
		}

		virtualServer := map[string]interface{}{
			"name":          i.Id,
			"allocation":    i.Allocation,
			"port":          i.Port,
			"service_group": serviceGroupsResult,
		}

		finalResult = append(finalResult, virtualServer)
	}

	return finalResult
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func resourceVirtualServerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["allocation"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))

	//if v, ok := m["service_group"]; ok {
	//	vs := v.(*schema.Set).List()
	//	for _, rawServiceGroup := range vs {
	//		buf.WriteString(fmt.Sprintf("%d-", rawServiceGroup.(map[string]interface{})["routing_method_id"].(int)))
	//		buf.WriteString(fmt.Sprintf("%d-", rawServiceGroup.(map[string]interface{})["routing_type_id"].(int)))
	//		if v1, ok := rawServiceGroup.(map[string]interface{})["service"]; ok {
	//			vs1 := v1.(*schema.Set).List()
	//			for _, rawService := range vs1 {
	//				buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["ip_address_id"].(int)))
	//				buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["port"].(int)))
	//				buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["health_check_type_id"].(int)))
	//				buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["weight"].(int)))
	//				buf.WriteString(fmt.Sprintf("%t-", rawService.(map[string]interface{})["enabled"].(bool)))
	//			}
	//		}
	//	}
	//}

	return hashcode.String(buf.String())
}

func resourceServiceGroupHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["routing_method_id"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["routing_type_id"].(int)))

	//if v, ok := m["service"]; ok {
	//	vs := v.(*schema.Set).List()
	//	for _, rawService := range vs {
	//		buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["ip_address_id"].(int)))
	//		buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["port"].(int)))
	//		buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["health_check_type_id"].(int)))
	//		buf.WriteString(fmt.Sprintf("%d-", rawService.(map[string]interface{})["weight"].(int)))
	//		buf.WriteString(fmt.Sprintf("%t-", rawService.(map[string]interface{})["enabled"].(bool)))
	//	}
	//}

	return hashcode.String(buf.String())
}

func resourceServiceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["ip_address_id"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["health_check_type_id"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["weight"].(int)))
	buf.WriteString(fmt.Sprintf("%t-", m["enabled"].(bool)))

	return hashcode.String(buf.String())
}
