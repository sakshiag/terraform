package ibmcloud

import (
	"fmt"
)

func validateServiceTags(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 2048 {
		errors = append(errors, fmt.Errorf(
			"%q must contain tags whose maximum length is 2048 characters", k))
	}
	return
}
