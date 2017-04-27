---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_service"
sidebar_current: "docs-kubernetes-resource-service"
description: |-
  The resource provides mechanisms to provide a abstraction which defines a logical set of Pods and a policy by which to access them.
---

# kubernetes_service

The resource provides mechanisms of grouping of pods that are running on the cluster.Services can efficiently power a microservice architecture.
Services provide important features that are standardized across the cluster: load-balancing, service discovery between applications, and features to support zero-downtime application deployments.
The resource will by default create a service in the specified (or default) namespace.

~> Read more about service : https://kubernetes.io/docs/concepts/services-networking/service/

## Example Usage

```
resource "kubernetes_service" "example" {
	metadata {
		labels {
			app  = "service_label"
		}
		name = "service-name"
	}
	spec {
		type = "NodePort"
		ports {
			name = "some-name"
			node_port = 30001
			port = 80
			protocol = "TCP"
			target_port = 8989
		}
	}
}
```

## Argument Reference

The following arguments are supported:

* `metadata` - (Required) Standard service metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `spec` - (Required) The Spec is the desired state of the service. 

## Nested Blocks

### `metadata`

#### Arguments

* `annotations` - (Optional) An unstructured key value map stored with the service that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
* `generate_name` - (Optional) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#idempotency
* `labels` - (Optional) Map of string keys and values that can be used to organize and categorize (scope and select) the service. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
* `name` - (Optional) Name of the service, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
* `namespace` - (Optional) Namespace defines the space within which name of the service must be unique.

### `spec`

#### Arguments

* `cluster_ip` - (Optional) ClusterIP is the IP address of the service and is usually assigned randomly by the master. If an address is specified manually and is not in use by others, it will be allocated to the service; otherwise, creation of the service will fail. This field can not be changed through updates.
* `external_ips` - (Optional) ExternalIPs is a list of IP addresses for which nodes in the cluster will also accept traffic for this service. These IPs are not managed by Kubernetes. The user is responsible for ensuring that traffic arrives at a node with this IP.
* `external_name` - (Optional) ExternalName is the external reference that kubedns or equivalent will return as a CNAME record for this service.
* `load_balancer_ip` - (Optional) Only applies to Service Type: LoadBalancer LoadBalancer will get created with the IP specified in this field.
* `session_affinity` - (Optional) Supports 'ClientIP' and 'None'. Used to maintain session affinity. Enable client IP based session affinity. Must be ClientIP or None. Defaults to None.
* `type` - (Optional) Type determines how the Service is exposed. Defaults to ClusterIP. Valid options are ExternalName, ClusterIP, NodePort, and LoadBalancer.
* `ports` - (Optional) The list of ports that are exposed by this service.


### `ports`

#### Arguments

* `name` - (Optional)The name of this port within the service. This must be a DNS_LABEL. All ports within a ServiceSpec must have unique names.
* `node_port` - (Optional)The port on each node on which this service is exposed when type=NodePort or LoadBalancer. Usually assigned by the system. If specified, it will be allocated to the service if unused or else creation of the service will fail.
* `port` - (Optional)The port that will be exposed by this service.
* `protocol` - (Optional)The IP protocol for this port. Supports 'TCP' and 'UDP'. Default is TCP.
* `target_port` - (Optional)Number or name of the port to access on the pods targeted by the service. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME. If this is a string, it will be looked up as a named port in the target Pod's container ports.


## Import

service can be imported using its name, e.g.

```
$ terraform import kubernetes_service.example service-name
```


