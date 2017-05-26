package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFSharedDomainDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFSharedDomainDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibmcloud_cf_shared_domain.testacc_domain", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFSharedDomainDataSourceConfig() string {
	return fmt.Sprintf(`
	
		data "ibmcloud_cf_shared_domain" "testacc_domain" {
			name = "mybluemix.net"
		}`)

}
