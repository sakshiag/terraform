package softlayer

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	datatypes "github.com/maximilien/softlayer-go/data_types"
)

func TestAccSoftLayerVirtualserver_Basic(t *testing.T) {
	var server datatypes.SoftLayer_Virtual_Guest

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckSoftLayerVirtualserverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualserverConfig_basic,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-1", &server),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "name", "terraform-test"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "domain", "bar.example.com"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "region", "ams01"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "public_network_speed", "10"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "hourly_billing", "true"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "private_network_only", "false"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "cpu", "1"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "ram", "1024"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "disks.0", "25"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "disks.1", "10"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "disks.2", "20"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "user_data", "{\"value\":\"newvalue\"}"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "local_disk", "false"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "dedicated_acct_host_only", "true"),
					// TODO: Will be changed in future, when the following issue is implemented: https://github.com/TheWeatherCompany/softlayer-go/issues/3.
					// TODO: For now, as agreed in issue https://github.com/TheWeatherCompany/terraform/issues/5, use hardcoded values for VLANs.
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "frontend_vlan_id", "1085155"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "backend_vlan_id", "1085157"),
				),
			},

			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualserverConfig_userDataUpdate,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-1", &server),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "user_data", "updatedData"),
				),
			},

			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualserverConfig_upgradeMemoryNetworkSpeed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-1", &server),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"softlayer_virtualserver.terraform-acceptance-test-1", "public_network_speed", "100"),
				),
			},

			// TODO: currently CPU upgrade test is disabled, due to unexpected behavior of field "dedicated_acct_host_only". For some reason it is reset by SoftLayer to "false". To be aligned with Daniel and Chris how to proceed with it.
//			resource.TestStep{
//				Config: testAccCheckSoftLayerVirtualserverConfig_vmUpgradeCPUs,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-1", &server),
//					resource.TestCheckResourceAttr(
//						"softlayer_virtualserver.terraform-acceptance-test-1", "cpu", "2"),
//				),
//			},

		},
	})
}

func TestAccSoftLayerVirtualserver_BlockDeviceTemplateGroup(t *testing.T) {
	var server datatypes.SoftLayer_Virtual_Guest

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckSoftLayerVirtualserverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualserverConfig_blockDeviceTemplateGroup,
				Check: resource.ComposeTestCheckFunc(
					// block_device_template_group_gid value is hardcoded. If it's valid then virtual server will be created well
					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-BDTGroup", &server),
				),
			},
		},
	})
}

func TestAccSoftLayerVirtualserver_postInstallScriptUri(t *testing.T) {
	var server datatypes.SoftLayer_Virtual_Guest

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckSoftLayerVirtualserverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerVirtualserverConfig_postInstallScriptUri,
				Check: resource.ComposeTestCheckFunc(
					// block_device_template_group_gid value is hardcoded. If it's valid then virtual server will be created well
					testAccCheckSoftLayerVirtualserverExists("softlayer_virtualserver.terraform-acceptance-test-pISU", &server),
				),
			},
		},
	})
}

func testAccCheckSoftLayerVirtualserverDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).virtualGuestService

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "softlayer_virtualserver" {
			continue
		}

		serverId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the server
		_, err := client.GetObject(serverId)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for server (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckSoftLayerVirtualserverExists(n string, server *datatypes.SoftLayer_Virtual_Guest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No virtual server ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*Client).virtualGuestService
		retrieveServer, err := client.GetObject(id)

		if err != nil {
			return err
		}

		fmt.Printf("The ID is %d", id)

		if retrieveServer.Id != id {
			return fmt.Errorf("Virtual server not found")
		}

		*server = retrieveServer

		return nil
	}
}

const testAccCheckSoftLayerVirtualserverConfig_basic = `
resource "softlayer_virtualserver" "terraform-acceptance-test-1" {
    name = "terraform-test"
    domain = "bar.example.com"
    image = "DEBIAN_7_64"
    region = "ams01"
    public_network_speed = 10
    hourly_billing = true
	private_network_only = false
    cpu = 1
    ram = 1024
    disks = [25, 10, 20]
    user_data = "{\"value\":\"newvalue\"}"
    dedicated_acct_host_only = true
    local_disk = false
    frontend_vlan_id = 1085155
	backend_vlan_id = 1085157
}
`

const testAccCheckSoftLayerVirtualserverConfig_userDataUpdate = `
resource "softlayer_virtualserver" "terraform-acceptance-test-1" {
    name = "terraform-test"
    domain = "bar.example.com"
    image = "DEBIAN_7_64"
    region = "ams01"
    public_network_speed = 10
    hourly_billing = true
    cpu = 1
    ram = 1024
    disks = [25, 10, 20]
    user_data = "updatedData"
    dedicated_acct_host_only = true
    local_disk = false
    frontend_vlan_id = 1085155
	backend_vlan_id = 1085157
}
`

const testAccCheckSoftLayerVirtualserverConfig_upgradeMemoryNetworkSpeed = `
resource "softlayer_virtualserver" "terraform-acceptance-test-1" {
    name = "terraform-test"
    domain = "bar.example.com"
    image = "DEBIAN_7_64"
    region = "ams01"
    public_network_speed = 100
    hourly_billing = true
    cpu = 1
    ram = 2048
    disks = [25, 10, 20]
    user_data = "updatedData"
    dedicated_acct_host_only = true
    local_disk = false
    frontend_vlan_id = 1085155
	backend_vlan_id = 1085157
}
`

const testAccCheckSoftLayerVirtualserverConfig_vmUpgradeCPUs = `
resource "softlayer_virtualserver" "terraform-acceptance-test-1" {
    name = "terraform-test"
    domain = "bar.example.com"
    image = "DEBIAN_7_64"
    region = "ams01"
    public_network_speed = 100
    hourly_billing = true
    cpu = 2
    ram = 2048
    disks = [25, 10, 20]
    user_data = "updatedData"
    dedicated_acct_host_only = true
    local_disk = false
    frontend_vlan_id = 1085155
	backend_vlan_id = 1085157
}
`

const testAccCheckSoftLayerVirtualserverConfig_postInstallScriptUri = `
resource "softlayer_virtualserver" "terraform-acceptance-test-pISU" {
    name = "terraform-test-pISU"
    domain = "bar.example.com"
    image = "DEBIAN_7_64"
    region = "ams01"
    public_network_speed = 10
    hourly_billing = true
	private_network_only = false
    cpu = 1
    ram = 1024
    disks = [25, 10, 20]
    user_data = "{\"value\":\"newvalue\"}"
    dedicated_acct_host_only = true
    local_disk = false
    post_install_script_uri = "https://www.google.com"
}
`

const testAccCheckSoftLayerVirtualserverConfig_blockDeviceTemplateGroup = `
resource "softlayer_virtualserver" "terraform-acceptance-test-BDTGroup" {
    name = "terraform-test-blockDeviceTemplateGroup"
    domain = "bar.example.com"
    region = "ams01"
    public_network_speed = 10
    hourly_billing = false
    cpu = 1
    ram = 1024
    local_disk = false
    block_device_template_group_gid = "ac2b413c-9893-4178-8e62-a24cbe2864db"
}
`
