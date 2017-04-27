package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

// These are used by both Container and PersistentVolumeSpec
// https://kubernetes.io/docs/api-reference/v1.5/#resourcerequirements-v1
// Hence putting in the seprate file as helper

func resourcesField() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"limits": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Describes the maximum amount of compute resources allowed. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
		},
		"requests": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
		},
	}
}

func selectorFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match_labels": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: `matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.`,
		},
		"match_expressions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "matchExpressions is a list of label selector requirements. The requirements are ANDed.",
			Elem: &schema.Resource{
				Schema: labelSelectorRequirementFields(),
			},
		},
	}
}

func labelSelectorRequirementFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "key is the label key that the selector applies to.",
		},
		"operator": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "operator represents a key's relationship to a set of values.",
		},
		"values": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
	}
}
