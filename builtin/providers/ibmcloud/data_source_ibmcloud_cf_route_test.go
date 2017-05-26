package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFRouteDataSource_basic(t *testing.T) {
	host := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFRouteDataSourceConfig(host),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibmcloud_cf_route.testacc_route", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFRouteDataSourceConfig(host string) string {
	return fmt.Sprintf(`
		data "ibmcloud_cf_space" "spacedata" {
			org    = "%s"
			space  = "%s"
		}
		
		data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}
		
		resource "ibmcloud_cf_route" "route" {
			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			host              = "%s"
			path              = "/app"
		}
		
		data "ibmcloud_cf_route" "testacc_route" {
			domain_guid       = "${ibmcloud_cf_route.route.domain_guid}"
			space_guid        = "${ibmcloud_cf_route.route.space_guid}"
			host              = "${ibmcloud_cf_route.route.host}"
			path              = "${ibmcloud_cf_route.route.path}"
	    }
	`, cfOrganization, cfSpace, host)

}
