package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudInfraVirtualGuestDataSource_basic(t *testing.T) {
	hostname := acctest.RandString(16)
	domain := "ds.terraform.ibm.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraVirtualGuestDataSourceConfigBasic(hostname, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_infra_virtual_guest.tf-vg-ds-acc-test", "power_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.ibmcloud_infra_virtual_guest.tf-vg-ds-acc-test", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraVirtualGuestDataSourceConfigBasic(hostname, domain string) string {
	return fmt.Sprintf(`
	resource "ibmcloud_infra_virtual_guest" "tf-vg-acc-test" {
    hostname = "%s"
    domain = "%s"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 10
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25, 10, 20]
    tags = ["data-source-test"]
    dedicated_acct_host_only = true
    local_disk = false
}
data "ibmcloud_infra_virtual_guest" "tf-vg-ds-acc-test" {
    hostname = "${ibmcloud_infra_virtual_guest.tf-vg-acc-test.hostname}"
	domain = "${ibmcloud_infra_virtual_guest.tf-vg-acc-test.domain}"
}`, hostname, domain)
}
