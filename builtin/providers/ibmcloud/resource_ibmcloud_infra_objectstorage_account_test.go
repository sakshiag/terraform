package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIBMCloudInfraObjectStorageAccount_Basic(t *testing.T) {
	var accountName string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraObjectStorageAccountDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraObjectStorageAccountConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraObjectStorageAccountExists("ibmcloud_infra_objectstorage_account.testacc_foobar", &accountName),
					testAccCheckIBMCloudInfraObjectStorageAccountAttributes(&accountName),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraObjectStorageAccountDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckIBMCloudInfraObjectStorageAccountExists(n string, accountName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		*accountName = rs.Primary.ID

		return nil
	}
}

func testAccCheckIBMCloudInfraObjectStorageAccountAttributes(accountName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *accountName == "" {
			return fmt.Errorf("No object storage account name")
		}

		return nil
	}
}

var testAccCheckIBMCloudInfraObjectStorageAccountConfig_basic = `
resource "ibmcloud_infra_objectstorage_account" "testacc_foobar" {
}`
