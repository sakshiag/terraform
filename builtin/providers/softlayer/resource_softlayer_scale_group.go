package softlayer

import (
	"fmt"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
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
				Type: schema.TypeMap,
				// TODO Make this required once the softlayer-go datatype is available
				// Until then, use health_check_id
				Optional: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
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

	// Build up our creation options
	opts := datatypes.SoftLayer_Scale_Group{
		Name:                       d.Get("name").(string),
		Cooldown:                   d.Get("cooldown").(int),
		MinimumMemberCount:         d.Get("minimum_member_count").(int),
		MaximumMemberCount:         d.Get("maximum_member_count").(int),
		SuspendedFlag:              false,
		VirtualGuestMemberTemplate: virtualGuestTemplateOpts,
	}

	opts.RegionalGroup = &datatypes.SoftLayer_Location_Group_Regional{
		Name: d.Get("regional_group").(string),
	}

	opts.TerminationPolicy = &datatypes.SoftLayer_Scale_Termination_Policy{
		KeyName: d.Get("termination_policy").(string),
	}

	healthCheck := d.Get("health_check").(map[string]interface{})
	healthCheckOpts := datatypes.SoftLayer_Health_Check{
		Name: healthCheck["type"].(string),
	}

	if healthCheckOpts.Name == "HTTP-CUSTOM" {
		// Validate and apply type-specific fields
		healthCheckMethod, ok := healthCheck["custom_method"]
		if !ok {
			return fmt.Errorf("\"custom_method\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckRequest, ok := healthCheck["custom_request"]
		if !ok {
			return fmt.Errorf("\"custom_request\" is required when HTTP-CUSTOM healthcheck is specified")
		}

		healthCheckResponse, ok := healthCheck["custom_response"]
		if !ok {
			return fmt.Errorf("\"custom_response\" is required when HTTP-CUSTOM healthcheck is specified")
		}

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
	//client := meta.(*Client).scaleGroupService

	return nil
}

func resourceSoftLayerScaleGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client).scaleGroupService

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
	//client := meta.(*Client).scaleGroupService

	return true, nil
}
