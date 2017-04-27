package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFServiceInstanceDataSource_basic(t *testing.T) {
	serviceName := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceInstanceDataSourceConfig(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_service_instance.testacc_ds_service_instance", "name", serviceName),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFServiceInstanceDataSourceConfig(serviceName string) string {
	return fmt.Sprintf(`
	data "ibmcloud_cf_space" "spacedata" {
			org    = "%s"
			space  = "%s"
		}
		
		resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service","cluster-bind"]
		}

	
		data "ibmcloud_cf_service_instance" "testacc_ds_service_instance" {
			name = "${ibmcloud_cf_service_instance.service.name}"
}`, cfOrganization, cfSpace, serviceName)

}
