package ibmcloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
)

func TestAccIBMCloudInfraProvisioningHook_Basic(t *testing.T) {
	var hook datatypes.Provisioning_Hook

	hookName1 := fmt.Sprintf("%s%s", "tfuathook", acctest.RandString(10))
	hookName2 := fmt.Sprintf("%s%s", "tfuathook", acctest.RandString(10))
	uri1 := "http://www.weather.com"
	uri2 := "https://www.ibm.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraProvisioningHookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraProvisioningHookConfig(hookName1, uri1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraProvisioningHookExists("ibmcloud_infra_provisioning_hook.test-provisioning-hook", &hook),
					testAccCheckIBMCloudInfraProvisioningHookAttributes(&hook, hookName1, uri1),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_provisioning_hook.test-provisioning-hook", "name", hookName1),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_provisioning_hook.test-provisioning-hook", "uri", uri1),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraProvisioningHookConfig(hookName2, uri2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraProvisioningHookExists("ibmcloud_infra_provisioning_hook.test-provisioning-hook", &hook),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_provisioning_hook.test-provisioning-hook", "name", hookName2),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_provisioning_hook.test-provisioning-hook", "uri", uri2),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraProvisioningHookDestroy(s *terraform.State) error {
	service := services.GetProvisioningHookService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_infra_provisioning_hook" {
			continue
		}

		hookId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the provisioning hook
		_, err := service.Id(hookId).GetObject()

		if err == nil {
			return fmt.Errorf("Provisioning Hook still exists")
		}
	}

	return nil
}

func testAccCheckIBMCloudInfraProvisioningHookAttributes(hook *datatypes.Provisioning_Hook, name, uri string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *hook.Name != name {
			return fmt.Errorf("Bad name: %s", *hook.Name)
		}

		if *hook.Uri != uri {
			return fmt.Errorf("Bad uri: %s", *hook.Uri)
		}

		return nil
	}
}

func testAccCheckIBMCloudInfraProvisioningHookExists(n string, hook *datatypes.Provisioning_Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		hookId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetProvisioningHookService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		foundHook, err := service.Id(hookId).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundHook.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*hook = foundHook

		return nil
	}
}

func testAccCheckIBMCloudInfraProvisioningHookConfig(name, uri string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_provisioning_hook" "test-provisioning-hook" {
    name = "%s"
    uri = "%s"
}`, name, uri)
}
