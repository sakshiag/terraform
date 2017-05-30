package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFAppDataSource_Basic(t *testing.T) {
	var conf cfv2.AppFields
	appName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	routeHostName := fmt.Sprintf("terraform-route-host-%d", acctest.RandInt())
	svcName := fmt.Sprintf("tfsvc-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFAppDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppDataSourceBasic(routeHostName, svcName, appName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", appName),
					resource.TestCheckResourceAttrSet("data.ibmcloud_cf_app.ds", "id"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "name", appName),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "buildpack", "sdk-for-nodejs"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "environment_json.%", "1"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "state", "STARTED"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "package_state", "STAGED"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "route_guid.#", "1"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "service_instance_guid.#", "1"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "memory", "128"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "instances", "1"),
					resource.TestCheckResourceAttr("data.ibmcloud_cf_app.ds", "disk_quota", "512"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFAppDataSourceBasic(routeHost, serviceInstanceName, appName string) (config string) {
	config = fmt.Sprintf(`
data "ibmcloud_cf_space" "space" {
  org   = "%s"
  space = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
  name = "mybluemix.net"
}

resource "ibmcloud_cf_route" "route" {
  domain_guid = "${data.ibmcloud_cf_shared_domain.domain.id}"
  space_guid  = "${data.ibmcloud_cf_space.space.id}"
  host        = "%s"
}

resource "ibmcloud_cf_service_instance" "service" {
  name       = "%s"
  space_guid = "${data.ibmcloud_cf_space.space.id}"
  service    = "cleardb"
  plan       = "spark"
  tags       = ["cluster-service"]
}

resource "ibmcloud_cf_app" "app" {
  name                  = "%s"
  space_guid            = "${data.ibmcloud_cf_space.space.id}"
  app_path              = "test-fixtures/app1.zip"
  wait_time_minutes     = 20
  buildpack             = "sdk-for-nodejs"
  instances             = 1
  route_guid            = ["${ibmcloud_cf_route.route.id}"]
  service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}"]
  disk_quota            = 512
  memory                = 128
  instances             = 1
  disk_quota            = 512

  environment_json = {
    "test" = "test1"
  }
}

data  "ibmcloud_cf_app" "ds" {
  name       = "${ibmcloud_cf_app.app.name}"
  space_guid = "${data.ibmcloud_cf_space.space.id}"
}
`, cfOrganization, cfSpace, routeHost, serviceInstanceName, appName)
	return
}
