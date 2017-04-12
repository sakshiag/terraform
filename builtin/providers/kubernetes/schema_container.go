package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

func resourcesField() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"limits": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Describes the maximum amount of compute resources allowed. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
		},
		"requests": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: "Describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
		},
	}
}

func volumeMountFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mount_path": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Path within the container at which the volume should be mounted. Must not contain ':'.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "This must match the Name of a Volume.",
		},
		"read_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.",
		},
		"sub_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: `Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).`,
		},
	}
}

func containerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"command": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Description: "Entrypoint array. Not executed within a shell. The docker image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands",
		},
		"args": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Description: "Arguments to the entrypoint. The docker image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands",
		},
		"image": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Docker image name. More info: http://kubernetes.io/docs/user-guide/images",
		},
		"image_pull_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/images#updating-images",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). Cannot be updated.",
		},
		"working_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated.",
		},
		"resources": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Compute Resources required by this container. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources",
			Elem: &schema.Resource{
				Schema: resourcesField(),
			},
		},
		"volume_mounts": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Pod volumes to mount into the container's filesystem. Cannot be updated.",
			Elem: &schema.Resource{
				Schema: volumeMountFields(),
			},
		},
		"ports": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: `List of ports to expose from the container. Exposing a port here gives the system additional information about the network connections a container uses, but is primarily informational. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Cannot be updated.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container_port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.",
					},
					"host_ip": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "What host IP to bind the external port to.",
					},
					"host_port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services",
					},
					"protocol": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `Protocol for port. Must be UDP or TCP. Defaults to "TCP".`,
						Default:     "TCP",
					},
				},
			},
		},
	}
}
