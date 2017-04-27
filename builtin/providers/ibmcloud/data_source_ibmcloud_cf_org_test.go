package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFOrgDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFOrgDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_org.testacc_ds_org", "org", cfOrganization),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFOrgDataSourceConfig() string {
	return fmt.Sprintf(`
	
data "ibmcloud_cf_org" "testacc_ds_org" {
    org = "%s"
}`, cfOrganization)

}
