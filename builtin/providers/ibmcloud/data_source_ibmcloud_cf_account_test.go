package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFAccountDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAccountDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_org.testacc_ds_org", "org", cfOrganization),
					resource.TestCheckResourceAttrSet(
						"data.ibmcloud_cf_account.testacc_acc", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFAccountDataSourceConfig() string {
	return fmt.Sprintf(`
	
data "ibmcloud_cf_org" "testacc_ds_org" {
    org = "%s"
}

data "ibmcloud_cf_account" "testacc_acc" {
    org_guid = "${data.ibmcloud_cf_org.testacc_ds_org.id}"
}`, cfOrganization)

}
