/*
* Licensed Materials - Property of IBM
* (C) Copyright IBM Corp. 2017. All Rights Reserved.
* US Government Users Restricted Rights - Use, duplication or
* disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 */

package ibmcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccIBMCloudInfraVlan_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraVlanConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "name", "test_vlan"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "datacenter", "lon02"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "type", "PUBLIC"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "softlayer_managed", "false"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "router_hostname", "fcr01a.lon02"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "subnet_size", "8"),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraVlanConfig_name_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_vlan.test_vlan", "name", "test_vlan_update"),
				),
			},
		},
	})
}

const testAccCheckIBMCloudInfraVlanConfig_basic = `
resource "ibmcloud_infra_vlan" "test_vlan" {
   name = "test_vlan"
   datacenter = "lon02"
   type = "PUBLIC"
   subnet_size = 8
   router_hostname = "fcr01a.lon02"
}`

const testAccCheckIBMCloudInfraVlanConfig_name_update = `
resource "ibmcloud_infra_vlan" "test_vlan" {
   name = "test_vlan_update"
   datacenter = "lon02"
   type = "PUBLIC"
   subnet_size = 8
   router_hostname = "fcr01a.lon02"
}`
