package kubernetes

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Flatteners

func flattenPodSpec(in v1.PodSpec, userSpec v1.PodSpec) []interface{} {
	att := make(map[string]interface{})
	if in.ActiveDeadlineSeconds != nil {
		att["active_deadline_seconds"] = *in.ActiveDeadlineSeconds
	}
	att["dns_policy"] = in.DNSPolicy
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

	if in.TerminationGracePeriodSeconds != nil {
		att["termination_grace_period_seconds"] = *in.TerminationGracePeriodSeconds
	}
	att["image_pull_secrets"] = flattenLocalObjectReferenceArray(in.ImagePullSecrets)
	att["containers"] = flattenContainers(in.Containers)

	//Always has one default one
	//TODO verify
	volume := userSpec.Volumes
	if len(in.Volumes) > 0 {
		att["volumes"] = flattenVolumes(in.Volumes, volume)
	}
	return []interface{}{att}
}

func flattenVolumes(volumes []v1.Volume, userVolume []v1.Volume) []interface{} {

	userVolumeNames := make(map[string]bool, len(userVolume))
	for _, v := range userVolume {
		userVolumeNames[v.Name] = true
	}

	att := make([]interface{}, len(userVolume))
	for i, v := range volumes {
		obj := map[string]interface{}{}

		if !userVolumeNames[v.Name] {
			continue
		}
		if v.Name != "" {
			obj["name"] = v.Name
		}
		if v.PersistentVolumeClaim != nil {
			obj["persistent_volume_claim"] = flattenPersistentVolumeClaimVolumeSource(v.PersistentVolumeClaim)
		}
		if v.Secret != nil {
			obj["secret"] = flattenSecretVolumeSource(v.Secret)
		}
		//More values needed here
		att[i] = obj
	}
	return att
}

func flattenPersistentVolumeClaimVolumeSource(in *v1.PersistentVolumeClaimVolumeSource) []interface{} {
	att := make(map[string]interface{})
	if in.ClaimName != "" {
		att["claim_name"] = in.ClaimName
	}
	if in.ReadOnly {
		att["read_only"] = in.ReadOnly
	}

	return []interface{}{att}
}

func flattenSecretVolumeSource(in *v1.SecretVolumeSource) []interface{} {
	att := make(map[string]interface{})
	if in.SecretName != "" {
		att["secret_name"] = in.SecretName
	}
	return []interface{}{att}
}

func flattenPodTemplateSpec(in v1.PodTemplateSpec, userSpec v1.PodTemplateSpec) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = flattenMetadata(in.ObjectMeta)
	att["spec"] = flattenPodSpec(in.Spec, userSpec.Spec)
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

	if v, ok := in["dns_policy"].(string); ok {
		obj.DNSPolicy = v1.DNSPolicy(v)
	}

	if v, ok := in["host_ipc"]; ok {
		obj.HostIPC = v.(bool)
	}

	if v, ok := in["host_network"]; ok {
		obj.HostNetwork = v.(bool)
	}

	if v, ok := in["hostname"]; ok {
		obj.Hostname = v.(string)
	}

	if v, ok := in["node_name"]; ok {
		obj.NodeName = v.(string)
	}

	if v, ok := in["node_selector"].(map[string]string); ok {
		obj.NodeSelector = v
	}

	if v, ok := in["restart_policy"].(string); ok {
		obj.RestartPolicy = v1.RestartPolicy(v)
	}

	if v, ok := in["service_account_name"].(string); ok {
		obj.ServiceAccountName = v
	}

	if v, ok := in["subdomain"].(string); ok {
		obj.Subdomain = v
	}

	if v, ok := in["termination_grace_period_seconds"].(*int64); ok {
		obj.TerminationGracePeriodSeconds = v
	}

	if v, ok := in["image_pull_secrets"].([]interface{}); ok {
		cs := expandLocalObjectReferenceArray(v)
		obj.ImagePullSecrets = cs
	}

	return obj, nil
}

func expandPersistentVolumeClaimVolumeSource(l []interface{}) *v1.PersistentVolumeClaimVolumeSource {
	if len(l) == 0 || l[0] == nil {
		return &v1.PersistentVolumeClaimVolumeSource{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.PersistentVolumeClaimVolumeSource{
		ClaimName: in["claim_name"].(string),
		ReadOnly:  in["read_only"].(bool),
	}
	return obj
}

func expandSecretVolumeSource(l []interface{}) *v1.SecretVolumeSource {
	if len(l) == 0 || l[0] == nil {
		return &v1.SecretVolumeSource{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.SecretVolumeSource{
		SecretName: in["secret_name"].(string),
	}
	return obj
}

func expandVolumes(volumes []interface{}) ([]v1.Volume, error) {
	if len(volumes) == 0 {
		return []v1.Volume{}, nil
	}
	vl := make([]v1.Volume, len(volumes))
	for i, c := range volumes {
		v := c.(map[string]interface{})

		if name, ok := v["name"]; ok {
			vl[i].Name = name.(string)
		}
		if pvc, ok := v["persistent_volume_claim"].([]interface{}); ok && len(pvc) > 0 {
			vl[i].PersistentVolumeClaim = expandPersistentVolumeClaimVolumeSource(pvc)
		}
		if secret, ok := v["secret"].([]interface{}); ok && len(secret) > 0 {
			vl[i].Secret = expandSecretVolumeSource(secret)
		}
	}
	return vl, nil
}

func expandPodTemplateSpec(p []interface{}) (v1.PodTemplateSpec, error) {
	obj := v1.PodTemplateSpec{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})
	meta := expandMetadata(in["metadata"].([]interface{}))
	obj.ObjectMeta = meta

	if v, ok := in["spec"].([]interface{}); ok && len(v) > 0 {
		var err error
		obj.Spec, err = expandPodSpec(v)
		if err != nil {
			return obj, err
		}
	}
	return obj, nil
}

func patchPodSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)
	prefix += ".0."

	if d.HasChange(prefix + "active_deadline_seconds") {

		v := d.Get(prefix + "active_deadline_seconds").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/activeDeadlineSeconds",
			Value: v,
		})
	}

	if d.HasChange(prefix + "containers") {
		containers := d.Get(prefix + "containers").([]interface{})
		value, _ := expandContainers(containers)

		for i, v := range value {
			ops = append(ops, &ReplaceOperation{
				Path:  pathPrefix + "/containers/" + strconv.Itoa(i) + "/image",
				Value: v.Image,
			})
		}

	}

	return ops, nil
}
