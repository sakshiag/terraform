package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

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
			Description: `Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).`,
		},
	}
}

func containerFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"command": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Entrypoint array. Not executed within a shell. The docker image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands",
		},
		"args": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Arguments to the entrypoint. The docker image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands",
		},

		"image": {
			Type:        schema.TypeString,
			Optional:    true,
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
		"env": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "List of environment variables to set in the container. Cannot be updated.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the environment variable. Must be a C_IDENTIFIER",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `Variable references $(VAR_NAME) are expanded using the previous defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "".`,
					},
					"value_from": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Source for the environment variable's value",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"config_map_key_ref": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Selects a key of a ConfigMap.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "The key to select.",
											},
											"name": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
											},
										},
									},
								},
								"field_ref": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Selects a field of the pod: supports metadata.name, metadata.namespace, metadata.labels, metadata.annotations, spec.nodeName, spec.serviceAccountName, status.podIP..",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"api_version": {
												Type:        schema.TypeString,
												Optional:    true,
												Default:     "v1",
												Description: `Version of the schema the FieldPath is written in terms of, defaults to "v1".`,
											},
											"field_path": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Path of the field to select in the specified API version",
											},
										},
									},
								},
								"resource_field_ref": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Selects a field of the pod: supports metadata.name, metadata.namespace, metadata.labels, metadata.annotations, spec.nodeName, spec.serviceAccountName, status.podIP..",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"container_name": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Container name: required for volumes, optional for env vars",
											},
											"resource": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Resource to select",
											},
										},
									},
								},
								"secret_key_ref": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Selects a field of the pod: supports metadata.name, metadata.namespace, metadata.labels, metadata.annotations, spec.nodeName, spec.serviceAccountName, status.podIP..",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "The key of the secret to select from. Must be a valid secret key.",
											},
											"name": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
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
		"resources": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			ForceNew:    true,
			Description: "Compute Resources required by this container. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources",
			Elem: &schema.Resource{
				Schema: resourcesField(),
			},
		},
		"volume_mounts": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
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
						Required:    true,
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
