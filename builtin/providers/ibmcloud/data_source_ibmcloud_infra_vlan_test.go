package ibmcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudInfraVlanDataSource_Basic(t *testing.T) {

	name := fmt.Sprintf("terraformuat_vlan_%s", acctest.RandString(2))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCloudInfraVlanDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraResources("data.ibmcloud_infra_vlan.tfacc_vlan", "number",
						"ibmcloud_infra_vlan.test_vlan_private", "vlan_number"),
					//resource.TestCheckResourceAttr("data.ibmcloud_infra_vlan.tfacc_vlan", "number", number),
					resource.TestCheckResourceAttr("data.ibmcloud_infra_vlan.tfacc_vlan", "name", name),
					resource.TestMatchResourceAttr("data.ibmcloud_infra_vlan.tfacc_vlan", "id", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraVlanDataSourceConfig(name string) string {
	return fmt.Sprintf(`
    resource "ibmcloud_infra_vlan" "test_vlan_private" {
    name            = "%s"
    datacenter      = "dal06"
    type            = "PRIVATE"
    subnet_size     = 8
    
}
data "ibmcloud_infra_vlan" "tfacc_vlan" {
    number = "${ibmcloud_infra_vlan.test_vlan_private.vlan_number}"
    name = "${ibmcloud_infra_vlan.test_vlan_private.name}"
}`, name)
}
