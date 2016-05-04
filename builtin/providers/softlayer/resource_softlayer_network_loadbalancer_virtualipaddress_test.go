package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

func TestAccSoftLayerVirtualIpAddress_Basic(t *testing.T) {
	var vip datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSoftLayerVirtualIpAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualIpAddressConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerVirtualIpAddressExists("softlayer_vip.testacc_vip", &vip),
					testAccCheckSoftLayerVirtualIpAddressAttributes(&vip),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "connection_limit", "2"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "load_balancing_method", "lc"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "notes", "test_notes"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "name", "test_load_balancer_vip"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "source_port", "80"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "type", "HTTP"),
					resource.TestCheckResourceAttr(
						"softlayer_network_loadbalancer_virtualipaddress.testacc_vip", "virtual_ip_address", "123.123.123.123"),
				),
			},
		},
	})
}

func testAccCheckSoftLayerVirtualIpAddressDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).networkApplicationDeliveryControllerService

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "softlayer_network_loadbalancer_virtualipaddress" {
			continue
		}

		nadcId, _ := strconv.Atoi(rs.Primary.Meta["nad_controller_id"].(string))
		vipName, _ := rs.Primary.Meta["name"].(string)

		// Try to find the vip
		_, err := client.GetVirtualIpAddress(nadcId, vipName)

		if err == nil {
			fmt.Errorf("Virtual ip address still exists")
		}
	}

	return nil
}

func testAccCheckSoftLayerVirtualIpAddressAttributes(vip *datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vip.Id <= 0 {
			return fmt.Errorf("Bad vip id: %d", vip.Id)
		}

		return nil
	}
}

func testAccCheckSoftLayerVirtualIpAddressExists(n string, vip *datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Meta["nad_controller_id"] == "" {
			return fmt.Errorf("No nadc ID is set")
		}

		if rs.Primary.Meta["name"] == "" {
			return fmt.Errorf("VIP name is not set")
		}

		client := testAccProvider.Meta().(*Client).networkApplicationDeliveryControllerService
		nadcId, _ := strconv.Atoi(rs.Primary.Meta["nad_controller_id"].(string))
		vipName, _ := rs.Primary.Meta["name"].(string)

		foundVip, err := client.GetVirtualIpAddress(nadcId, vipName)

		if err != nil {
			return err
		}

		if strconv.Itoa(int(foundVip.Id)) != rs.Primary.ID {
			return fmt.Errorf("Vip not found")
		}

		*vip = foundVip

		return nil
	}
}

var testAccCheckSoftLayerVirtualIpAddressConfig_basic = `
resource "softlayer_network_loadbalancer_virtualipaddress" "testacc_vip" {
    name = "test_load_balancer_vip"
    nad_controller_id = 15293
    connection_limit = 2
    load_balancing_method = "lc"
    notes = "test_notes"
    source_port = 80
    type = "HTTP"
    virtual_ip_address = "123.123.123.123"
}`
