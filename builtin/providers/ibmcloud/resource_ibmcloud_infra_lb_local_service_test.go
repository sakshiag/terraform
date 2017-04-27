package ibmcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccBluemixIBMCloudInfraLbLocalService_Basic(t *testing.T) {
	hostname := acctest.RandString(16)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckBluemixIBMCloudInfraLbLocalServiceConfig_basic(hostname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service.test_service", "port", "80"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service.test_service", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service.test_service", "weight", "1"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service.test_service", "health_check_type", "DNS"),
				),
			},
		},
	})
}

func testAccCheckBluemixIBMCloudInfraLbLocalServiceConfig_basic(hostname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_virtual_guest" "test_server_1" {
    hostname = "%s"
    domain = "terraformuat.ibm.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 10
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25, 10, 20]
    user_metadata = "{\"value\":\"newvalue\"}"
    dedicated_acct_host_only = true
    local_disk = false
}

resource "ibmcloud_infra_lb_local" "testacc_foobar_lb" {
    connections = 250
    datacenter    = "dal06"
    ha_enabled  = false
    dedicated = false
}

resource "ibmcloud_infra_lb_local_service_group" "test_service_group" {
    port = 82
    routing_method = "CONSISTENT_HASH_IP"
    routing_type = "HTTP"
    load_balancer_id = "${ibmcloud_infra_lb_local.testacc_foobar_lb.id}"
    allocation = 100
}

resource "ibmcloud_infra_lb_local_service" "test_service" {
    port = 80
    enabled = true
    service_group_id = "${ibmcloud_infra_lb_local_service_group.test_service_group.service_group_id}"
    weight = 1
    health_check_type = "DNS"
    ip_address_id = "${ibmcloud_infra_virtual_guest.test_server_1.ip_address_id}"
}`, hostname)
}
