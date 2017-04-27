package ibmcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSoftLayerDnsDomainDataSource_Basic(t *testing.T) {

	var domainName = acctest.RandString(16) + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckSoftLayerDnsDomainDataSourceConfig_basic, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibmcloud_infra_dns_domain.domain_id", "name", domainName),
					resource.TestMatchResourceAttr("data.ibmcloud_infra_dns_domain.domain_id", "id", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

// The datasource to apply
const testAccCheckSoftLayerDnsDomainDataSourceConfig_basic = `
resource "ibmcloud_infra_dns_domain" "ds_domain_test" {
	name = "%s"
}
data "ibmcloud_infra_dns_domain" "domain_id" {
    name = "${ibmcloud_infra_dns_domain.ds_domain_test.name}"
}
`
