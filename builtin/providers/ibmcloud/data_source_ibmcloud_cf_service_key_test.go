package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFServiceKeyDataSource_basic(t *testing.T) {
	serviceName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	serviceKey := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceKeyDataSourceConfig(serviceName, serviceKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cf_service_key.testacc_ds_service_key", "name", serviceKey),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFServiceKeyDataSourceConfig(serviceName, serviceKey string) string {
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

		resource "ibmcloud_cf_service_key" "servicekey" {
			name = "%s"
			service_instance_guid = "${ibmcloud_cf_service_instance.service.id}"
		}
		
		data "ibmcloud_cf_service_key" "testacc_ds_service_key" {
			name = "${ibmcloud_cf_service_key.servicekey.name}"
			service_instance_name = "${ibmcloud_cf_service_instance.service.name}"
}`, cfOrganization, cfSpace, serviceName, serviceKey)

}
