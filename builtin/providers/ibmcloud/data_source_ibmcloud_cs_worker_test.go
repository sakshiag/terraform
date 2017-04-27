package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCSWorkerDataSource_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCSWorkerDataSourceConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_cs_worker.testacc_ds_worker", "state", "normal"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCSWorkerDataSourceConfig(clusterName string) string {
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
    name = "%s"
    datacenter = "dal10"
    workers = [{
    name = "worker1"
    action = "add"
  },]
	machine_type = "free"
	isolation = "public"
	public_vlan_id = "vlan"
	private_vlan_id = "vlan"

    org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"
}
data "ibmcloud_cs_cluster" "testacc_ds_cluster" {
	org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"
    cluster_name_id = "${ibmcloud_cs_cluster.testacc_cluster.id}"
}
data "ibmcloud_cs_worker" "testacc_ds_worker" {
	org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"
    worker_id = "${data.ibmcloud_cs_cluster.testacc_ds_cluster.workers[0]}"
}
`, cfOrganization, cfOrganization, cfSpace, clusterName)
}
