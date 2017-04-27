package kubernetes

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/util/intstr"
)

func expandBackend(b map[string]interface{}) *v1beta1.IngressBackend {
	backend := v1beta1.IngressBackend{}

	if name, ok := b["service_name"].(string); ok {
		backend.ServiceName = name

	}
	if port, ok := b["service_port"].(string); ok {
		i, err := strconv.Atoi(port)
		if err != nil {
			backend.ServicePort = intstr.IntOrString{
				Type:   intstr.String,
				StrVal: port,
			}
		} else {
			backend.ServicePort = intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(i),
			}
		}

	}

	return &backend
}

func expandIngressHTTPPath(path map[string]interface{}) v1beta1.HTTPIngressPath {
	ingressPath := v1beta1.HTTPIngressPath{}

	if bknd := path["backend"]; bknd != nil {
		b := bknd.([]interface{})[0].(map[string]interface{})
		ingressPath.Backend = *expandBackend(b)
	}
	if path, ok := path["path"].(string); ok {
		ingressPath.Path = path
	}

	return ingressPath
}

func expandIngressRule(rule map[string]interface{}) v1beta1.IngressRule {
	ingressRule := v1beta1.IngressRule{}

	if host, ok := rule["host"].(string); ok {
		ingressRule.Host = host
	}
	if http, ok := rule["http"].([]interface{}); ok {
		ingressRule.HTTP = expandHTTPIngressRuleValue(http)
	}

	return ingressRule
}

func expandHTTPIngressRuleValue(http []interface{}) *v1beta1.HTTPIngressRuleValue {
	httpIngressRuleValue := v1beta1.HTTPIngressRuleValue{}
	httpMap := http[0].(map[string]interface{})
	if paths, ok := httpMap["paths"].([]interface{}); ok {
		ps := make([]v1beta1.HTTPIngressPath, len(paths))
		for i, path := range paths {
			ps[i] = expandIngressHTTPPath(path.(map[string]interface{}))
		}
		httpIngressRuleValue.Paths = ps
	}
	return &httpIngressRuleValue
}

func expandIngressTLS(rule map[string]interface{}) v1beta1.IngressTLS {
	ingressTLS := v1beta1.IngressTLS{}

	if hosts, ok := rule["hosts"].(*schema.Set); ok {
		ingressTLS.Hosts = schemaSetToStringArray(hosts)
	}
	if secretName, ok := rule["secret_name"].(string); ok {
		ingressTLS.SecretName = secretName
	}

	return ingressTLS

}
func expandIngressSpec(in []interface{}) v1beta1.IngressSpec {
	ingressSpec := v1beta1.IngressSpec{}
	if len(in) < 1 {
		return ingressSpec
	}
	s := in[0].(map[string]interface{})

	if len(s["backend"].([]interface{})) >= 1 {
		bMap := s["backend"].([]interface{})[0].(map[string]interface{})
		ingressSpec.Backend = expandBackend(bMap)
	}

	rulesArray := s["rules"].([]interface{})
	rules := make([]v1beta1.IngressRule, len(rulesArray))
	for i, r := range rulesArray {
		rules[i] = expandIngressRule(r.(map[string]interface{}))
	}
	ingressSpec.Rules = rules
	if s["tls"] != nil {
		tlsArray := s["tls"].([]interface{})
		tls := make([]v1beta1.IngressTLS, len(tlsArray))
		for i, r := range tlsArray {
			tls[i] = expandIngressTLS(r.(map[string]interface{}))
		}
		ingressSpec.TLS = tls
	}
	return ingressSpec
}

// Flatteners
func flattenIngressSpec(in v1beta1.IngressSpec) []interface{} {
	att := make(map[string]interface{})
	if in.Backend != nil {
		att["backend"] = flattenBackend(in.Backend)
	}
	if in.Rules != nil {
		att["rules"] = flattenRules(in.Rules)
	}

	if in.TLS != nil {
		att["tls"] = flattenTLS(in.TLS)
	}
	return []interface{}{att}
}

func flattenBackend(in *v1beta1.IngressBackend) []interface{} {
	att := make(map[string]interface{})

	att["service_name"] = in.ServiceName
	att["service_port"] = in.ServicePort

	return []interface{}{att}
}

func flattenRules(in []v1beta1.IngressRule) []interface{} {
	att := make([]map[string]interface{}, len(in))
	for i, r := range in {
		m := map[string]interface{}{}
		m["host"] = r.Host
		m["http"] = flattenHTTPIngressRuleValue(r.HTTP)
		att[i] = m
	}
	return []interface{}{att}
}

func flattenHTTPIngressRuleValue(in *v1beta1.HTTPIngressRuleValue) []interface{} {
	att := make(map[string]interface{})
	att["paths"] = flattenHTTPIngressPath(in.Paths)
	return []interface{}{att}
}

func flattenHTTPIngressPath(in []v1beta1.HTTPIngressPath) []interface{} {
	att := make([]map[string]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		m["backend"] = flattenBackend(&v.Backend)
		if v.Path != "" {
			m["path"] = v.Path
		}
		att[i] = m
	}
	return []interface{}{att}
}

func flattenTLS(in []v1beta1.IngressTLS) []interface{} {

	att := make([]map[string]interface{}, len(in))
	for i, t := range in {
		m := map[string]interface{}{}
		m["hosts"] = t.Hosts
		m["secret_name"] = t.SecretName
		att[i] = m
	}
	return []interface{}{att}
}

func patchIngressSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)
	prefix += ".0."

	if d.HasChange(prefix + "backend") {
		v := d.Get(prefix + "backend").([]interface{})[0].(map[string]interface{})
		backend := expandBackend(v)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/backend",
			Value: backend,
		})
	}
	if d.HasChange(prefix + "rules") {
		rulesArray := d.Get(prefix + "rules").([]interface{})
		rules := make([]v1beta1.IngressRule, len(rulesArray))
		for i, r := range rulesArray {
			rules[i] = expandIngressRule(r.(map[string]interface{}))
		}
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/rules",
			Value: rules,
		})
	}

	if d.HasChange(prefix + "tls") {
		tlsArray := d.Get(prefix + "tls").([]interface{})
		tls := make([]v1beta1.IngressTLS, len(tlsArray))
		for i, r := range tlsArray {
			tls[i] = expandIngressTLS(r.(map[string]interface{}))
		}
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/tls",
			Value: tls,
		})
	}

	return ops, nil
}
