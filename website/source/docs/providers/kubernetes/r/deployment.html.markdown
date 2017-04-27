---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_deployment"
sidebar_current: "docs-kubernetes-resource-deployment"
description: |-
  A Deployment ensures that a specified number of pod “replicas” are running at any one time in the cluster provisioned by an administrator.
---

# kubernetes_deployment

Deployment enables declarative updates for Pods and ReplicaSets.

More info: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

## Example Usage

```
resource "kubernetes_deployment" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			app  = "test-for-deployment"
		}
		name = "test-deployment"
	}
	spec {
		replicas =  3
		template {
			metadata {
					labels {
						app  = "tempalteapp"
					}
			}
			spec {
				containers = [{
					image = "nginx:1.7.9"
					name = "test1"
				}]
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

* `replicas` - (Optional) Number of desired pods. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.
* `min_ready_seconds` - (Optional) Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)
* `template` - (Optional) Template is the object that describes the pod that will be created if insufficient replicas are detected. This takes precedence over a TemplateRef. More info: http://kubernetes.io/docs/user-guide/replication-controller#pod-template
* `paused` - (Optional) Indicates that the deployment is paused and will not be processed by the deployment controller
* `progress_deadline_seconds` - (Optional) The maximum time in seconds for a deployment to make progress before it is considered to be failed. The deployment controller will continue to process failed deployments and a condition with a ProgressDeadlineExceeded reason will be surfaced in the deployment status. Once autoRollback is implemented, the deployment controller will automatically rollback failed deployments. Note that progress will not be estimated during the time a deployment is paused. This is not set by default.
* `revision_history_limit` - (Optional) The number of old ReplicaSets to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified.
* `rollback_to` - (Optional) The config this deployment is rolling back to. Will be cleared after rollback is done.
* `selector` - (Optional) Label selector for pods. Existing ReplicaSets whose pods are selected by this will be the ones affected by this deployment.
* `strategy` - (Optional) The deployment strategy to use to replace existing pods with new ones.

### `rollback_to`

#### Arguments

* `revision` - (Optional) The revision to rollback to. If set to 0, rollbck to the last revision.

### `selector`

#### Arguments

* `match_expressions` - matchExpressions is a list of label selector requirements. The requirements are ANDed.
* `match_labels ` - matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.

### `match_expressions`

#### Arguments

* `key` - key is the label key that the selector applies to.
* `operator` - operator represents a key's relationship to a set of values. Valid operators ard In, NotIn, Exists and DoesNotExist.
* `values` -  values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.

### `strategy`

#### Arguments

* `rolling_update` - Rolling update config params. Present only if DeploymentStrategyType = RollingUpdate.

### `rolling_update`

#### Arguments

* `max_surge` - The maximum number of pods that can be scheduled above the desired number of pods.
* `max_unavailable` - The maximum number of pods that can be unavailable during the update.

### `delete_options`

#### Arguments

* `orphan_dependents` - (Optional) Should the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added to/removed from the object's finalizers list.
 


Deployment can be imported using its name, e.g.

```
$ terraform import kubernetes_deployment.test test-deployment
```
