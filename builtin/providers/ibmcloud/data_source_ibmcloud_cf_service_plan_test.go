package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFServicePlanDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServicePlanDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_service_plan.testacc_ds_service_plan", "service", "cleardb"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_service_plan.testacc_ds_service_plan", "plan", "spark"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFServicePlanDataSourceConfig() string {
	return fmt.Sprintf(`
	
data "ibmcloud_cf_service_plan" "testacc_ds_service_plan" {
    service = "cleardb"
	plan = "spark"
}`,
	)

}
