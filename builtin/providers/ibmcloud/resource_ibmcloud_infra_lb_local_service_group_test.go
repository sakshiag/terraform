package ibmcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccIBMCloudInfraIBMCloudInfraLbLocalServiceGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraIBMCloudInfraLbLocalServiceGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group1", "port", "82"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group1", "routing_method", "CONSISTENT_HASH_IP"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group1", "routing_type", "HTTP"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group1", "allocation", "50"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group2", "port", "83"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group2", "routing_method", "CONSISTENT_HASH_IP"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group2", "routing_type", "TCP"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_lb_local_service_group.test_service_group2", "allocation", "50"),
				),
			},
		},
	})
}

const testAccCheckIBMCloudInfraIBMCloudInfraLbLocalServiceGroupConfig_basic = `
resource "ibmcloud_infra_lb_local" "testacc_foobar_lb" {
    connections = 250
    datacenter    = "tor01"
    ha_enabled  = false
}

resource "ibmcloud_infra_lb_local_service_group" "test_service_group1" {
    port = 82
    routing_method = "CONSISTENT_HASH_IP"
    routing_type = "HTTP"
    load_balancer_id = "${ibmcloud_infra_lb_local.testacc_foobar_lb.id}"
    allocation = 50
}

resource "ibmcloud_infra_lb_local_service_group" "test_service_group2" {
    port = 83
    routing_method = "CONSISTENT_HASH_IP"
    routing_type = "TCP"
    load_balancer_id = "${ibmcloud_infra_lb_local.testacc_foobar_lb.id}"
    allocation = 50
}
`
