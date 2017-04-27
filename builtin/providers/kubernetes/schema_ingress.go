package kubernetes

import "github.com/hashicorp/terraform/helper/schema"

func ingressBackendFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the name of the referenced service.",
		},
		"service_port": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the port of the referenced service.",
		},
	}
}

func ingressHTTPPathFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backend": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Backend defines the referenced service endpoint to which the traffic will be forwarded to.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: ingressBackendFields(),
			},
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path is an extended POSIX regex as defined by IEEE Std 1003.1,matched against the path of an incoming request. Currently it can contain characters disallowed from the conventional 'path' part of a URL as defined by RFC 3986. Paths must begin with a '/'. If unspecified, the path defaults to a catch all sending traffic to the backend.",
		},
	}
}

func httpIngressRuleValueFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"paths": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A collection of paths that map requests to backends.",
			Elem: &schema.Resource{
				Schema: ingressHTTPPathFields(),
			},
		},
	}
}

func ingressRuleFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Host is the fully qualified domain name of a network host, as defined by RFC 3986.",
		},

		"http": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: httpIngressRuleValueFields(),
			},
		},
	}
}

func ingressTLSFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hosts": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Hosts are a list of hosts included in the TLS certificate. The values in this list must match the name/s used in the tlsSecret. Defaults to the wildcard host setting for the loadbalancer controller fulfilling this Ingress, if left unspecified.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
		},
		"secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SecretName is the name of the secret used to terminate SSL traffic on 443. Field is left optional to allow SSL routing based on SNI hostname alone. If the SNI host in a listener conflicts with the 'Host' header field used by an IngressRule, the SNI host is used for termination and value of the Host header is used for routing.",
		},
	}
}

func ingressSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backend": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified. This field is optional to allow the loadbalancer controller or defaulting logic to specify a global default.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: ingressBackendFields(),
			},
		},

		"rules": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A list of host rules used to configure the Ingress. If unspecified, or no rule matches, all traffic is sent to the default backend.",
			Elem: &schema.Resource{
				Schema: ingressRuleFields(),
			},
		},

		"tls": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "TLS configuration. Currently the Ingress only supports a single TLS port, 443. If multiple members of this list specify different hosts, they will be multiplexed on the same port according to the hostname specified through the SNI TLS extension, if the ingress controller fulfilling the ingress supports SNI.",
			Elem: &schema.Resource{
				Schema: ingressTLSFields(),
			},
		},
	}
}
