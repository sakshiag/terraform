package softlayer

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccSoftLayerLoadBalancer_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"softlayer_load_balancer.testacc_foobar_lb", "connections", "15000"),
					resource.TestCheckResourceAttr(
						"softlayer_load_balancer.testacc_foobar_lb", "location", "tok02"),
					resource.TestCheckResourceAttr(
						"softlayer_load_balancer.testacc_foobar_lb", "ha_enabled", "false"),
				),
			},
		},
	})
}

const testAccCheckSoftLayerLoadBalancerConfig_basic = `
resource "softlayer_load_balancer" "testacc_foobar_lb" {
    connections = 15000
    location    = "tok02"
    ha_enabled  = false
}`
