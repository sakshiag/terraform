package ibmcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var cfOrganization string
var cfSpace string

func init() {
	cfOrganization = os.Getenv("IBMCLOUD_ORG")
	if cfOrganization == "" {
		fmt.Println("[WARN] Set the environment variable IBMCLOUD_ORG for testing ibmcloud_cf_space  resource Some tests for that resource will fail if this is not set correctly")
	}
	cfSpace = os.Getenv("IBMCLOUD_SPACE")
	if cfSpace == "" {
		fmt.Println("[WARN] Set the environment variable IBMCLOUD_SPACE for testing ibmcloud_cf_space  resource Some tests for that resource will fail if this is not set correctly")
	}
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ibmcloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {

	requiredEnv := map[string]string{
		"ibmid":    "IBMID",
		"password": "IBMID_PASSWORD",
	}

	for _, param := range []string{"ibmid", "ibmid_password"} {
		value, _ := testAccProvider.Schema[param].DefaultFunc()
		if value == "" {
			t.Fatalf("%s must be set for acceptance test", requiredEnv[param])
		}
	}
}
