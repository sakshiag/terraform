package ibmcloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
)

func TestAccIBMCloudInfraLbVpx_Basic(t *testing.T) {
	var nadc datatypes.Network_Application_Delivery_Controller

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraLbVpxConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraLbVpxExists("ibmcloud_infra_lb_vpx.testacc_foobar_vpx", &nadc),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "type", "NetScaler VPX"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "datacenter", "dal09"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "speed", "10"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "plan", "Standard"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "ip_count", "2"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "version", "10.1"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_vpx.testacc_foobar_vpx", "vip_pool.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraLbVpxExists(n string, nadc *datatypes.Network_Application_Delivery_Controller) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		nadcId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetNetworkApplicationDeliveryControllerService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		found, err := service.Id(nadcId).GetObject()
		if err != nil {
			return err
		}

		if strconv.Itoa(int(*found.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*nadc = found

		return nil
	}
}

const testAccCheckIBMCloudInfraLbVpxConfig_basic = `
resource "ibmcloud_infra_lb_vpx" "testacc_foobar_vpx" {
    datacenter = "dal09"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
}`
