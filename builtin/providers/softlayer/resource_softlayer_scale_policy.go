package softlayer

import (
	"fmt"
	"log"
	"strconv"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceSoftLayerScalePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerScalePolicyCreate,
		Read:   resourceSoftLayerScalePolicyRead,
		Update: resourceSoftLayerScalePolicyUpdate,
		Delete: resourceSoftLayerScalePolicyDelete,
		Exists: resourceSoftLayerScalePolicyExists,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scale_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scale_amount": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"cooldown": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"triggers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						// Conditionally-required fields, based on value of "type"
						"watches": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:	  schema.TypeInt,
										Computed: true,
									},
									"metric": &schema.Schema{
										Type:	  schema.TypeString,
										Required: true,
									},
									"operator": &schema.Schema{
										Type:	  schema.TypeString,
										Required: true,
									},
									"value": &schema.Schema{
										Type:	  schema.TypeString,
										Required: true,
									},
									"period": &schema.Schema{
										Type:	  schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"date": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"schedule": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

		},
	}
}

func resourceSoftLayerScalePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scalePolicyService

	// Build up our creation options
	opts := datatypes.SoftLayer_Scale_Policy{
		Name:         d.Get("name").(string),
		ScaleGroupId: d.Get("scale_group_id").(int),
		Cooldown:     d.Get("cooldown").(int),
	}

	opts.ScaleActions = []datatypes.SoftLayer_Scale_Policy_Action{{
		TypeId:       1,
		Amount:       d.Get("scale_amount").(int),
		ScaleType:    d.Get("scale_type").(string),
	},
	}

	if triggers, ok := d.GetOk("triggers"); ok {
		opts.OneTimeTriggers = prepareOneTimeTriggers(triggers.([]interface{}))
		opts.RepeatingTriggers = prepareRepeatingTriggers(triggers.([]interface{}))
		opts.ResourceUseTriggers = prepareResourceUseTriggers(triggers.([]interface{}))
	}

	res, err := client.CreateObject(opts)
	if err != nil {
		return fmt.Errorf("Error creating Scale Policy: %s", err)
	}

	d.SetId(strconv.Itoa(res.Id))
	log.Printf("[INFO] Scale Polocy: %d", res.Id)

	return resourceSoftLayerScalePolicyRead(d, meta)
}

func resourceSoftLayerScalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client).scalePolicyService

	return nil
}

func resourceSoftLayerScalePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*Client).scalePolicyService

	return nil
}

func resourceSoftLayerScalePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scalePolicyService

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting scale policy: %s", err)
	}

	log.Printf("[INFO] Deleting scale policy: %d", id)
	_, err = client.DeleteObject(id)
	if err != nil {
		return fmt.Errorf("Error deleting scale policy: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceSoftLayerScalePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	//client := meta.(*Client).scaleGroupService

	return true, nil
}

func prepareOneTimeTriggers(raw_triggers []interface{}) []datatypes.SoftLayer_Scale_Policy_Trigger_OneTime {
	sl_onetime_triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_OneTime, 0)
	for _, raw_trigger := range raw_triggers {
		trigger := raw_trigger.(map[string]interface{})

		if trigger["type"].(string) == "ONE_TIME" {
			var sl_onetime_trigger datatypes.SoftLayer_Scale_Policy_Trigger_OneTime
			sl_onetime_trigger.TypeId = 3
			timeStampString := trigger["date"].(string)
			timeStamp, _ := time.Parse(time.RFC3339Nano, timeStampString)
			sl_onetime_trigger.Date = &timeStamp
			sl_onetime_triggers = append(sl_onetime_triggers, sl_onetime_trigger)
		}
	}
	return sl_onetime_triggers
}

func prepareRepeatingTriggers(raw_triggers []interface{}) []datatypes.SoftLayer_Scale_Policy_Trigger_Repeating {
	sl_repeating_triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_Repeating, 0)
	for _, raw_trigger := range raw_triggers {
		trigger := raw_trigger.(map[string]interface{})

		if trigger["type"].(string) == "REPEATING" {
			var sl_repeating_trigger datatypes.SoftLayer_Scale_Policy_Trigger_Repeating
			sl_repeating_trigger.TypeId = 2
			sl_repeating_trigger.Schedule = trigger["schedule"].(string)
			sl_repeating_triggers = append(sl_repeating_triggers, sl_repeating_trigger)
		}
	}
	return sl_repeating_triggers
}

func prepareResourceUseTriggers(raw_triggers []interface{}) []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse {
	sl_resourceuse_triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse, 0)
	for _, raw_trigger := range raw_triggers {
		trigger := raw_trigger.(map[string]interface{})

		if trigger["type"].(string) == "RESOURCE_USE" {
			var sl_resourceuse_trigger datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse
			sl_resourceuse_trigger.TypeId = 1
			sl_resourceuse_trigger.Watches = prepareWatches(trigger["watches"].([]interface{}))
			sl_resourceuse_triggers = append(sl_resourceuse_triggers, sl_resourceuse_trigger)
		}
	}
	return sl_resourceuse_triggers
}

func prepareWatches(raw_watches []interface{}) []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch {
	sl_watches := make([]datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch, 0)
	for _, raw_watch := range raw_watches {
		var sl_watch datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch
		watch := raw_watch.(map[string]interface{})

		sl_watch.Metric = watch["metric"].(string)
		sl_watch.Operator = watch["operator"].(string)
		sl_watch.Period = watch["period"].(int)
		sl_watch.Value = watch["value"].(string)
		sl_watch.Algorithm = "EWMA"

		sl_watches = append(sl_watches, sl_watch)
	}
	return sl_watches
}