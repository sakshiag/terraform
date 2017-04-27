package ibmcloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/services"
)

func TestAccIBMCloudInfraFileStorage_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraFileStorageConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					testAccCheckIBMCloudInfraFileStorageExists("ibmcloud_infra_file_storage.fs_endurance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_endurance", "type", "Endurance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_endurance", "capacity", "20"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_endurance", "iops", "0.25"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_endurance", "snapshot_capacity", "10"),
					testAccCheckIBMCloudInfraResources("ibmcloud_infra_file_storage.fs_endurance", "datacenter",
						"ibmcloud_infra_virtual_guest.storagevm1", "datacenter"),
					// Performance Storage
					testAccCheckIBMCloudInfraFileStorageExists("ibmcloud_infra_file_storage.fs_performance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_performance", "type", "Performance"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_performance", "capacity", "20"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_file_storage.fs_performance", "iops", "200"),
					testAccCheckIBMCloudInfraResources("ibmcloud_infra_file_storage.fs_performance", "datacenter",
						"ibmcloud_infra_virtual_guest.storagevm1", "datacenter"),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraFileStorageConfig_update,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "allowed_virtual_guest_ids.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "allowed_subnets.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "allowed_ip_addresses.#", "1"),
					// Performance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_performance", "allowed_virtual_guest_ids.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_performance", "allowed_subnets.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_performance", "allowed_ip_addresses.#", "1"),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraFileStorageConfig_enablesnapshot,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.#", "3"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.scheduleType", "WEEKLY"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.retentionCount", "5"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.minute", "2"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.hour", "13"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.dayOfWeek", "SUNDAY"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.1.scheduleType", "HOURLY"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.1.retentionCount", "5"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.1.minute", "30"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.scheduleType", "DAILY"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.retentionCount", "6"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.minute", "2"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.hour", "15"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraFileStorageConfig_updatesnapshot,
				Check: resource.ComposeTestCheckFunc(
					// Endurance Storage
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.#", "3"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.retentionCount", "2"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.minute", "2"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.hour", "13"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.0.dayOfWeek", "MONDAY"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.1.retentionCount", "3"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.1.minute", "40"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.retentionCount", "5"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.minute", "2"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.hour", "15"),
					resource.TestCheckResourceAttr("ibmcloud_infra_file_storage.fs_endurance", "snapshot_schedule.2.enable", "false"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraFileStorageExists(n string) resource.TestCheckFunc {
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

const testAccCheckIBMCloudInfraFileStorageConfig_basic = `
resource "ibmcloud_infra_virtual_guest" "storagevm1" {
    hostname = "storagevm1"
    domain = "terraformuat.ibm.com"
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

resource "ibmcloud_infra_file_storage" "fs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 0.25
        snapshot_capacity = 10
}

resource "ibmcloud_infra_file_storage" "fs_performance" {
        type = "Performance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 200
}
`
const testAccCheckIBMCloudInfraFileStorageConfig_update = `
resource "ibmcloud_infra_virtual_guest" "storagevm1" {
    hostname = "storagevm1"
    domain = "terraformuat.ibm.com"
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

resource "ibmcloud_infra_file_storage" "fs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 0.25
        allowed_virtual_guest_ids = [ "${ibmcloud_infra_virtual_guest.storagevm1.id}" ]
        allowed_subnets = [ "${ibmcloud_infra_virtual_guest.storagevm1.private_subnet}" ]
        allowed_ip_addresses = [ "${ibmcloud_infra_virtual_guest.storagevm1.ipv4_address_private}" ]
        snapshot_capacity = 10
}

resource "ibmcloud_infra_file_storage" "fs_performance" {
        type = "Performance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 100
        allowed_virtual_guest_ids = [ "${ibmcloud_infra_virtual_guest.storagevm1.id}" ]
        allowed_subnets = [ "${ibmcloud_infra_virtual_guest.storagevm1.private_subnet}" ]
        allowed_ip_addresses = [ "${ibmcloud_infra_virtual_guest.storagevm1.ipv4_address_private}" ]
}
`

const testAccCheckIBMCloudInfraFileStorageConfig_enablesnapshot = `
resource "ibmcloud_infra_virtual_guest" "storagevm1" {
    hostname = "storagevm1"
    domain = "terraformuat.ibm.com"
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

resource "ibmcloud_infra_file_storage" "fs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 0.25
        snapshot_capacity = 10
        snapshot_schedule = [
  		{
			scheduleType="WEEKLY",
			retentionCount= 5,
			minute= 2,
			hour= 13,
			dayOfWeek= "SUNDAY",
			enable= true
		},
		{
			scheduleType="HOURLY",
			retentionCount= 5,
			minute= 30,
			enable= true
		},
		
		{
			scheduleType="DAILY",
			retentionCount= 6,
			minute= 2,
			hour= 15
			enable= true
		},
 		]
}
`
const testAccCheckIBMCloudInfraFileStorageConfig_updatesnapshot = `
resource "ibmcloud_infra_virtual_guest" "storagevm1" {
    hostname = "storagevm1"
    domain = "terraformuat.ibm.com"
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

resource "ibmcloud_infra_file_storage" "fs_endurance" {
        type = "Endurance"
        datacenter = "${ibmcloud_infra_virtual_guest.storagevm1.datacenter}"
        capacity = 20
        iops = 0.25
        snapshot_capacity = 10
        snapshot_schedule = [
  		{
			scheduleType="WEEKLY",
			retentionCount= 2,
			minute= 2,
			hour= 13,
			dayOfWeek= "MONDAY",
			enable= true
		},
		{
			scheduleType="HOURLY",
			retentionCount= 3,
			minute= 40,
			enable= true
		},
		
		{
			scheduleType="DAILY",
			retentionCount= 5,
			minute= 2,
			hour= 15
			enable= false
		},
 		]
}
`
