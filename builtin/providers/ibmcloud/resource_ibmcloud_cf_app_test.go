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
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
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
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 90
	buildpack = "sdk-for-nodejs"
}`, cfOrganization, cfSpace, name)

}

func testAccCheckIBMCloudCFAppUpdate(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 90
	buildpack = "sdk-for-nodejs"
	instances = 1
	disk_quota = 512
	memory = 128
	ports = [9080]
	instances = 1
	disk_quota = 512
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, name)

}
