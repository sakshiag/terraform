---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_ingress"
sidebar_current: "docs-kubernetes-resource-ingress"
description: |-
  The resource provides mechanisms to provide a collection of rules that allow inbound connections to reach the cluster services.
---

# kubernetes_ingress

The resource provides mechanisms to provide a collection of rules that allow inbound connections to reach the cluster services.
Ingress can be configured to give services externally-reachable urls, load balance traffic, terminate SSL, offer name based virtual hosting etc. 
The resource will by default create a ingress in the specified (or default) namespace.

~> Read more about ingress : https://kubernetes.io/docs/concepts/services-networking/ingress/

## Example Usage

```
resource "kubernetes_ingress" "example" {
   metadata {
    labels {
      app = "service_label"
    }
    name = "ingress-name"
  }
  spec {
    rules =[
      {
        host = "foo.bar.com"
        http {
          paths = [
            {
              path = "/foo"
              backend {
                service_name = "echoheaders-x"
                service_port = 80
              }
            }
          ]
        }
      },
      {
        host= "car.baz.com"
        http {
          paths = [
            {
              path = "/bar"
              backend {
                service_name = "echoheaders-y"
                service_port = 80
              }
            }
          ]
        }
      }
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `metadata` - (Required) Standard ingress metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `spec` - (Required) The Spec is the desired state of the Ingress. 

## Nested Blocks

### `metadata`

#### Arguments

* `annotations` - (Optional) An unstructured key value map stored with the ingress that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
* `generate_name` - (Optional) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#idempotency
* `labels` - (Optional) Map of string keys and values that can be used to organize and categorize (scope and select) the ingress. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
* `name` - (Optional) Name of the ingress, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
* `namespace` - (Optional) Namespace defines the space within which name of the ingress must be unique.

### `spec`

#### Arguments

* `backend` - (Optional) A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified. More info: https://kubernetes.io/docs/federation/api-reference/extensions/v1beta1/definitions/#_v1beta1_ingressbackend
* `rules` - (Required) A list of host rules used to configure the Ingress. If unspecified, or no rule matches, all traffic is sent to the default backend. Read more: https://kubernetes.io/docs/federation/api-reference/extensions/v1beta1/definitions/#_v1beta1_ingressrule
* `tls` - (Optional) TLS configuration. Currently the Ingress only supports a single TLS port, 443.. More info:https://kubernetes.io/docs/federation/api-reference/extensions/v1beta1/definitions/#_v1beta1_ingresstls

### `backend`

#### Arguments

* `service_name` - (Required) Specifies the name of the referenced service.
* `service_port` - (Required) Specifies the port of the referenced service.

### `rules`

#### Arguments

* `host` - (Optional) Host is the fully qualified domain name of a network host, as defined by RFC 3986.
* `http` - (Required)

### `http`

#### Arguments

* `paths` - (Required) A collection of paths that map requests to backends.

### `paths`

#### Arguments

* `backend` - (Optional) Backend defines the referenced service endpoint to which the traffic will be forwarded to.
* `path` - (Optional) Path is an extended POSIX regex as defined by IEEE Std 1003.1, (i.e this follows the egrep/unix syntax, not the perl syntax) matched against the path of an incoming request.

### `tls`

#### Arguments

* `hosts` - (Optional) osts are a list of hosts included in the TLS certificate. The values in this list must match the name/s used in the tlsSecret. 
* `secret_name ` - (Optional) SecretName is the name of the secret used to terminate SSL traffic on 443. Field is left optional to allow SSL routing based on SNI hostname alone. If the SNI host in a listener conflicts with the "Host" header field used by an IngressRule, the SNI host is used for termination and value of the Host header is used for routing.


## Import

Ingress can be imported using its name, e.g.

```
$ terraform import kubernetes_ingress.example ingress-name
```


