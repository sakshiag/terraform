package ibmcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccIBMCloudInfraIBMCloudInfraLbLocalShared_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraIBMCloudInfraLbLocalConfigShared_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "connections", "250"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "datacenter", "dal09"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "ha_enabled", "false"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "dedicated", "false"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "ssl_enabled", "false"),
				),
			},
		},
	})
}

func TestAccIBMCloudInfraIBMCloudInfraLbLocalDedicated_Basic(t *testing.T) {
	t.SkipNow()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraIBMCloudInfraLbLocalDedicatedConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "connections", "15000"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "datacenter", "dal09"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "ha_enabled", "false"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "dedicated", "true"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local.testacc_foobar_lb", "ssl_enabled", "true"),
				),
			},
		},
	})
}

const testAccCheckIBMCloudInfraIBMCloudInfraLbLocalConfigShared_basic = `
resource "ibmcloud_infra_lb_local" "testacc_foobar_lb" {
    connections = 250
    datacenter    = "dal09"
    ha_enabled  = false
}`

const testAccCheckIBMCloudInfraIBMCloudInfraLbLocalDedicatedConfig_basic = `
resource "ibmcloud_infra_lb_local" "testacc_foobar_lb" {
    connections = 15000
    datacenter    = "dal09"
    ha_enabled  = false
    dedicated = true	
}`
