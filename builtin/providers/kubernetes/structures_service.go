package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/util/intstr"
)

// Flatteners
func flattenServiceSpec(in v1.ServiceSpec) []interface{} {
	att := make(map[string]interface{})

	if in.ClusterIP != "" {
		att["cluster_ip"] = in.ClusterIP
	}

	if len(in.ExternalIPs) > 0 {
		att["external_ips"] = newStringSet(schema.HashString, in.ExternalIPs)
	}

	if in.ExternalName != "" {
		att["external_name"] = in.ExternalName
	}

	if in.LoadBalancerIP != "" {
		att["load_balancer_ip"] = in.LoadBalancerIP
	}

	if len(in.LoadBalancerSourceRanges) > 0 {
		att["load_balancer_source_ranges"] = newStringSet(schema.HashString, in.LoadBalancerSourceRanges)
	}

	if in.SessionAffinity != "" {
		att["session_affinity"] = in.SessionAffinity
	}

	if len(in.Selector) > 0 {
		att["selector"] = in.Selector
	}

	if in.Type != "" {
		att["type"] = in.Type
	}

	att["ports"] = flattenServicePorts(in.Ports)

	return []interface{}{att}
}

func flattenServicePorts(in []v1.ServicePort) []interface{} {
	att := make([]map[string]interface{}, len(in))

	for i, v := range in {
		m := map[string]interface{}{}

		if v.Name != "" {
			m["name"] = v.Name
		}

		m["node_port"] = v.NodePort

		m["port"] = v.Port

		if v.Protocol != "" {
			m["protocol"] = v.Protocol
		}

		m["target_port"] = v.TargetPort

		att[i] = m
	}
	return []interface{}{att}
}

//expanders
func expandServiceSpec(d []interface{}) (v1.ServiceSpec, error) {
	if len(d) == 0 || d[0] == nil {
		return v1.ServiceSpec{}, nil
	}
	in := d[0].(map[string]interface{})
	obj := v1.ServiceSpec{}

	if clusterIP, ok := in["cluster_ip"].(string); ok {
		obj.ClusterIP = clusterIP
	}

	if externalIPs, ok := in["external_ips"]; ok {
		obj.ExternalIPs = schemaSetToStringArray(externalIPs.(*schema.Set))
	}

	if externalName, ok := in["external_name"].(string); ok {
		obj.ExternalName = externalName
	}

	if loadBalancerIP, ok := in["load_balancer_ip"].(string); ok {
		obj.LoadBalancerIP = loadBalancerIP
	}

	if loadBalancerSourceRanges, ok := in["load_balancer_source_ranges"]; ok {
		obj.LoadBalancerSourceRanges = schemaSetToStringArray(loadBalancerSourceRanges.(*schema.Set))
	}

	if sessionAffinity, ok := in["session_affinity"]; ok {
		obj.SessionAffinity = v1.ServiceAffinity(sessionAffinity.(string))
	}

	if selector, ok := in["selector"]; ok {
		obj.Selector = expandStringMap(selector.(map[string]interface{}))
	}

	if typeValue, ok := in["type"]; ok {
		obj.Type = v1.ServiceType(typeValue.(string))
	}

	if ports, ok := in["ports"].([]interface{}); ok && len(ports) > 0 {
		var err error
		obj.Ports, err = expandServicePort(ports)
		if err != nil {
			return obj, err
		}
	}

	return obj, nil
}

func expandServicePort(in []interface{}) ([]v1.ServicePort, error) {
	if len(in) == 0 {
		return []v1.ServicePort{}, nil
	}
	ports := make([]v1.ServicePort, len(in))
	for i, c := range in {
		p := c.(map[string]interface{})

		if name, ok := p["name"]; ok {
			ports[i].Name = name.(string)
		}

		if nodePort, ok := p["node_port"]; ok {
			ports[i].NodePort = int32(nodePort.(int))
		}

		if port, ok := p["port"]; ok {
			ports[i].Port = int32(port.(int))
		}

		if protocol, ok := p["protocol"]; ok {
			ports[i].Protocol = v1.Protocol(protocol.(string))
		}

		if targetPort, ok := p["target_port"]; ok {
			ports[i].TargetPort = intstr.IntOrString{
				IntVal: int32(targetPort.(int)),
			}
		}
	}
	return ports, nil
}

func patchServiceSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)
	prefix += ".0."

	if d.HasChange(prefix + "external_ips") {
		v := d.Get(prefix + "external_ips").(*schema.Set)
		externalIPs := schemaSetToStringArray(v)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/externalIPs",
			Value: externalIPs,
		})
	}
	if d.HasChange(prefix + "external_name") {
		externalName := d.Get(prefix + "external_name").(string)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/externalName",
			Value: externalName,
		})

	}

	if d.HasChange(prefix + "load_balancer_ip") {
		loadBalancerIP := d.Get(prefix + "load_balancer_ip").(string)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/loadBalancerIP",
			Value: loadBalancerIP,
		})
	}

	if d.HasChange(prefix + "load_balancer_source_ranges") {
		v := d.Get(prefix + "load_balancer_source_ranges").(*schema.Set)
		loadBalancerSourceRanges := schemaSetToStringArray(v)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/loadBalancerSourceRanges",
			Value: loadBalancerSourceRanges,
		})
	}

	if d.HasChange(prefix + "session_affinity") {
		sessionAffinity := v1.ServiceAffinity(d.Get(prefix + "session_affinity").(string))
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/sessionAffinity",
			Value: sessionAffinity,
		})
	}
	if d.HasChange(prefix + "type") {
		typeValue := v1.ServiceType(d.Get(prefix + "type").(string))
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/type",
			Value: typeValue,
		})
	}
	if d.HasChange(prefix + "selector") {
		selector := d.Get(prefix + "selector").(map[string]interface{})
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/selector",
			Value: expandStringMap(selector),
		})
	}
	if d.HasChange(prefix + "ports") {
		v := d.Get(prefix + "ports").([]interface{})
		ports, err := expandServicePort(v)
		if err != nil {
			return ops, err
		}
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/ports",
			Value: ports,
		})
	}

	return ops, nil
}
