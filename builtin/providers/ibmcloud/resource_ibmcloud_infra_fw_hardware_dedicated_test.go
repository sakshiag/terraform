package ibmcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccIBMCloudInfraFwHardwareDedicated_Basic(t *testing.T) {
	hostname := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraFwHardwareDedicated_basic(hostname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_fw_hardware_dedicated.accfw", "ha_enabled", "false"),
					testAccCheckIBMCloudInfraResources("ibmcloud_infra_fw_hardware_dedicated.accfw", "public_vlan_id",
						"ibmcloud_infra_virtual_guest.fwvm1", "public_vlan_id"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraFwHardwareDedicated_basic(hostname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_virtual_guest" "fwvm1" {
    hostname = "%s"
    domain = "terraformuat.ibm.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "sjc01"
    network_speed = 10
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "ibmcloud_infra_fw_hardware_dedicated" "accfw" {
  ha_enabled = false
  public_vlan_id = "${ibmcloud_infra_virtual_guest.fwvm1.public_vlan_id}"
}`, hostname)
}
