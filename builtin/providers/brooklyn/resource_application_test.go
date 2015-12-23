package brooklyn

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
	"os"
	"fmt"
)

const (
	TEST_YAML_FILE = "BROOKLYN_TEST_YAML_FILE"
)

func TestAccArtifact_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testApplicationCreationPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccApplication_basic, os.Getenv(TEST_YAML_FILE)),
				// Check only that application is running, as other fields depend on the YAML file itself
				Check: resource.TestCheckResourceAttr("brooklyn_application.application_test", "status", "RUNNING"),
			},
		},
	})
}

func testApplicationCreationPreCheck(t *testing.T) {
	// Check basic env properties
	testAccPreCheck(t)
	// Check additional env properties for Application creation
	if v := os.Getenv(TEST_YAML_FILE); v == "" {
		t.Fatal("BROOKLYN_TEST_YAML_FILE must be set for application creation acceptance tests")
	}
}

const testAccApplication_basic = `
resource "brooklyn_application" "application_test" {
	application_spec = "%s"
}`

