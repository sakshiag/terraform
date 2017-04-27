package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudClusterBindService_basic(t *testing.T) {
	serviceName := "testService"
	serviceKey := "testKey"
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCloudInfraClusterBindService_basic(clusterName, serviceName, serviceKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_cs_cluster_service_bind.bind_service", "namespace_id", "default"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraClusterBindService_basic(clusterName, serviceName, serviceKey string) string {
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
}


resource "ibmcloud_cf_service_instance" "service" {
  name       = "%s"
  space_guid = "${data.ibmcloud_cf_space.space.id}"
  service    = "cloudantNoSQLDB"
  plan       = "Lite"
  tags       = ["cluster-service", "cluster-bind"]
}

resource "ibmcloud_cf_service_key" "serviceKey" {
	name = "%s"
	service_instance_guid = "${ibmcloud_cf_service_instance.service.id}"
}

resource "ibmcloud_cs_cluster_service_bind" "bind_service" {
  cluster_name_id          = "${ibmcloud_cs_cluster.testacc_cluster.name}"
  service_instance_space_guid              = "${data.ibmcloud_cf_space.space.id}"
  service_instance_name_id = "${ibmcloud_cf_service_instance.service.id}"
  namespace_id 			   = "default"
  org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"
}
	`, cfOrganization, cfOrganization, cfSpace, clusterName, serviceName, serviceKey)
}
