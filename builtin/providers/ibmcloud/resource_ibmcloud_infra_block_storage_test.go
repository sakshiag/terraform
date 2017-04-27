package ibmcloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/services"
)

func TestAccIBMCloudInfraBlockStorage_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraBlockStorageConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					testAccCheckIBMCloudInfraBlockStorageExists("ibmcloud_infra_block_storage.bs_endurance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "type", "Endurance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "capacity", "20"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "iops", "0.25"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "snapshot_capacity", "10"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "os_format_type", "Linux"),
					testAccCheckIBMCloudInfraResources("ibmcloud_infra_block_storage.bs_endurance", "datacenter",
						"ibmcloud_infra_virtual_guest.storagevm2", "datacenter"),
					// Performance Storage
					testAccCheckIBMCloudInfraBlockStorageExists("ibmcloud_infra_block_storage.bs_performance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_performance", "type", "Performance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_performance", "capacity", "20"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_performance", "iops", "100"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_block_storage.bs_endurance", "os_format_type", "Linux"),
					testAccCheckIBMCloudInfraResources("ibmcloud_infra_block_storage.bs_performance", "datacenter",
						"ibmcloud_infra_virtual_guest.storagevm2", "datacenter"),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraBlockStorageConfig_update,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_block_storage.bs_endurance", "allowed_virtual_guest_ids.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_block_storage.bs_endurance", "allowed_ip_addresses.#", "1"),
					// Performance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_block_storage.bs_performance", "allowed_virtual_guest_ids.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_block_storage.bs_performance", "allowed_ip_addresses.#", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraBlockStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		storageId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetNetworkStorageService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		foundStorage, err := service.Id(storageId).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundStorage.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

const testAccCheckIBMCloudInfraBlockStorageConfig_basic = `
resource "ibmcloud_infra_virtual_guest" "storagevm2" {
    hostname = "storagevm2"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "ibmcloud_infra_block_storage" "bs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm2.datacenter}"
        capacity = 20
        iops = 0.25
        snapshot_capacity = 10
        os_format_type = "Linux"
}

resource "ibmcloud_infra_block_storage" "bs_performance" {
        type = "Performance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm2.datacenter}"
        capacity = 20
        iops = 100
        os_format_type = "Linux"
}
`
const testAccCheckIBMCloudInfraBlockStorageConfig_update = `
resource "ibmcloud_infra_virtual_guest" "storagevm2" {
    hostname = "storagevm2"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "ibmcloud_infra_block_storage" "bs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm2.datacenter}"
        capacity = 20
        iops = 0.25
        os_format_type = "Linux"
        allowed_virtual_guest_ids = [ "${ibmcloud_infra_virtual_guest.storagevm2.id}" ]
        allowed_ip_addresses = [ "${ibmcloud_infra_virtual_guest.storagevm2.ipv4_address_private}" ]
        snapshot_capacity = 10
}

resource "ibmcloud_infra_block_storage" "bs_performance" {
        type = "Performance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm2.datacenter}"
        capacity = 20
        iops = 100
        os_format_type = "Linux"
        allowed_virtual_guest_ids = [ "${ibmcloud_infra_virtual_guest.storagevm2.id}" ]
        allowed_ip_addresses = [ "${ibmcloud_infra_virtual_guest.storagevm2.ipv4_address_private}" ]
}
`
