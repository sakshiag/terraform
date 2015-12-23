package brooklyn

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	"os"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"brooklyn": testAccProvider,
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
	if v := os.Getenv("BROOKLYN_URL"); v == "" {
		t.Fatal("BROOKLYN_URL must be set for acceptance tests")
	}
	if v := os.Getenv("BROOKLYN_USERNAME"); v == "" {
		t.Fatal("BROOKLYN_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("BROOKLYN_PASSWORD"); v == "" {
		t.Fatal("BROOKLYN_USERNAME must be set for acceptance tests")
	}
}