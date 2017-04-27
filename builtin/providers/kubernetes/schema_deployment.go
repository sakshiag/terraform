package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

func deploymentSpecFileds() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"min_ready_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)",
		},
		"replicas": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Number of desired pods.",
		},
		"selector": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Label selector for pods. Existing ReplicaSets whose pods are selected by this will be the ones affected by this deployment.",
			Elem: &schema.Resource{
				Schema: selectorFields(),
			},
		},

		"strategy": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The deployment strategy to use to replace existing pods with new ones.",
			Elem: &schema.Resource{
				Schema: strategySchema(),
			},
		},
		"revision_history_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The number of old ReplicaSets to retain to allow rollback.",
		},
		"pause": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates that the deployment is paused and will not be processed by the deployment controller.",
		},

		"progress_deadline_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The maximum time in seconds for a deployment to make progress before it is considered to be failed.",
		},
		"rollback_to": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The config this deployment is rolling back to. Will be cleared after rollback is done.",
			Elem: &schema.Resource{
				Schema: rollbackToSchema(),
			},
		},
		"template": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Template describes the pods that will be created.",
			Elem: &schema.Resource{
				Schema: podTemplateSpecSchema(),
			},
		},
	}
}

func strategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "RollingUpdate",
			Description: "Type of deployment. Can be Recreate or RollingUpdate. Default is RollingUpdate.",
		},
		"rolling_update": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Update this to follow our convention for oneOf, whatever we decide it to be.",
			Elem: &schema.Resource{
				Schema: rollingUpdateSchema(),
			},
		},
	}
}

func rollingUpdateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"max_unavailable": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The maximum number of pods that can be unavailable during the update.",
		},
		"max_surge": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The maximum number of pods that can be scheduled above the desired number of pods.",
		},
	}
}

func rollbackToSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"revision": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The revision to rollback to. If set to 0, rollbck to the last revision.",
		},
	}
}
