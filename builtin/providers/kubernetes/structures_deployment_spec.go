package kubernetes

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/util/intstr"
)

// Flatteners
func flattenDeploymentSpec(in v1beta1.DeploymentSpec, userSpec v1beta1.DeploymentSpec) []interface{} {

	att := make(map[string]interface{})

	if in.Selector != nil {
		att["selector"] = flattenLabelSelector(in.Selector)
	}

	att["min_ready_seconds"] = in.MinReadySeconds

	if in.Replicas != nil {
		att["replicas"] = *in.Replicas
	}

	att["template"] = flattenPodTemplateSpec(in.Template, userSpec.Template)

	att["strategy"] = flattenStrategy(in.Strategy)

	if in.RevisionHistoryLimit != nil {

		att["revision_history_limit"] = *in.RevisionHistoryLimit

	}

	if in.Paused {

		att["pause"] = in.Paused
	}

	if in.ProgressDeadlineSeconds != nil {

		att["progress_deadline_seconds"] = *in.ProgressDeadlineSeconds
	}

	if in.RollbackTo != nil {

		att["rollback_to"] = flattenRollbackTo(in.RollbackTo)
	}

	return []interface{}{att}
}

func flattenRollbackTo(l *v1beta1.RollbackConfig) []interface{} {

	att := make(map[string]interface{})
	att["revision"] = l.Revision

	return []interface{}{att}
}

func flattenStrategy(l v1beta1.DeploymentStrategy) []interface{} {

	att := make(map[string]interface{})
	att["type"] = l.Type
	if l.RollingUpdate != nil {

		att["rolling_update"] = flattenRollingUpdate(l.RollingUpdate)
	}

	return []interface{}{att}
}

func flattenRollingUpdate(l *v1beta1.RollingUpdateDeployment) []interface{} {

	att := make(map[string]interface{})

	if l.MaxUnavailable != nil {
		att["max_unavailable"] = l.MaxUnavailable
	}

	if l.MaxSurge != nil {

		att["max_surge"] = l.MaxSurge
	}

	return []interface{}{att}
}

//expanders
func expandDeploymentSpec(d []interface{}) (v1beta1.DeploymentSpec, error) {
	if len(d) == 0 || d[0] == nil {
		return v1beta1.DeploymentSpec{}, nil
	}
	in := d[0].(map[string]interface{})
	obj := v1beta1.DeploymentSpec{}

	if v, ok := in["replicas"]; ok {
		obj.Replicas = ptrToInt32(int32(v.(int)))
	}

	if v, ok := in["min_ready_seconds"]; ok {
		obj.MinReadySeconds = int32(v.(int))
	}

	if v, ok := in["selector"].([]interface{}); ok && len(v) > 0 {

		obj.Selector = expandLabelSelector(v)

	}

	if v, ok := in["strategy"].([]interface{}); ok && len(v) > 0 {
		var err error
		obj.Strategy, err = expandStrategy(v)
		if err != nil {
			return obj, err
		}
	}

	if v, ok := in["revision_history_limit"]; ok {
		obj.RevisionHistoryLimit = ptrToInt32(int32(v.(int)))
	}

	if v, ok := in["pause"]; ok {
		obj.Paused = v.(bool)
	}

	if v, ok := in["progress_deadline_seconds"].(int); ok && v > 0 {
		obj.ProgressDeadlineSeconds = ptrToInt32(int32(v))
	}

	if v, ok := in["rollback_to"].([]interface{}); ok && len(v) > 0 {
		var err error
		obj.RollbackTo, err = expandRollbackTo(v)
		if err != nil {
			return obj, err
		}
	}

	if v, ok := in["template"].([]interface{}); ok {
		pd, err := expandPodTemplateSpec(v)
		if err != nil {
			return obj, err
		}
		obj.Template = pd
	}
	return obj, nil
}

func expandStrategy(d []interface{}) (v1beta1.DeploymentStrategy, error) {

	if len(d) == 0 || d[0] == nil {
		return v1beta1.DeploymentStrategy{}, nil
	}
	in := d[0].(map[string]interface{})
	obj := v1beta1.DeploymentStrategy{}
	if v, ok := in["type"]; ok {
		obj.Type = v.(v1beta1.DeploymentStrategyType)
	}
	if v, ok := in["rolling_update"].([]interface{}); ok && len(v) > 0 {
		var err error
		obj.RollingUpdate, err = expandRollingUpdate(v)
		if err != nil {
			return obj, err
		}
	}
	return obj, nil
}

func expandRollingUpdate(d []interface{}) (*v1beta1.RollingUpdateDeployment, error) {
	if len(d) == 0 || d[0] == nil {
		return &v1beta1.RollingUpdateDeployment{}, nil
	}
	in := d[0].(map[string]interface{})
	obj := &v1beta1.RollingUpdateDeployment{}
	if v, ok := in["max_surge"].(string); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			obj.MaxSurge = &intstr.IntOrString{
				Type:   intstr.String,
				StrVal: v,
			}
		} else {
			obj.MaxSurge = &intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(i),
			}
		}

	}
	if v, ok := in["max_unavailable"].(string); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			obj.MaxUnavailable = &intstr.IntOrString{
				Type:   intstr.String,
				StrVal: v,
			}
		} else {
			obj.MaxUnavailable = &intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(i),
			}
		}
	}

	return obj, nil
}

func expandRollbackTo(d []interface{}) (*v1beta1.RollbackConfig, error) {
	if len(d) == 0 || d[0] == nil {
		return &v1beta1.RollbackConfig{}, nil
	}
	in := d[0].(map[string]interface{})
	obj := &v1beta1.RollbackConfig{}
	if v, ok := in["revision"]; ok {
		obj.Revision = int64(v.(int))
	}
	return obj, nil
}

func patchDeploymentSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {

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

	if d.HasChange(prefix + "revision_history_limit") {

		v := d.Get(prefix + "revision_history_limit").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/revisionHistoryLimit",
			Value: v,
		})
	}

	if d.HasChange(prefix + "pause") {

		v := d.Get(prefix + "pause").(bool)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/paused",
			Value: v,
		})
	}

	if d.HasChange(prefix + "progress_deadline_seconds") {

		v := d.Get(prefix + "progress_deadline_seconds").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/progressDeadlineSeconds",
			Value: v,
		})
	}

	if d.HasChange(prefix + "strategy") {
		ops = append(ops, patchStrategy(
			pathPrefix+"/strategy",
			prefix+"strategy.0.",
			d,
		)...)
	}

	if d.HasChange(prefix + "rollback_to") {
		ops = append(ops, patchRollbackTo(
			pathPrefix+"/rollbackTo",
			prefix+"rollback_to.0.",
			d,
		)...)
	}

	if d.HasChange(prefix + "template") {
		ops = append(ops, patchTemplate(
			pathPrefix+"/template",
			prefix+"template.0.",
			d,
		)...)
	}

	if d.HasChange(prefix + "selector") {
		ops = append(ops, patchSelector(
			pathPrefix+"/selector",
			prefix+"selector.0.",
			d,
		)...)
	}

	return ops, nil
}

func patchSelector(pathPrefix, prefix string, d *schema.ResourceData) []PatchOperation {

	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "match_labels") {

		v := d.Get(prefix + "match_labels")
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/matchLabels",
			Value: v,
		})
	}

	if d.HasChange(prefix + "match_expressions") {
		ops = append(ops, patchMatchExpression(
			pathPrefix+"/matchExpressions",
			prefix+"match_expressions.0.",
			d,
		)...)
	}

	return ops
}

func patchMatchExpression(pathPrefix, prefix string, d *schema.ResourceData) []PatchOperation {

	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "key") {

		v := d.Get(prefix + "key")
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/key",
			Value: v,
		})
	}

	if d.HasChange(prefix + "operator") {

		v := d.Get(prefix + "operator")
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/operator",
			Value: v,
		})
	}

	if d.HasChange(prefix + "values") {

		v := d.Get(prefix + "values")
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/values",
			Value: v,
		})
	}

	return ops
}

func patchRollbackTo(pathPrefix, prefix string, d *schema.ResourceData) []PatchOperation {

	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "revision") {

		v := d.Get(prefix + "revision").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/revision",
			Value: v,
		})
	}

	return ops
}

func patchStrategy(pathPrefix, prefix string, d *schema.ResourceData) []PatchOperation {

	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "type") {

		v := d.Get(prefix + "type").(string)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/type",
			Value: v,
		})
	}

	if d.HasChange(prefix + "rolling_update") {
		ops = append(ops, patchRollingUpdate(
			pathPrefix+"/rollingUpdate",
			prefix+"rolling_update.0.",
			d,
		)...)
	}

	return ops
}

func patchRollingUpdate(pathPrefix, prefix string, d *schema.ResourceData) []PatchOperation {

	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "max_unavailable") {

		v := d.Get(prefix + "max_unavailable").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/maxUnavailable",
			Value: v,
		})
	}

	if d.HasChange(prefix + "max_surge") {

		v := d.Get(prefix + "max_surge").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/maxSurge",
			Value: v,
		})
	}

	return ops
}
