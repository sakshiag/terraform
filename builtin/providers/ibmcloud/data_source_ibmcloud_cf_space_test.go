package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFSpaceDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFSpaceDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_space.testacc_ds_space", "org", cfOrganization),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_space.testacc_ds_space", "space", cfSpace),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFSpaceDataSourceConfig() string {
	return fmt.Sprintf(`
	
data "ibmcloud_cf_space" "testacc_ds_space" {
    org = "%s"
	space = "%s"
}`, cfOrganization, cfSpace)

}
