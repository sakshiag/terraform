package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCluster_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCloudInfraCluster_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_cs_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibmcloud_cs_cluster.testacc_cluster", "worker_num", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraCluster_basic(clusterName string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_org" "org" {
    org = "%s"
}

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_account" "acc" {
   org_guid = "${data.ibmcloud_cf_org.org.id}"
}

resource "ibmcloud_cs_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "dal10"

  org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"

  workers = [{
    name = "worker1"

    action = "add"
  }]

  machine_type    = "free"
  isolation       = "public"
  public_vlan_id  = "vlan"
  private_vlan_id = "vlan"
}	`, cfOrganization, cfOrganization, cfSpace, clusterName)
}
