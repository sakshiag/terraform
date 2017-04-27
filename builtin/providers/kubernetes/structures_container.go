package kubernetes

import "k8s.io/kubernetes/pkg/api/v1"

func flattenConfigMapKeyRef(in *v1.ConfigMapKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	return []interface{}{att}
}

func flattenFieldRef(in *v1.ObjectFieldSelector) []interface{} {
	att := make(map[string]interface{})

	if in.APIVersion != "" {
		att["api_version"] = in.APIVersion
	}
	if in.FieldPath != "" {
		att["field_path"] = in.FieldPath
	}
	return []interface{}{att}
}

func flattenResourceFieldRef(in *v1.ResourceFieldSelector) []interface{} {
	att := make(map[string]interface{})

	if in.ContainerName != "" {
		att["container_name"] = in.ContainerName
	}
	if in.Resource != "" {
		att["resource"] = in.Resource
	}
	return []interface{}{att}
}

func flattenSecretKeyRef(in *v1.SecretKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	return []interface{}{att}
}

func flattenValueFrom(in *v1.EnvVarSource) []interface{} {
	att := make(map[string]interface{})

	if in.ConfigMapKeyRef != nil {
		att["config_map_key_ref"] = flattenConfigMapKeyRef(in.ConfigMapKeyRef)
	}
	if in.ResourceFieldRef != nil {
		att["resource_field_ref"] = flattenResourceFieldRef(in.ResourceFieldRef)
	}
	if in.SecretKeyRef != nil {
		att["secret_key_ref"] = flattenSecretKeyRef(in.SecretKeyRef)
	}
	if in.FieldRef != nil {
		att["field_ref"] = flattenFieldRef(in.FieldRef)
	}
	return []interface{}{att}
}

func flattenContainerVolumeMounts(in []v1.VolumeMount) []interface{} {
	att := make([]map[string]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		if v.MountPath != "" {
			m["mount_path"] = v.MountPath
		}
		if v.Name != "" {
			m["name"] = v.Name
		}
		if v.SubPath != "" {
			m["sub_path"] = v.SubPath

		}

		att[i] = m
	}
	return []interface{}{att}
}

func flattenContainerEnvs(in []v1.EnvVar) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		if v.Name != "" {
			m["name"] = v.Name
		}
		if v.Value != "" {
			m["value"] = v.Value
		}
		if v.ValueFrom != nil {
			m["value_from"] = flattenValueFrom(v.ValueFrom)
		}

		att[i] = m
	}
	return att
}

func flattenContainerPorts(in []v1.ContainerPort) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		m["container_port"] = v.ContainerPort
		if v.HostIP != "" {
			m["host_ip"] = v.HostIP
		}
		m["host_port"] = v.HostPort
		if v.Name != "" {
			m["name"] = v.Name
		}
		if v.Protocol != "" {
			m["protocol"] = v.Protocol
		}
		att[i] = m
	}
	return att
}

func flattenContainers(in []v1.Container) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		c := make(map[string]interface{})
		c["image"] = v.Image
		c["name"] = v.Name
		if len(v.Command) > 0 {
			c["command"] = v.Command
		}
		if len(v.Args) > 0 {
			c["args"] = v.Args
		}

		c["image_pull_policy"] = v.ImagePullPolicy
		c["working_dir"] = v.WorkingDir

		if len(v.Resources.Limits) == 0 && len(v.Resources.Requests) == 0 {
			c["resources"] = []interface{}{}
		} else {
			c["resources"] = flattenResourceRequirements(v.Resources)
		}

		c["ports"] = flattenContainerPorts(v.Ports)
		c["env"] = flattenContainerEnvs(v.Env)
		att[i] = c
	}
	return att
}

func expandContainers(ctrs []interface{}) ([]v1.Container, error) {
	if len(ctrs) == 0 {
		return []v1.Container{}, nil
	}
	cs := make([]v1.Container, len(ctrs))
	for i, c := range ctrs {
		ctr := c.(map[string]interface{})

		if image, ok := ctr["image"]; ok {
			cs[i].Image = image.(string)
		}
		if name, ok := ctr["name"]; ok {
			cs[i].Name = name.(string)
		}
		if command, ok := ctr["command"].([]interface{}); ok {
			cs[i].Command = expandStringSlice(command)
		}
		if args, ok := ctr["args"].([]interface{}); ok {
			cs[i].Args = expandStringSlice(args)
		}

		if v, ok := ctr["resources"].([]interface{}); ok && len(v) > 0 {

			var err error
			cs[i].Resources, err = expandResourceRequirements(v)
			if err != nil {
				return cs, err
			}
		}

		if v, ok := ctr["ports"].([]interface{}); ok && len(v) > 0 {
			var err error
			cs[i].Ports, err = expandContainerPort(v)
			if err != nil {
				return cs, err
			}
		}
		if v, ok := ctr["env"].([]interface{}); ok && len(v) > 0 {
			var err error
			cs[i].Env, err = expandContainerEnv(v)
			if err != nil {
				return cs, err
			}
		}

		if policy, ok := ctr["image_pull_policy"]; ok {
			cs[i].ImagePullPolicy = v1.PullPolicy(policy.(string))
		}

		if v, ok := ctr["volume_mounts"].([]interface{}); ok && len(v) > 0 {
			var err error
			cs[i].VolumeMounts, err = expandContainerVolumeMounts(v)
			if err != nil {
				return cs, err
			}
		}
	}
	return cs, nil
}

func expandContainerVolumeMounts(in []interface{}) ([]v1.VolumeMount, error) {
	if len(in) == 0 {
		return []v1.VolumeMount{}, nil
	}
	vmp := make([]v1.VolumeMount, len(in))
	for i, c := range in {
		p := c.(map[string]interface{})
		if mountPath, ok := p["mount_path"]; ok {
			vmp[i].MountPath = mountPath.(string)
		}
		if name, ok := p["name"]; ok {
			vmp[i].Name = name.(string)
		}
		if readOnly, ok := p["read_only"]; ok {
			vmp[i].ReadOnly = readOnly.(bool)
		}
		if subPath, ok := p["sub_path"]; ok {
			vmp[i].SubPath = subPath.(string)
		}
	}
	return vmp, nil
}

func expandContainerEnv(in []interface{}) ([]v1.EnvVar, error) {
	if len(in) == 0 {
		return []v1.EnvVar{}, nil
	}
	envs := make([]v1.EnvVar, len(in))
	for i, c := range in {
		p := c.(map[string]interface{})
		if name, ok := p["name"]; ok {
			envs[i].Name = name.(string)
		}
		if value, ok := p["value"]; ok {
			envs[i].Value = value.(string)
		}
		if v, ok := p["value_from"].([]interface{}); ok && len(v) > 0 {
			var err error
			envs[i].ValueFrom, err = expandEnvValueFrom(v)
			if err != nil {
				return envs, err
			}
		}
	}
	return envs, nil
}

func expandContainerPort(in []interface{}) ([]v1.ContainerPort, error) {
	if len(in) == 0 {
		return []v1.ContainerPort{}, nil
	}
	ports := make([]v1.ContainerPort, len(in))
	for i, c := range in {
		p := c.(map[string]interface{})
		if containerPort, ok := p["container_port"]; ok {
			ports[i].ContainerPort = int32(containerPort.(int))
		}
		if hostIP, ok := p["host_ip"]; ok {
			ports[i].HostIP = hostIP.(string)
		}
		if hostPort, ok := p["host_port"]; ok {
			ports[i].HostPort = int32(hostPort.(int))
		}
		if name, ok := p["name"]; ok {
			ports[i].Name = name.(string)
		}
		if protocol, ok := p["protocol"]; ok {
			ports[i].Protocol = v1.Protocol(protocol.(string))
		}
	}
	return ports, nil
}

func expandConfigMapKeyRef(r []interface{}) (*v1.ConfigMapKeySelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ConfigMapKeySelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ConfigMapKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	return obj, nil

}
func expandFieldRef(r []interface{}) (*v1.ObjectFieldSelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ObjectFieldSelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ObjectFieldSelector{}

	if v, ok := in["api_version"].(string); ok {
		obj.APIVersion = v
	}
	if v, ok := in["field_path"].(string); ok {
		obj.FieldPath = v
	}
	return obj, nil
}
func expandResourceFieldRef(r []interface{}) (*v1.ResourceFieldSelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ResourceFieldSelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ResourceFieldSelector{}

	if v, ok := in["container_name"].(string); ok {
		obj.ContainerName = v
	}
	if v, ok := in["resource"].(string); ok {
		obj.Resource = v
	}
	return obj, nil
}
func expandSecretKeyRef(r []interface{}) (*v1.SecretKeySelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.SecretKeySelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.SecretKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	return obj, nil
}

func expandEnvValueFrom(r []interface{}) (*v1.EnvVarSource, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.EnvVarSource{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.EnvVarSource{}

	var err error
	if v, ok := in["config_map_key_ref"].([]interface{}); ok && len(v) > 0 {
		obj.ConfigMapKeyRef, err = expandConfigMapKeyRef(v)
		if err != nil {
			return obj, err
		}
	}
	if v, ok := in["field_ref"].([]interface{}); ok && len(v) > 0 {
		obj.FieldRef, err = expandFieldRef(v)
		if err != nil {
			return obj, err
		}
	}
	if v, ok := in["secret_key_ref"].([]interface{}); ok && len(v) > 0 {
		obj.SecretKeyRef, err = expandSecretKeyRef(v)
		if err != nil {
			return obj, err
		}
	}
	if v, ok := in["resource_field_ref"].([]interface{}); ok && len(v) > 0 {
		obj.ResourceFieldRef, err = expandResourceFieldRef(v)
		if err != nil {
			return obj, err
		}
	}
	return obj, nil

}
