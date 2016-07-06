package softlayer

import (
	"fmt"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSoftLayerScaleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerScaleGroupCreate,
		Read:   resourceSoftLayerScaleGroupRead,
		Update: resourceSoftLayerScaleGroupUpdate,
		Delete: resourceSoftLayerScaleGroupDelete,
		Exists: resourceSoftLayerScaleGroupExists,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"regional_group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"minimum_member_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"maximum_member_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"cooldown": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"termination_policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"virtual_server_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"health_check": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
						},

						// Conditionally-required fields, based on value of "type"
						"custom_method": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							// TODO: Must be GET or HEAD
						},

						"custom_request": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"custom_response": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// This has to be a TypeList, because TypeMap does not handle non-primitive
			// members properly.
			// TODO Validate that only one template is provided
			"virtual_guest_member_template": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     resourceSoftLayerVirtualGuest(),
			},

			"network_vlans": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_number": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"primary_router_hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Helper method to parse healthcheck data in the resource schema format to the SoftLayer datatypes
func buildHealthCheckFromResourceData(d map[string]interface{}) (datatypes.SoftLayer_Health_Check, error) {
	healthCheckOpts := datatypes.SoftLayer_Health_Check{
		Type: datatypes.SoftLayer_Health_Check_Type{
			KeyName: d["type"].(string),
		},
	}

	if healthCheckOpts.Type.KeyName == datatypes.HEALTH_CHECK_TYPE_HTTP_CUSTOM {
		// Validate and apply type-specific fields
		healthCheckMethod, ok := d["custom_method"]
		if !ok {
			return datatypes.SoftLayer_Health_Check{}, fmt.Errorf("\"custom_method\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckRequest, ok := d["custom_request"]
		if !ok {
			return datatypes.SoftLayer_Health_Check{}, fmt.Errorf("\"custom_request\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckResponse, ok := d["custom_response"]
		if !ok {
			return datatypes.SoftLayer_Health_Check{}, fmt.Errorf("\"custom_response\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		// HTTP-CUSTOM values are represented as an array of SoftLayer_Health_Check_Attributes
		healthCheckOpts.Attributes = make([]datatypes.SoftLayer_Health_Check_Attribute, 3)
		healthCheckOpts.Attributes[0] = datatypes.SoftLayer_Health_Check_Attribute{
			Type: &datatypes.SoftLayer_Health_Check_Attribute_Type{
				Keyname: "HTTP_CUSTOM_TYPE",
			},
			Value: healthCheckMethod.(string),
		}
		healthCheckOpts.Attributes[1] = datatypes.SoftLayer_Health_Check_Attribute{
			Type: &datatypes.SoftLayer_Health_Check_Attribute_Type{
				Keyname: "LOCATION",
			},
			Value: healthCheckRequest.(string),
		}
		healthCheckOpts.Attributes[2] = datatypes.SoftLayer_Health_Check_Attribute{
			Type: &datatypes.SoftLayer_Health_Check_Attribute_Type{
				Keyname: "EXPECTED_RESPONSE",
			},
			Value: healthCheckResponse.(string),
		}

	}

	return healthCheckOpts, nil
}

// Helper method to parse network vlan information in the resource schema format to the SoftLayer datatypes
func buildScaleVlansFromResourceData(d *schema.Set, meta interface{}) ([]datatypes.SoftLayer_Scale_Network_Vlan, error) {
	client := meta.(*Client).accountService

	scaleNetworkVlans := make([]datatypes.SoftLayer_Scale_Network_Vlan, 0, d.Len())

	for _, elem := range d.List() {
		elem := elem.(map[string]interface{})

		vlanNumber, err := strconv.Atoi(elem["vlan_number"].(string))
		if err != nil {
			return nil, fmt.Errorf("Vlan number must be an integer: %s", elem["vlan_number"])
		}

		primaryRouterHostname := elem["primary_router_hostname"].(string)

		mask := []string{
			"id",
		}

		filter := fmt.Sprintf(
			"{\"networkVlans\":{\"primaryRouter\":{\"hostname\":{\"operation\":\"%s\"}},"+
				"\"vlanNumber\":{\"operation\":%d}}}",
			primaryRouterHostname,
			vlanNumber)

		networkVlan, err := client.GetNetworkVlans(mask, filter)

		if err != nil {
			return nil, fmt.Errorf("Error looking up Vlan: %s", err)
		}

		if len(networkVlan) < 1 {
			return nil, fmt.Errorf(
				"Unable to locate a vlan matching the provided router hostname and vlan number: %s/%s",
				primaryRouterHostname,
				vlanNumber)
		}

		scaleNetworkVlans = append(
			scaleNetworkVlans,
			datatypes.SoftLayer_Scale_Network_Vlan{
				NetworkVlanId: networkVlan[0].Id,
			})
	}

	return scaleNetworkVlans, nil
}

func resourceSoftLayerScaleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scaleGroupService

	// Retrieve the map of virtual_guest_member_template attributes
	// Note: Because 'virtual_guest_member_template' is defined using TypeList, a slice is returned.  We assert
	// that only one element exists, therefore we get the first element in the slice, which contains the actual
	// map we care about.
	vGuestMap := d.Get("virtual_guest_member_template").([]interface{})[0].(map[string]interface{})

	// Create an empty ResourceData instance for a SoftLayer_Virtual_Guest resource
	vGuestResourceData := resourceSoftLayerVirtualGuest().Data(nil)

	// For each item in the map, call Set on the ResourceData.  This handles
	// validation and yields a completed ResourceData object
	for k, v := range vGuestMap {
		log.Printf("****** %s: %#v", k, v)
		err := vGuestResourceData.Set(k, v)
		if err != nil {
			return fmt.Errorf("Error while parsing virtual_guest_member_template values: %s", err)
		}
	}

	// Get the virtual guest creation template from the completed resource data object
	virtualGuestTemplateOpts, _ := GetVirtualGuestTemplateFromResourceData(vGuestResourceData)

	scaleNetworkVlans, err := buildScaleVlansFromResourceData(d.Get("network_vlans").(*schema.Set), meta)
	if err != nil {
		return fmt.Errorf("Error while parsing network vlan values: %s", err)
	}

	// Build up our creation options
	opts := datatypes.SoftLayer_Scale_Group{
		Name:                       d.Get("name").(string),
		Cooldown:                   d.Get("cooldown").(int),
		MinimumMemberCount:         d.Get("minimum_member_count").(int),
		MaximumMemberCount:         d.Get("maximum_member_count").(int),
		SuspendedFlag:              false,
		VirtualGuestMemberTemplate: virtualGuestTemplateOpts,
		NetworkVlans:               scaleNetworkVlans,
	}

	opts.RegionalGroup = &datatypes.SoftLayer_Location_Group_Regional{
		Name: d.Get("regional_group").(string),
	}

	opts.TerminationPolicy = &datatypes.SoftLayer_Scale_Termination_Policy{
		KeyName: d.Get("termination_policy").(string),
	}

	healthCheckOpts, err := buildHealthCheckFromResourceData(d.Get("health_check").(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("Error while parsing health check options: %s", err)
	}

	opts.LoadBalancers = make([]datatypes.SoftLayer_Scale_LoadBalancer, 1)
	opts.LoadBalancers[0] = datatypes.SoftLayer_Scale_LoadBalancer{
		HealthCheck:     &healthCheckOpts,
		Port:            d.Get("port").(int),
		VirtualServerId: d.Get("virtual_server_id").(int),
	}

	res, err := client.CreateObject(opts)
	if err != nil {
		return fmt.Errorf("Error creating Scale Group: %s", err)
	}

	d.SetId(strconv.Itoa(res.Id))
	log.Printf("[INFO] Scale Group ID: %d", res.Id)

	return resourceSoftLayerScaleGroupRead(d, meta)
}

func resourceSoftLayerScaleGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scaleGroupService

	groupId, _ := strconv.Atoi(d.Id())

	slGroupObj, err := client.GetObject(groupId)
	if err != nil {
		// If the scale group is somehow already destroyed, mark as successfully gone
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving SoftLayer Scale Group: %s", err)
	}

	d.Set("id", slGroupObj.Id)
	d.Set("name", slGroupObj.Name)
	d.Set("regional_group", slGroupObj.RegionalGroup.Name)
	d.Set("minimum_member_count", slGroupObj.MinimumMemberCount)
	d.Set("maximum_member_count", slGroupObj.MaximumMemberCount)
	d.Set("cooldown", slGroupObj.Cooldown)
	d.Set("termination_policy", slGroupObj.TerminationPolicy.KeyName)
	d.Set("virtual_server_id", slGroupObj.LoadBalancers[0].VirtualServerId)
	d.Set("port", slGroupObj.LoadBalancers[0].Port)

	healthCheckObj := slGroupObj.LoadBalancers[0].HealthCheck
	currentHealthCheck := d.Get("health_check").(map[string]interface{})

	currentHealthCheck["type"] = healthCheckObj.Type.KeyName

	if healthCheckObj.Type.KeyName == "HTTP-CUSTOM" {
		for _, elem := range healthCheckObj.Attributes {
			switch elem.Type.Keyname {
			case "HTTP_CUSTOM_TYPE":
				currentHealthCheck["custom_method"] = elem.Value
			case "LOCATION":
				currentHealthCheck["custom_request"] = elem.Value
			case "EXPECTED_RESPONSE":
				currentHealthCheck["custom_response"] = elem.Value
			}
		}
	}

	d.Set("health_check", currentHealthCheck)

	log.Printf("******** healthCheckObj: %#v", healthCheckObj)
	log.Printf("******** currentHealthCheck: %#v", currentHealthCheck)

	return nil
}

func resourceSoftLayerScaleGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scaleGroupService

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID. Must be an integer: %s", err)
	}

	// Fetch the complete object from SoftLayer, update with current values from the configuration, and send the
	// whole thing back to SoftLayer (effectively, a PUT)
	groupObj, err := client.GetObject(groupId)
	if err != nil {
		return fmt.Errorf("Error retrieving softlayer_scale_group resource: %s", err)
	}

	groupObj.Name = d.Get("name").(string)
	groupObj.MinimumMemberCount = d.Get("minimum_member_count").(int)
	groupObj.MaximumMemberCount = d.Get("maximum_member_count").(int)
	groupObj.Cooldown = d.Get("cooldown").(int)
	groupObj.TerminationPolicy.KeyName = d.Get("termination_policy").(string)
	groupObj.LoadBalancers[0].VirtualServerId = d.Get("virtual_server_id").(int)
	groupObj.LoadBalancers[0].Port = d.Get("port").(int)

	healthCheck, err := buildHealthCheckFromResourceData(d.Get("health_check").(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("Unable to parse health check options: %s", err)
	}

	groupObj.LoadBalancers[0].HealthCheck = &healthCheck

	_, err = client.EditObject(groupId, groupObj)
	if err != nil {
		return fmt.Errorf("Error received while editing softlayer_scale_group: %s", err)
	}

	return nil
}

func resourceSoftLayerScaleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scaleGroupService

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting scale group: %s", err)
	}

	log.Printf("[INFO] Deleting scale group: %d", id)
	_, err = client.ForceDeleteObject(id)
	if err != nil {
		return fmt.Errorf("Error deleting scale group: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSoftLayerScaleGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).scaleGroupService

	if client == nil {
		return false, fmt.Errorf("The client was nil.")
	}

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := client.GetObject(groupId)
	return result.Id == groupId && err == nil, nil
}
