package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/kubernetes/pkg/api/resource"
	apiValidation "k8s.io/kubernetes/pkg/api/validation"
	utilValidation "k8s.io/kubernetes/pkg/util/validation"
)

func validateAnnotations(value interface{}, key string) (ws []string, es []error) {
	m := value.(map[string]interface{})
	for k, _ := range m {
		errors := utilValidation.IsQualifiedName(strings.ToLower(k))
		if len(errors) > 0 {
			for _, e := range errors {
				es = append(es, fmt.Errorf("%s (%q) %s", key, k, e))
			}
		}
	}
	return
}

func validateName(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)

	errors := apiValidation.NameIsDNSLabel(v, false)
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

func validateGenerateName(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)

	errors := apiValidation.NameIsDNSLabel(v, true)
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

func validateLabels(value interface{}, key string) (ws []string, es []error) {
	m := value.(map[string]interface{})
	for k, v := range m {
		for _, msg := range utilValidation.IsQualifiedName(k) {
			es = append(es, fmt.Errorf("%s (%q) %s", key, k, msg))
		}
		val := v.(string)
		for _, msg := range utilValidation.IsValidLabelValue(val) {
			es = append(es, fmt.Errorf("%s (%q) %s", key, val, msg))
		}
	}
	return
}

func validateResourceList(value interface{}, key string) (ws []string, es []error) {
	m := value.(map[string]interface{})
	for k, value := range m {
		if _, ok := value.(int); ok {
			continue
		}

		if v, ok := value.(string); ok {
			_, err := resource.ParseQuantity(v)
			if err != nil {
				es = append(es, fmt.Errorf("%s.%s (%q): %s", key, k, v, err))
			}
			continue
		}

		err := "Value can be either string or int"
		es = append(es, fmt.Errorf("%s.%s (%#v): %s", key, k, value, err))
	}
	return
}

func validateActiveDeadlineSeconds(value interface{}, key string) (ws []string, es []error) {
	v := value.(int)
	if v <= 0 {
		es = append(es, fmt.Errorf("%s must be positive", key))
	}
	return
}

func validateDNSPolicy(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	if v != "ClusterFirst" && v != "Default" {
		es = append(es, fmt.Errorf("%s must be either ClusterFirst or Default", key))
	}
	return
}

func validateRestartPolicy(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	switch v {
	case "Always", "OnFailure", "Never":
		return
	default:
		es = append(es, fmt.Errorf("%s must be one of Always, OnFailure or Never ", key))
	}
	return
}

func validateTerminationGracePeriodSeconds(value interface{}, key string) (ws []string, es []error) {
	v := value.(int)
	if v < 0 {
		es = append(es, fmt.Errorf("%s must be non-negative", key))
	}
	return
}
