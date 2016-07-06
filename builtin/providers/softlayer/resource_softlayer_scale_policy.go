package softlayer

import (
	"fmt"
	"log"
	"strconv"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						// Conditionally-required fields, based on value of "type"
						"watches": &schema.Schema{
							Type:     schema.TypeSet,
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
							Set : resourceSoftLayerScalePolicyHandlerHash,
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
				Set : resourceSoftLayerScalePolicyTriggerHash,
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

	if _, ok := d.GetOk("triggers"); ok {
		opts.OneTimeTriggers = prepareOneTimeTriggers(d)
		opts.RepeatingTriggers = prepareRepeatingTriggers(d)
		opts.ResourceUseTriggers = prepareResourceUseTriggers(d)
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
	triggers := make([]map[string]interface{}, 0)
	triggers = append(triggers, readOneTimeTriggers(scalePolicy.OneTimeTriggers)...)
	triggers = append(triggers, readRepeatingTriggers(scalePolicy.RepeatingTriggers)...)
	triggers = append(triggers, readResourceUseTriggers(scalePolicy.ResourceUseTriggers)...)

	d.Set("triggers", triggers)

	return nil
}

func resourceSoftLayerScalePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).scalePolicyService
	triggerClient := meta.(*Client).scalePolicyTriggerService

	scalePolicyId, _ := strconv.Atoi(d.Id())

	scalePolicy, err := client.GetObject(scalePolicyId)
	if err != nil {
		return fmt.Errorf("Error retrieving scalePolicy: %s", err)
	}

	var template datatypes.SoftLayer_Scale_Policy

	template.Id, _ = strconv.Atoi(d.Id())

	if d.HasChange("name") {
		template.Name = d.Get("name").(string)
	}

	if d.HasChange("scale_type") || d.HasChange("scale_amount") {
		template.ScaleActions = []datatypes.SoftLayer_Scale_Policy_Action{{
			Id : scalePolicy.ScaleActions[0].Id,
			TypeId: 1,
		}}
	}
	if d.HasChange("scale_type") {
		template.ScaleActions[0].ScaleType = d.Get("scale_type").(string)
	}

	if d.HasChange("scale_amount") {
		template.ScaleActions[0].Amount = d.Get("scale_amount").(int)
	}

	if d.HasChange("cooldown") {
		template.Cooldown = d.Get("cooldown").(int)
	}

	for _, triggerList := range scalePolicy.Triggers {
		triggerClient.DeleteObject(triggerList.Id)
	}

	time.Sleep(60)
	if _, ok := d.GetOk("triggers"); ok {
		template.OneTimeTriggers = prepareOneTimeTriggers(d)
		template.RepeatingTriggers = prepareRepeatingTriggers(d)
		template.ResourceUseTriggers = prepareResourceUseTriggers(d)
	}

	_, err = client.EditObject(scalePolicyId, template)

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

func prepareOneTimeTriggers(d *schema.ResourceData) []datatypes.SoftLayer_Scale_Policy_Trigger_OneTime {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_OneTime, 0)
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "ONE_TIME" {
			var oneTimeTrigger datatypes.SoftLayer_Scale_Policy_Trigger_OneTime
			oneTimeTrigger.TypeId = datatypes.SOFTLAYER_SCALE_POLICY_TRIGGER_TYPE_ID_ONE_TIME
			timeStampString := trigger["date"].(string)
			timeStamp, _ := time.Parse(time.RFC3339Nano, timeStampString)
			oneTimeTrigger.Date = &timeStamp
			triggers = append(triggers, oneTimeTrigger)
		}
	}
	return triggers
}

func prepareRepeatingTriggers(d *schema.ResourceData) []datatypes.SoftLayer_Scale_Policy_Trigger_Repeating {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_Repeating, 0)
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "REPEATING" {
			var repeatingTrigger datatypes.SoftLayer_Scale_Policy_Trigger_Repeating
			repeatingTrigger.TypeId = datatypes.SOFTLAYER_SCALE_POLICY_TRIGGER_TYPE_ID_REPEATING
			repeatingTrigger.Schedule = trigger["schedule"].(string)
			triggers = append(triggers, repeatingTrigger)
		}
	}
	return triggers
}

func prepareResourceUseTriggers(d *schema.ResourceData) []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse, 0)
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "RESOURCE_USE" {
			var resourceUseTrigger datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse
			resourceUseTrigger.TypeId = datatypes.SOFTLAYER_SCALE_POLICY_TRIGGER_TYPE_ID_RESOURCE_USE
			resourceUseTrigger.Watches = prepareWatches(trigger["watches"].(*schema.Set))
			triggers = append(triggers, resourceUseTrigger)
		}
	}
	return triggers
}

func prepareWatches(d *schema.Set) []datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch {
	watchLists := d.List()
	watches := make([]datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch, 0)
	for _, watcheList := range watchLists {
		var watch datatypes.SoftLayer_Scale_Policy_Trigger_ResourceUse_Watch
		watchMap := watcheList.(map[string]interface{})

		watch.Metric = watchMap["metric"].(string)
		watch.Operator = watchMap["operator"].(string)
		watch.Period = watchMap["period"].(int)
		watch.Value = watchMap["value"].(string)
		watch.Algorithm = "EWMA"

		watches = append(watches, watch)
	}
	return watches
}

func readOneTimeTriggers(list []datatypes.SoftLayer_Scale_Policy_Trigger_OneTime) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = trigger.Id
		t["type"] = "ONE_TIME"
//		t["date"] = trigger.Date.Format(time.RFC3339Nano)
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

func resourceSoftLayerScalePolicyTriggerHash(v interface{}) int {
	var buf bytes.Buffer
	trigger := v.(map[string]interface{})
	if  trigger["type"].(string) == "ONE_TIME" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", trigger["date"].(string)))
	}
	if  trigger["type"].(string) == "REPEATING" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", trigger["schedule"].(string)))
	}
	if  trigger["type"].(string) == "RESOURCE_USE" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		for _, watchList := range trigger["watches"].(*schema.Set).List() {
			watch := watchList.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", watch["metric"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", watch["operator"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", watch["value"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", watch["period"].(int)))
		}
	}
	return hashcode.String(buf.String())
}

func resourceSoftLayerScalePolicyHandlerHash(v interface{}) int {
	var buf bytes.Buffer
	watch := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", watch["metric"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", watch["operator"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", watch["value"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", watch["period"].(int)))
	return hashcode.String(buf.String())
}