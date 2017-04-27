package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

func servicePortFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of this port within the service. This must be a DNS_LABEL. All ports within a ServiceSpec must have unique names. This maps to the 'Name' field in EndpointPort objects. Optional if only one ServicePort is defined on this service.",
		},
		"node_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The port on each node on which this service is exposed when type=NodePort or LoadBalancer. Usually assigned by the system. If specified, it will be allocated to the service if unused or else creation of the service will fail. Default is to auto-allocate a port if the ServiceType of this Service requires one.",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The port that will be exposed by this service.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IP protocol for this port. Supports 'TCP' and 'UDP'. Default is TCP.",
		},
		"target_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number or name of the port to access on the pods targeted by the service. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME. If this is a string, it will be looked up as a named port in the target Pod's container ports. If this is not specified, the value of the 'port' field is used (an identity map). This field is ignored for services with clusterIP=None, and should be omitted or set equal to the 'port' field.",
		},
	}
}

func serviceSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_iP": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Description: "ClusterIP is the IP address of the service and is usually assigned randomly by the master. If an address is specified manually and is not in use by others, it will be allocated to the service; otherwise, creation of the service will fail. This field can not be changed through updates.",
			Optional:    true,
		},
		"external_ips": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "ExternalIPs is a list of IP addresses for which nodes in the cluster will also accept traffic for this service. These IPs are not managed by Kubernetes. The user is responsible for ensuring that traffic arrives at a node with this IP.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
		},
		"external_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ExternalName is the external reference that kubedns or equivalent will return as a CNAME record for this service. No proxying will be involved. Must be a valid DNS name and requires Type to be ExternalName.",
		},
		"load_balancer_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only applies to Service Type: LoadBalancer LoadBalancer will get created with the IP specified in this field. This feature depends on whether the underlying cloud-provider supports specifying the loadBalancerIP when a load balancer is created. This field will be ignored if the cloud-provider does not support the feature.",
		},
		"load_balancer_source_ranges": {
			Type:        schema.TypeSet,
			Optional:    true,
			ForceNew:    true,
			Description: "If specified and supported by the platform, this will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs. This field will be ignored if the cloud-provider does not support the feature.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
		},
		"session_affinity": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Supports 'ClientIP' and 'None'. Used to maintain session affinity. Enable client IP based session affinity. Must be ClientIP or None. Defaults to None.",
		},
		"selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Route service traffic to pods with label keys and values matching this selector. If empty or not present, the service is assumed to have an external process managing its endpoints, which Kubernetes will not modify. Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if type is ExternalName.",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type determines how the Service is exposed. Defaults to ClusterIP. Valid options are ExternalName, ClusterIP, NodePort, and LoadBalancer.",
		},
		"ports": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: servicePortFields(),
			},
		},
	}
}
