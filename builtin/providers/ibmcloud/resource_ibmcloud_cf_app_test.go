package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIBMCloudCFApp_Basic(t *testing.T) {
	var conf cfv2.AppFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	serviceName := fmt.Sprintf("ter_service_%d", acctest.RandInt())
	updatedName := fmt.Sprintf("terraform_updated_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFAppDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppUpdate(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "124"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "buildpack", "nodejs_buildpack"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppUpdateRoute(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppRemoveRoute(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppBindService(updatedName, serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "service_instance_guid.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppAddDeleteService(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "service_instance_guid.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFAppDestroy(s *terraform.State) error {

	appClient := testAccProvider.Meta().(ClientSession).CloudFoundryAppClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_app" {
			continue
		}
		appGUID := rs.Primary.ID

		_, err := appClient.Get(appGUID)
		if err == nil {
			return fmt.Errorf("App still exists: %s", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckIBMCloudCFAppExists(n string, obj *cfv2.AppFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		appClient := testAccProvider.Meta().(ClientSession).CloudFoundryAppClient()
		appGUID := rs.Primary.ID

		app, err := appClient.Get(appGUID)
		if err != nil {
			return err
		}

		*obj = *app
		return nil
	}
}

func testAccCheckIBMCloudCFAppCreate(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}

resource "ibmcloud_cf_route" "route" {

			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			host              = "%s"
			
		}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello.zip"
	ports = [9080]
	wait_timeout = 90
	buildpack = "sdk-for-nodejs"
}`, cfOrganization, cfSpace, name, name)

}

func testAccCheckIBMCloudCFAppUpdate(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}

resource "ibmcloud_cf_route" "route" {

			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			host              = "%s"
			
		}
	
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello1.zip"
	instances = 1
	memory = 124
	disk_quota = 512
	buildpack = "nodejs_buildpack"
}`, cfOrganization, cfSpace, name, name)

}

func testAccCheckIBMCloudCFAppUpdateRoute(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}

resource "ibmcloud_cf_route" "route" {

			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			host              = "%s"
	
}

resource "ibmcloud_cf_route" "route1" {

			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			host              = "%s"
			
		}
	
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["${ibmcloud_cf_route.route.id}","${ibmcloud_cf_route.route1.id}"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello1.zip"
	instances = 1
	memory = 124
	disk_quota = 512
	buildpack = "nodejs_buildpack"
}`, cfOrganization, cfSpace, name, name+"-new", name)

}

func testAccCheckIBMCloudCFAppRemoveRoute(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

resource "ibmcloud_cf_route" "route1" {

			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			host              = "%s"
			
		}
	
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["${ibmcloud_cf_route.route1.id}"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello1.zip"
	instances = 1
	memory = 124
	disk_quota = 512
	buildpack = "nodejs_buildpack"
}`, cfOrganization, cfSpace, name+"-new", name)

}

// Service App Use cases

func testAccCheckIBMCloudCFAppBindService(name, serviceName string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cloudantNOSQLDB"
			plan              = "Lite"
			tags               = ["cluster-service"]
		}

	
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["4f1c7625-0096-4831-9817-e70541c15347"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello.zip"
	service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}"]
	ports = [9080]
	wait_timeout = 90
	buildpack = "sdk-for-nodejs"
}`, cfOrganization, cfSpace, serviceName, name)

}

func testAccCheckIBMCloudCFAppAddDeleteService(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

resource "ibmcloud_cf_service_instance" "service1" {
			name              = "terraform_service1"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cloudantNOSQLDB"
			plan              = "Lite"
			tags               = ["cluster-service1"]
		}
resource "ibmcloud_cf_service_instance" "service2" {
			name              = "terraform_service2"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service2"]
		}
	
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	route_guid = ["4f1c7625-0096-4831-9817-e70541c15347"]
	app_path = "/Users/sakshi/Documents/AlchemyProject/go_workspace/src/github.com/IBM-Bluemix/hello.zip"
	service_instance_guid = ["${ibmcloud_cf_service_instance.service1.id}","${ibmcloud_cf_service_instance.service2.id}"]
	ports = [9080]
	wait_timeout = 90
	buildpack = "sdk-for-nodejs"
}`, cfOrganization, cfSpace, name)

}
