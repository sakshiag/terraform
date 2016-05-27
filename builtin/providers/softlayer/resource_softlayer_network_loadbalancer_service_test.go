package softlayer

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSoftLayerLoadBalancerService_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerLoadBalancerServiceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_service.testacc_service", "name", "test_load_balancer_service"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_service.testacc_service", "destination_ip_address", "192.155.238.31"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_service.testacc_service", "destination_port", "8000"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_service.testacc_service", "weight", "55"),
				),
			},
		},
	})
}

var testAccCheckSoftLayerLoadBalancerServiceConfig_basic = `
resource "softlayer_network_application_delivery_controller" "testacc_foobar_nadc" {
    type = "Netscaler VPX"
    datacenter = "DALLAS06"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
}

resource "softlayer_network_loadbalancer_virtualipaddress" "testacc_vip" {
    name = "test_load_balancer_vip"
    nad_controller_id = "${softlayer_network_application_delivery_controller.testacc_foobar_nadc.id}"
    load_balancing_method = "lc"
    source_port = 80
    type = "HTTP"
    virtual_ip_address = "23.246.204.65"
}

resource "softlayer_network_loadbalancer_service" "testacc_service" {
  name = "test_load_balancer_service"
  vip_id = "${softlayer_network_loadbalancer_virtualipaddress.testacc_vip.id}"
  destination_ip_address = "192.155.238.31"
  destination_port = 8000
  weight = 55
}
`
