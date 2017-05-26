package ibmcloud

import (
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func validateServiceTags(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 2048 {
		errors = append(errors, fmt.Errorf(
			"%q must contain tags whose maximum length is 2048 characters", k))
	}
	return
}

func validateAllowedStringValue(validValues []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		input := v.(string)
		existed := false
		for _, s := range validValues {
			if s == input {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a value from %#v, got %q",
				k, validValues, input))
		}
		return

	}
}

func validateRoutePath(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if (len(value) < 2) || (len(value) > 128) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must contain from 2 to 128 characters ", k, value))
	}
	if !(strings.HasPrefix(value, "/")) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must start with a forward slash '/'", k, value))

	}
	if strings.Contains(value, "?") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must not contain a '?'", k, value))
	}

	return
}

func validateRoutePort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value < 1024) || (value > 65535) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must be in the range of 1024 to 65535", k, value))
	}
	return
}

func validateAppQuota(v interface{}, k string) (ws []string, errors []error) {
	memoryInMB := float64(v.(int))

	// Validate memory to match gigs format
	remaining := math.Mod(memoryInMB, 1024)
	if remaining > 0 {
		suggested := math.Ceil(memoryInMB/1024) * 1024
		errors = append(errors, fmt.Errorf(
			"Invalid 'memory' value %d megabytes, must be a multiple of 1024 (e.g. use %d)", int(memoryInMB), int(suggested)))
	}

	return
}

func validateAppInstance(v interface{}, k string) (ws []string, errors []error) {
	instances := v.(int)
	if instances < 0 {
		errors = append(errors, fmt.Errorf(
			"%q (%q) must be greater than 0", k, instances))
	}
	return

}
