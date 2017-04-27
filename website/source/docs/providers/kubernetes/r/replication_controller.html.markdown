---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_replication_controller"
sidebar_current: "docs-kubernetes-resource-replication-controller"
description: |-
  A ReplicationController ensures that a specified number of pod “replicas” are running at any one time in the cluster provisioned by an administrator.
---

# kubernetes_replication_controller

The resource ensures that a specified number of pod “replicas” are running at any one time in the cluster provisioned by an administrator. It is a resource in the cluster just like a node is a cluster resource. 

More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/

## Example Usage

```
resource "kubernetes_replication_controller" "replication-controller" {
	metadata {
	 labels {
		app  = "RC_for_UAT_test"
	 }
	 name = "replication-controller"
	}
	spec {
	 min_ready_seconds = 60
	 replicas = 2
	 template {
	  metadata {
		labels {
		  app  = "replicationUATAPP"
		}
	  }
	  spec {
		containers {
		  image = "nginx:1.7.9"
		  name = "uattest"
			resources{
				limits {
					cpu = "500m"
					memory = "128Mi"
				}
				requests {
					   memory = "64Mi"
        		cpu =  "250m"
				}
			}
		}
	  }
	}
   }
   delete_options {
   	orphan_dependents = false
   }
}


```

## Argument Reference

The following arguments are supported:

* `metadata` - (Required) Standard replication controller's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `spec` - (Required) Spec of the  replication controller owned by the cluster. See below.
* `delete_options` - (Optional) DeleteOptions may be provided when deleting an API object

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

* `replicas` - (Optional) Replicas is the number of desired replicas. This is a pointer to distinguish between explicit zero and unspecified. More info: http://kubernetes.io/docs/user-guide/replication-controller#what-is-a-replication-controller
* `min_ready_seconds` - (Optional) Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)
* `selector` - (Optional) Selector is a label query over pods that should match the Replicas count. If Selector is empty, it is defaulted to the labels present on the Pod template. Label keys and values that must match in order to be controlled by this replication controller, if empty defaulted to labels on Pod template. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors
* `template` - (Optional) Template is the object that describes the pod that will be created if insufficient replicas are detected. This takes precedence over a TemplateRef. More info: http://kubernetes.io/docs/user-guide/replication-controller#pod-template


### `selector`

#### Arguments

* `match_expressions` - matchExpressions is a list of label selector requirements. The requirements are ANDed.
* `match_labels ` - matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.

### `match_expressions`

#### Arguments

* `key` - key is the label key that the selector applies to.
* `operator` - operator represents a key's relationship to a set of values. Valid operators ard In, NotIn, Exists and DoesNotExist.
* `values` -  values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.

### `delete_options`

#### Arguments

* `orphan_dependents` - (Optional) Should the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added to/removed from the object's finalizers list.

Replication Controller can be imported using its name, e.g.

```
$ terraform import kubernetes_replication_controller.replication-controller replication-controller
```
