package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "k8s.io/kubernetes/pkg/api/v1"
)

func expandReplicationControllerSpec(in []interface{}) api.ReplicationControllerSpec {
	rcs := api.ReplicationControllerSpec{}
	if len(in) < 1 {
		return rcs
	}
	s := in[0].(map[string]interface{})
	if v, ok := s["replicas"]; ok {
		va := int32(v.(int))
		rcs.Replicas = &va
	}
	if v, ok := s["min_ready_seconds"]; ok {
		rcs.MinReadySeconds = int32(v.(int))
	}
	if v, ok := s["selector"]; ok {
		rcs.Selector = expandStringMap(v.(map[string]interface{}))
	}
	t := s["template"].([]interface{})
	pd, _ := expandPodTemplateSpec(t)
	rcs.Template = &pd
	return rcs
}

func flattenReplicationControllerSpec(in api.ReplicationControllerSpec, userSpec api.ReplicationControllerSpec) []interface{} {

	att := make(map[string]interface{})

	att["replicas"] = *in.Replicas
	att["min_ready_seconds"] = in.MinReadySeconds
	if len(in.Selector) > 0 {
		att["selector"] = in.Selector
	}
	att["template"] = flattenPodTemplateSpec(*in.Template, *userSpec.Template)
	return []interface{}{att}
}

func patchReplicationControllerSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)
	prefix += ".0."

	if d.HasChange(prefix + "replicas") {

		v := d.Get(prefix + "replicas").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/replicas",
			Value: v,
		})
	}

	if d.HasChange(prefix + "min_ready_seconds") {

		v := d.Get(prefix + "min_ready_seconds").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/minReadySeconds",
			Value: v,
		})
	}

	if d.HasChange(prefix + "template") {
		ops = append(ops, patchTemplate(
			pathPrefix+"/template",
			prefix+"template.0.",
			d,
		)...)
	}

	return ops, nil
}
