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
	client := meta.(*Client).scalePolicyService
	scalePolicyId, _ := strconv.Atoi(d.Id())

	scalePolicy, err := client.GetObject(scalePolicyId)
	if err != nil {
		return fmt.Errorf("Error retrieving Scale Policy: %s", err)
	}
	d.Set("id", scalePolicy.Id)
	d.Set("name", scalePolicy.Name)
	d.Set("cooldown", scalePolicy.Cooldown)
	d.Set("scaleGroupId", scalePolicy.ScaleGroupId)
	d.Set("oneTimeTriggers", readOneTimeTriggers(scalePolicy.OneTimeTriggers))
	d.Set("repeatingTriggers", readRepeatingTriggers(scalePolicy.RepeatingTriggers))
	d.Set("resourceUseTriggers", readResourceUseTriggers(scalePolicy.ResourceUseTriggers))

	return nil
}

func resourceSoftLayerScalePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scalePolicyService

	scalePolicyId := d.Get("id").(int)
	result, err := client.GetObject(scalePolicyId)
	if err != nil {
		return fmt.Errorf("Error retrieving scale policy: %s", err)
	}

	if d.HasChange("name") {
		result.Name = d.Get("name").(string)
	}

	if d.HasChange("scale_type") {
		result.ScaleActions[0].ScaleType = d.Get("scale_type").(string)
	}

	if d.HasChange("scale_amount") {
		result.ScaleActions[0].Amount = d.Get("scale_amount").(int)
	}

	if d.HasChange("cooldown") {
		result.Cooldown = d.Get("cooldown").(int)
	}

	triggers := d.Get("triggers").([]interface{})

	count := 0
	countOneTime := 0
	countRepeating := 0
	countResourceUse := 0

	for _, _ = range triggers {
		if d.Get("triggers."+strconv.Itoa(count)+".type") == "ONE_TIME" {
			if d.HasChange("triggers."+strconv.Itoa(count)+".data") {
				timeStamp, _ := time.Parse(time.RFC3339Nano, d.Get("triggers."+strconv.Itoa(count)+".date").(string))
				result.OneTimeTriggers[countOneTime].Date = &timeStamp
			}
			countOneTime++
		}
		if d.Get("triggers."+strconv.Itoa(count)+".type") == "REPEATING" {
			if d.HasChange("triggers."+strconv.Itoa(count)+".schedule") {
				result.RepeatingTriggers[countOneTime].Schedule = d.Get("triggers."+strconv.Itoa(count)+".schedule").(string)
			}
			countRepeating++
		}
		if d.Get("triggers."+strconv.Itoa(count)+".type") == "RESOURCE_USE" {
			watches := d.Get("triggers."+strconv.Itoa(count)+".watches").([]interface{})
			for _, _ = range watches {
				if d.HasChange("triggers." + strconv.Itoa(count) + ".watches." + strconv.Itoa(countResourceUse) + ".period") {
					result.ResourceUseTriggers[count].Watches[countResourceUse].Period = d.Get("triggers." + strconv.Itoa(count) + ".watches." + strconv.Itoa(countResourceUse) + ".period").(int)
//					if d.HasChange("triggers." + strconv.Itoa(0) + ".watches." + strconv.Itoa(0) + ".period") {
//						result.ResourceUseTriggers[0].Watches[0].Period = d.Get("triggers." + strconv.Itoa(0) + ".watches." + strconv.Itoa(0) + ".period").(int)

				}
			}
			countResourceUse++
		}
		count++
	}

	_, err = client.EditObject(scalePolicyId, result)

	if err != nil {
		return fmt.Errorf("Error updating scalie policy: %s", err)
	}

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

func readOneTimeTriggers(list []datatypes.SoftLayer_Scale_Policy_Trigger_OneTime) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = trigger.Id
		t["type"] = "ONE_TIME"
		t["date"] = trigger.Date.String()
		triggers = append(triggers, t)
	}
	return triggers
}

func readRepeatingTriggers(list []datatypes.SoftLayer_Scale_Policy_Trigger_Repeating) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = trigger.Id
		t["type"] = "REPEATING"
		t["schedule"] = trigger.Schedule
		triggers = append(triggers, t)
	}
	return triggers
}

func readResourceUseTriggers(list []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = trigger.Id
		t["type"] = "RESOURCE_USE"
		t["watches"] = readResourceUseWatches(trigger.Watches)
		triggers = append(triggers, t)
	}
	return triggers
}

func readResourceUseWatches(list []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch) []map[string]interface{} {
	watches := make([]map[string]interface{}, 0, len(list))
	for _, watch := range list {
		w := make(map[string]interface{})
		w["id"] = watch.Id
		w["metric"] = watch.Metric
		w["operator"] = watch.Operator
		w["period"] = watch.Period
		w["value"] = watch.Value
		watches = append(watches, w)
	}
	return watches
}