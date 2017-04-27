package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

func podTemplateSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": namespacedMetadataSchema("pod template spec", true),
		"spec": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Specification of the desired behavior of the pod. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status",
			Elem: &schema.Resource{
				Schema: podSpecFileds(),
			},
		},
	}
}

func podSpecFileds() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"active_deadline_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer.",
		},
		"dns_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ClusterFirst",
			Description: "Set DNS policy for containers within the pod. One of 'ClusterFirst' or 'Default'. Defaults to 'ClusterFirst'.",
		},
		"host_ipc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use the host's ipc namespace. Optional: Default to false.",
		},
		"host_network": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Host networking requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified.",
		},
		"hostname": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.",
		},
		"node_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements.",
		},
		"node_selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: http://kubernetes.io/docs/user-guide/node-selection.",
		},
		"restart_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Always",
			Description: "Restart policy for all containers within the pod. One of Always, OnFailure, Never. More info: http://kubernetes.io/docs/user-guide/pod-states#restartpolicy.",
		},
		"service_account_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: http://releases.k8s.io/HEAD/docs/design/service_accounts.md.",
		},
		"subdomain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `If specified, the fully qualified Pod hostname will be "...svc.". If not specified, the pod will not have a domainname at all..`,
		},
		"termination_grace_period_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     30,
			Description: "Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process.",
		},
		"image_pull_secrets": {
			Type:        schema.TypeList,
			Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: http://kubernetes.io/docs/user-guide/images#specifying-imagepullsecrets-on-a-pod",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
						Optional:    true,
					},
				},
			},
		},
		"volumes": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of volumes that can be mounted by containers belonging to the pod. More info: http://kubernetes.io/docs/user-guide/volumes",
			Elem:        volumeSchema(),
		},
		"containers": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers",
			Elem: &schema.Resource{
				Schema: containerFields(),
			},
		},
	}
}

func volumeSchema() *schema.Resource {
	v := map[string]*schema.Schema{}

	v["persistent_volume_claim"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "The specification of a persistent volume.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"claim_name": {
					Type:        schema.TypeString,
					Description: "ClaimName is the name of a PersistentVolumeClaim in the same ",
					Optional:    true,
				},
				"read_only": {
					Type:        schema.TypeBool,
					Description: "Will force the ReadOnly setting in VolumeMounts.",
					Optional:    true,
					Default:     false,
				},
			},
		},
	}

	v["secret"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Secret represents a secret that should populate this volume. More info: http://kubernetes.io/docs/user-guide/volumes#secrets",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"secret_name": {
					Type:        schema.TypeString,
					Description: "Name of the secret in the pod's namespace to use. More info: http://kubernetes.io/docs/user-guide/volumes#secrets",
					Optional:    true,
				},
			},
		},
	}
	v["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Volume's name. Must be a DNS_LABEL and unique within the pod. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: v,
	}
}
