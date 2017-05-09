package kubernetes

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Flatteners

// userSpec is the  pod spec as specified in the user tf file configuration
// to compare with the actual PodSpec received from the API
func flattenPodSpec(in v1.PodSpec, userSpec v1.PodSpec) []interface{} {
	att := make(map[string]interface{})
	if in.ActiveDeadlineSeconds != nil {
		att["active_deadline_seconds"] = *in.ActiveDeadlineSeconds
	}
	att["containers"] = flattenContainers(in.Containers)

	att["dns_policy"] = in.DNSPolicy

	att["host_ipc"] = in.HostIPC
	att["host_network"] = in.HostNetwork
	att["host_pid"] = in.HostPID

	if in.Hostname != "" {
		att["hostname"] = in.Hostname
	}
	att["image_pull_secrets"] = flattenLocalObjectReferenceArray(in.ImagePullSecrets)

	if in.NodeName != "" {
		att["node_name"] = in.NodeName
	}
	if len(in.NodeSelector) > 0 {
		att["node_selector"] = in.NodeSelector
	}
	if in.RestartPolicy != "" {
		att["restart_policy"] = in.RestartPolicy
	}

	if in.SecurityContext != nil {
		att["security_context"] = flattenPodSecurityContext(in.SecurityContext)
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

	volume := userSpec.Volumes
	if len(in.Volumes) > 0 {
		att["volumes"] = flattenVolumes(in.Volumes, volume)
	}
	return []interface{}{att}
}

func flattenPodSecurityContext(in *v1.PodSecurityContext) []interface{} {
	att := make(map[string]interface{})
	if in.FSGroup != nil {
		att["fs_group"] = *in.FSGroup
	}

	if in.RunAsNonRoot != nil {
		att["run_as_non_root"] = *in.RunAsNonRoot
	}

	if in.RunAsUser != nil {
		att["run_as_user"] = *in.RunAsUser
	}

	if len(in.SupplementalGroups) > 0 {
		att["supplemental_groups"] = newInt64Set(schema.HashSchema(&schema.Schema{
			Type: schema.TypeInt,
		}), in.SupplementalGroups)
	}
	if in.SELinuxOptions != nil {
		att["se_linux_options"] = flattenSeLinuxOptions(in.SELinuxOptions)
	}
	return []interface{}{att}
}

func flattenSeLinuxOptions(in *v1.SELinuxOptions) []interface{} {
	att := make(map[string]interface{})
	if in.User != "" {
		att["user"] = in.User
	}
	if in.Role != "" {
		att["role"] = in.Role
	}
	if in.User != "" {
		att["type"] = in.Type
	}
	if in.Level != "" {
		att["level"] = in.Level
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
			//Found a volume which was not provided by user in the tf configuration file and must be internal
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
		if v.GCEPersistentDisk != nil {
			obj["gce_persistent_disk"] = flattenGCEPersistentDiskVolumeSource(v.GCEPersistentDisk)
		}
		if v.AWSElasticBlockStore != nil {
			obj["aws_elastic_block_store"] = flattenAWSElasticBlockStoreVolumeSource(v.AWSElasticBlockStore)
		}
		if v.HostPath != nil {
			obj["host_path"] = flattenHostPathVolumeSource(v.HostPath)
		}
		if v.Glusterfs != nil {
			obj["glusterfs"] = flattenGlusterfsVolumeSource(v.Glusterfs)
		}
		if v.NFS != nil {
			obj["nfs"] = flattenNFSVolumeSource(v.NFS)
		}
		if v.RBD != nil {
			obj["rbd"] = flattenRBDVolumeSource(v.RBD)
		}
		if v.ISCSI != nil {
			obj["iscsi"] = flattenISCSIVolumeSource(v.ISCSI)
		}
		if v.Cinder != nil {
			obj["cinder"] = flattenCinderVolumeSource(v.Cinder)
		}
		if v.CephFS != nil {
			obj["ceph_fs"] = flattenCephFSVolumeSource(v.CephFS)
		}
		if v.FC != nil {
			obj["fc"] = flattenFCVolumeSource(v.FC)
		}
		if v.Flocker != nil {
			obj["flocker"] = flattenFlockerVolumeSource(v.Flocker)
		}
		if v.FlexVolume != nil {
			obj["flex_volume"] = flattenFlexVolumeSource(v.FlexVolume)
		}
		if v.AzureFile != nil {
			obj["azure_file"] = flattenAzureFileVolumeSource(v.AzureFile)
		}
		if v.VsphereVolume != nil {
			obj["vsphere_volume"] = flattenVsphereVirtualDiskVolumeSource(v.VsphereVolume)
		}
		if v.Quobyte != nil {
			obj["quobyte"] = flattenQuobyteVolumeSource(v.Quobyte)
		}
		if v.AzureDisk != nil {
			obj["azure_disk"] = flattenAzureDiskVolumeSource(v.AzureDisk)
		}
		if v.PhotonPersistentDisk != nil {
			obj["photon_persistent_disk"] = flattenPhotonPersistentDiskVolumeSource(v.PhotonPersistentDisk)
		}
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

	if v, ok := in["active_deadline_seconds"].(int); ok && v > 0 {
		obj.ActiveDeadlineSeconds = ptrToInt64(int64(v))
	}

	if v, ok := in["containers"].([]interface{}); ok && len(v) > 0 {
		cs, err := expandContainers(v)
		if err != nil {
			return obj, err
		}
		obj.Containers = cs
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

	if v, ok := in["host_pid"]; ok {
		obj.HostPID = v.(bool)
	}

	if v, ok := in["hostname"]; ok {
		obj.Hostname = v.(string)
	}

	if v, ok := in["image_pull_secrets"].([]interface{}); ok {
		cs := expandLocalObjectReferenceArray(v)
		obj.ImagePullSecrets = cs
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

	if v, ok := in["security_context"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityContext = expandPodSecurityContext(v)
	}

	if v, ok := in["service_account_name"].(string); ok {
		obj.ServiceAccountName = v
	}

	if v, ok := in["subdomain"].(string); ok {
		obj.Subdomain = v
	}

	if v, ok := in["termination_grace_period_seconds"].(int); ok {
		obj.TerminationGracePeriodSeconds = ptrToInt64(int64(v))
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

func expandPodSecurityContext(l []interface{}) *v1.PodSecurityContext {
	if len(l) == 0 || l[0] == nil {
		return &v1.PodSecurityContext{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.PodSecurityContext{}
	if v, ok := in["fs_group"].(int); ok {
		obj.FSGroup = ptrToInt64(int64(v))
	}
	if v, ok := in["run_as_non_root"].(bool); ok {
		obj.RunAsNonRoot = ptrToBool(v)
	}
	if v, ok := in["run_as_user"].(int); ok {
		obj.RunAsUser = ptrToInt64(int64(v))
	}
	if v, ok := in["supplemental_groups"].(*schema.Set); ok {
		obj.SupplementalGroups = schemaSetToInt64Array(v)
	}

	if v, ok := in["se_linux_options"].([]interface{}); ok && len(v) > 0 {
		obj.SELinuxOptions = expandSeLinuxOptions(v)
	}

	return obj
}

func expandSeLinuxOptions(l []interface{}) *v1.SELinuxOptions {
	if len(l) == 0 || l[0] == nil {
		return &v1.SELinuxOptions{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.SELinuxOptions{}
	if v, ok := in["level"]; ok {
		obj.Level = v.(string)
	}
	if v, ok := in["role"]; ok {
		obj.Role = v.(string)
	}
	if v, ok := in["type"]; ok {
		obj.Type = v.(string)
	}
	if v, ok := in["user"]; ok {
		obj.User = v.(string)
	}
	return obj
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
		m := c.(map[string]interface{})

		if value, ok := m["name"]; ok {
			vl[i].Name = value.(string)
		}
		if value, ok := m["persistent_volume_claim"].([]interface{}); ok && len(value) > 0 {
			vl[i].PersistentVolumeClaim = expandPersistentVolumeClaimVolumeSource(value)
		}
		if value, ok := m["secret"].([]interface{}); ok && len(value) > 0 {
			vl[i].Secret = expandSecretVolumeSource(value)
		}
		if v, ok := m["gce_persistent_disk"].([]interface{}); ok && len(v) > 0 {
			vl[i].GCEPersistentDisk = expandGCEPersistentDiskVolumeSource(v)
		}
		if v, ok := m["aws_elastic_block_store"].([]interface{}); ok && len(v) > 0 {
			vl[i].AWSElasticBlockStore = expandAWSElasticBlockStoreVolumeSource(v)
		}
		if v, ok := m["host_path"].([]interface{}); ok && len(v) > 0 {
			vl[i].HostPath = expandHostPathVolumeSource(v)
		}
		if v, ok := m["glusterfs"].([]interface{}); ok && len(v) > 0 {
			vl[i].Glusterfs = expandGlusterfsVolumeSource(v)
		}
		if v, ok := m["nfs"].([]interface{}); ok && len(v) > 0 {
			vl[i].NFS = expandNFSVolumeSource(v)
		}
		if v, ok := m["rbd"].([]interface{}); ok && len(v) > 0 {
			vl[i].RBD = expandRBDVolumeSource(v)
		}
		if v, ok := m["iscsi"].([]interface{}); ok && len(v) > 0 {
			vl[i].ISCSI = expandISCSIVolumeSource(v)
		}
		if v, ok := m["cinder"].([]interface{}); ok && len(v) > 0 {
			vl[i].Cinder = expandCinderVolumeSource(v)
		}
		if v, ok := m["ceph_fs"].([]interface{}); ok && len(v) > 0 {
			vl[i].CephFS = expandCephFSVolumeSource(v)
		}
		if v, ok := m["fc"].([]interface{}); ok && len(v) > 0 {
			vl[i].FC = expandFCVolumeSource(v)
		}
		if v, ok := m["flocker"].([]interface{}); ok && len(v) > 0 {
			vl[i].Flocker = expandFlockerVolumeSource(v)
		}
		if v, ok := m["flex_volume"].([]interface{}); ok && len(v) > 0 {
			vl[i].FlexVolume = expandFlexVolumeSource(v)
		}
		if v, ok := m["azure_file"].([]interface{}); ok && len(v) > 0 {
			vl[i].AzureFile = expandAzureFileVolumeSource(v)
		}
		if v, ok := m["vsphere_volume"].([]interface{}); ok && len(v) > 0 {
			vl[i].VsphereVolume = expandVsphereVirtualDiskVolumeSource(v)
		}
		if v, ok := m["quobyte"].([]interface{}); ok && len(v) > 0 {
			vl[i].Quobyte = expandQuobyteVolumeSource(v)
		}
		if v, ok := m["azure_disk"].([]interface{}); ok && len(v) > 0 {
			vl[i].AzureDisk = expandAzureDiskVolumeSource(v)
		}
		if v, ok := m["photon_persistent_disk"].([]interface{}); ok && len(v) > 0 {
			vl[i].PhotonPersistentDisk = expandPhotonPersistentDiskVolumeSource(v)
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
