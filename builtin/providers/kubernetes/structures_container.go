package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/v1"
)

func flattenContainerVolumeMounts(in []v1.VolumeMount) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		m["mount_path"] = v.MountPath
		m["name"] = v.Name
		m["sub_path"] = v.SubPath
		m["read_only"] = v.ReadOnly

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
			c["command"] = newStringSet(schema.HashString, v.Command)
		}
		if len(v.Args) > 0 {
			c["args"] = newStringSet(schema.HashString, v.Args)
		}

		c["image_pull_policy"] = v.ImagePullPolicy
		c["working_dir"] = v.WorkingDir

		c["resources"] = flattenResourceRequirements(v.Resources)
		c["ports"] = flattenContainerPorts(v.Ports)
		if len(v.VolumeMounts) > 0 {
			c["volume_mounts"] = flattenContainerVolumeMounts(v.VolumeMounts)
		}

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
		if command, ok := ctr["command"]; ok {
			cs[i].Command = schemaSetToStringArray(command.(*schema.Set))
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
