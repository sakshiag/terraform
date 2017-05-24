package ibmcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
)

//HashInt ...
func HashInt(v interface{}) int { return v.(int) }

//Expanders

func expandIntList(input []interface{}) []int {
	vs := make([]int, len(input))
	for i, v := range input {
		vs[i] = v.(int)
	}
	return vs
}

//Flattners

func flattenIntList(list []int) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

func flattenStorageID(in []datatypes.Network_Storage) *schema.Set {
	var out = make([]interface{}, len(in))
	for i, v := range in {
		out[i] = *v.Id
	}
	return schema.NewSet(HashInt, out)
}
