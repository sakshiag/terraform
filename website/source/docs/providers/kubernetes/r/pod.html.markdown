---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_pod"
sidebar_current: "docs-kubernetes-pod"
description: |-
  Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts.
---

# kubernetes_pod

Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts.

More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/

## Example Usage

```

resource "kubernetes_secret" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			TestLabelOne = "one"
			TestLabelThree = "three"
		}
		name = "secret"
	}
	data {
		one = "first"
		two = "second"
		nine = "ninth"
	}
}

resource "kubernetes_config_map" "test" {
	metadata {
		name = "some-config-map"
	}
	data {
		one = "ONE"
	}
}

resource "kubernetes_pod" "test" {
	metadata {
		labels {
			app  = "pod_label"
		}
		name = "test-pod"
	}
	spec {
		containers {
			image = "nginx:1.7.9"
			name = "containername"
			env = [{
				name = "EXPORTED_VARIBALE_FROM_SECRET"
				value_from {
					secret_key_ref{
						name = "${kubernetes_secret.test.metadata.0.name}"
						key = "one"
					}
				}
			},
			{
				name = "EXPORTED_VARIBALE_FROM_CONFIG_MAP"
				value_from {
					config_map_key_ref{
						name = "${kubernetes_config_map.test.metadata.0.name}"
						key = "one"
					}
				}
			}]
		}
		volumes =  [{
        name = "mycloudant",
          secret =  {
            secret_name =  "${kubernetes_secret.test.metadata.0.name}"
          }
         }]
	}
}

```

## Argument Reference

The following arguments are supported:

* `metadata` - (Required) Standard replication controller's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `spec` - (Required) Spec of the  replication controller owned by the cluster. See below.

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

* `active_deadline_seconds` - (Optional) Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer.
* `dns_policy` - (Optional) Set DNS policy for containers within the pod. One of 'ClusterFirst' or 'Default'. Defaults to "ClusterFirst".Cannot be updated.
* `host_ipc` - (Optional) Use the host's ipc namespace. Optional: Default to false.Cannot be updated.
* `host_network` - (Optional) Host networking requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false.Cannot be updated.
* `host_pid` - (Optional) Use the host's pid namespace. Optional: Default to false.Cannot be updated.
* `hostname` - (Optional) Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.Cannot be updated.
* `node_name` - (Optional) NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements.Cannot be updated.
* `node_selector` - (Optional) Cannot be updated.NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: http://kubernetes.io/docs/user-guide/node-selection
* `restart_policy` - (Optional) Cannot be updated.Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always. More info: http://kubernetes.io/docs/user-guide/pod-states#restartpolicy
* `service_account_name` - (Optional) Cannot be updated.ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: http://releases.k8s.io/HEAD/docs/design/service_accounts.md
* `subdomain` - (Optional) Cannot be updated.If specified, the fully qualified Pod hostname will be "...svc.". If not specified, the pod will not have a domainname at all.
* `termination_grace_period_seconds` - (Optional)Cannot be updated. Optional duration in seconds the pod needs to terminate gracefully. 
* `containers` - (Required) List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers
* `image_pull_secrets` - (Optional) Cannot be updated.ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: http://kubernetes.io/docs/user-guide/images#specifying-imagepullsecrets-on-a-pod
* `security_context` - (Optional)Cannot be updated. SecurityContext holds pod-level security attributes and common container settings. Optional: Defaults to empty. See type description for default values of each field.
* `volumes` - (Optional) Cannot be updated.List of volumes that can be mounted by containers belonging to the pod. More info: http://kubernetes.io/docs/user-guide/volumes

### `containers`

#### Arguments

* `args` - (Optional) Arguments to the entrypoint. The docker image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment.Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands
* `command` - (Optional) Entrypoint array. Not executed within a shell. The docker image's ENTRYPOINT is used if this is not provided. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/containers#containers-and-commands
* `env` - (Optional) List of environment variables to set in the container. Cannot be updated.
* `image` - (Required) Docker image name. More info: http://kubernetes.io/docs/user-guide/images
* `image_pull_policy` - (Optional) Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/images#updating-images
* `name` - (Required) Name of the container specified Name of the container specified as a (DNS LABEL). Each container in a pod must have a unique name (DNS LABEL). Cannot be updated.
* `ports` - (Optional) List of ports to expose from the container. Exposing a port here gives the system additional information about the network connections a container uses, but is primarily informational. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Cannot be updated.
* `resources` - (Optional) Compute Resources required by this container. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources
* `volume_mounts` - (Optional) Pod volumes to mount into the container's filesystem. Cannot be updated.
* `workingDir` - (Optional) Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated.


### `env`

#### Arguments

* `name` - (Optional) Name of the environment variable. Must be a C_IDENTIFIER.
* `value` - (Optional) Variable references are expanded using the previous defined environment variables in the container and any service environment variables.
* `valueFrom` - (Optional) Source for the environment variable's value. Cannot be used if value is not empty.

### `valueFrom`

#### Arguments

* `configMapKeyRef` - (Optional) Selects a key of a ConfigMap.
* `fieldRef` - (Optional) Selects a field of the pod: supports metadata.name, metadata.namespace, metadata.labels, metadata.annotations, spec.nodeName, spec.serviceAccountName, status.podIP.
* `resourceFieldRef` - (Optional) Selects a resource of the container: only resources limits and requests (limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.
* `secretKeyRef` - (Optional) Selects a key of a secret in the pod's namespace.

### `configMapKeyRef`

#### Arguments

* `key` - (Optional) The key to select.
* `name` - (Optional) Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names


### `fieldRef`

#### Arguments

* `apiVersion` - (Optional) Version of the schema the FieldPath is written in terms of, defaults to "v1".
* `fieldPath` - (Optional) Path of the field to select in the specified API version.


### `resourceFieldRef`

#### Arguments

* `containerName` - (Optional) Container name: required for volumes, optional for env vars
* `divisor` - (Optional) Specifies the output format of the exposed resources, defaults to "1"
* `resource` - (Required) Required: resource to select


### `secretKeyRef`

#### Arguments

* `key` - (Optional) The key of the secret to select from. Must be a valid secret key.
* `name` - (Optional) Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names

### `ports`

#### Arguments

* `container_port ` - (Required) Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
* `host_ip` - (Optional) What host IP to bind the external port to.
* `host_port` - (Optional) Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.
* `name` - (Optional) If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services
* `protocol` - (Optional) Protocol for port. Must be UDP or TCP. Defaults to "TCP". 

### `resources`

#### Arguments

* `limits` - (Optional) Describes the maximum amount of compute resources allowed. More info: http://kubernetes.io/docs/user-guide/compute-resources/
* `requests` - (Required) Describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: http://kubernetes.io/docs/user-guide/compute-resources/

### `volume_mounts`

#### Arguments

* `mount_path` - (Required) Path within the container at which the volume should be mounted. Must not contain ':'.
* `name` - (Required) This must match the Name of a Volume.
* `read_only` - (Optional) Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.
* `sub_path` - (Optional) Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).



Pods can be imported using its name, e.g.

```
$ terraform import kubernetes_pod.test test-pod
```
