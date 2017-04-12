package kubernetes

import (
	"k8s.io/kubernetes/pkg/api/v1"
)

// Flatteners
func flattenPodSpec(in v1.PodSpec) []interface{} {
	att := make(map[string]interface{})
	if in.ActiveDeadlineSeconds != nil {
		att["active_deadline_seconds"] = *in.ActiveDeadlineSeconds
	}
	if in.DNSPolicy != "" {
		att["dns_policy"] = in.DNSPolicy
	}
	if in.HostIPC {
		att["host_ipc"] = in.HostIPC
	}
	if in.HostNetwork {
		att["host_network"] = in.HostNetwork
	}
	if in.Hostname != "" {
		att["hostname"] = in.Hostname
	}
	if in.NodeName != "" {
		att["node_name"] = in.NodeName
	}
	if len(in.NodeSelector) > 0 {
		att["node_selector"] = in.NodeSelector
	}
	if in.RestartPolicy != "" {
		att["restart_policy"] = in.RestartPolicy
	}
	if in.ServiceAccountName != "" {
		att["service_account_name"] = in.ServiceAccountName
	}
	if in.Subdomain != "" {
		att["subdomain"] = in.Subdomain
	}

	if len(in.NodeSelector) > 0 {
		att["node_selector"] = in.NodeSelector
	}
	if in.TerminationGracePeriodSeconds != nil {
		att["termination_grace_period_seconds"] = *in.TerminationGracePeriodSeconds
	}
	att["image_pull_secrets"] = flattenLocalObjectReferenceArray(in.ImagePullSecrets)
	att["containers"] = flattenContainers(in.Containers)

	if len(in.Volumes) > 0 {
		att["volumes"] = flattenVolumes(in.Volumes)
	}
	return []interface{}{att}
}

// Expanders
func expandPodSpec(p []interface{}) (v1.PodSpec, error) {
	obj := v1.PodSpec{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["active_deadline_seconds"].(*int64); ok {
		obj.ActiveDeadlineSeconds = v
	}

	if v, ok := in["dns_policy"].(string); ok {
		obj.DNSPolicy = v1.DNSPolicy(v)
	}
	if v, ok := in["host_ipc"].(bool); ok {
		obj.HostIPC = v
	}
	if v, ok := in["host_network"].(bool); ok {
		obj.HostNetwork = v
	}
	if v, ok := in["hostname"].(string); ok {
		obj.Hostname = v
	}
	if v, ok := in["node_name"].(string); ok {
		obj.NodeName = v
	}
	if v, ok := in["node_selector"].(map[string]interface{}); ok {
		obj.NodeSelector = expandStringMap(v)
	}
	if v, ok := in["restart_policy"].(string); ok {
		obj.RestartPolicy = v1.RestartPolicy(v)
	}
	if v, ok := in["service_account_name"].(string); ok {
		obj.ServiceAccountName = v
	}
	if v, ok := in["termination_grace_period_seconds"].(*int64); ok {
		obj.TerminationGracePeriodSeconds = v
	}

	if v, ok := in["containers"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandContainers(v)
		if err != nil {
			return obj, err
		}
		obj.Containers = cs
	}
	if v, ok := in["volumes"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandVolumes(v)
		if err != nil {
			return obj, err
		}
		obj.Volumes = cs
	}

	return obj, nil
}
