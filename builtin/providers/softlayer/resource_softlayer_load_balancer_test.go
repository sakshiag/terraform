package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSoftLayerLoadBalancer_Basic(t *testing.T) {
	var lb datatypes.SoftLayer_Load_Balancer

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerLoadBalancerExists("softlayer_load_balancer.testacc_foobar_nadc", &lb),
				),
			},
		},
	})
}

func testAccCheckSoftLayerLoadBalancerExists(n string, lb *datatypes.SoftLayer_Load_Balancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		lbId, _ := strconv.Atoi(rs.Primary.ID)

		client := testAccProvider.Meta().(*Client).loadBalancerService
		found, err := client.GetObject(lbId)

		if err != nil {
			return err
		}

		if strconv.Itoa(int(found.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*lb = found

		return nil
	}
}

const testAccCheckSoftLayerLoadBalancerConfig_basic = `
resource "softlayer_load_balancer" "testacc_foobar_lb" {
    connections = 15000
    location    = "tok02"
    ha_enabled  = false
}`
