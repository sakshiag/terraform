package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
)

func TestAccSoftLayerDnsDomain_Basic(t *testing.T) {
	var dns_domain datatypes.SoftLayer_Dns_Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSoftLayerDnsDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerDnsDomainConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerDnsDomainExists("softlayer_dns_domain.test_dns_domain-1", &dns_domain),
					testAccCheckSoftLayerDnsDomainAttributes(&dns_domain),
				),
			},
		},
	})
}

func testAccCheckSoftLayerDnsDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).dnsDomainService

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "softlayer_dns_domain" {
			continue
		}

		dnsId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the domain
		_, err := client.GetObject(dnsId)

		if err == nil {
			fmt.Errorf("Dns Domain with id %d does not exist", dnsId)
		}
	}

	return nil
}

func testAccCheckSoftLayerDnsDomainAttributes(dns *datatypes.SoftLayer_Dns_Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dns.Name != test_dns_domain_name {
			return fmt.Errorf("Bad dns domain name: %s  Expected: %s", dns.Name, test_dns_domain_name)
		}

		if dns.Serial == 0 {
			return fmt.Errorf("Bad dns domain serial: %d", dns.Serial)
		}

		if dns.Id == 0 {
			return fmt.Errorf("Bad dns domain id: %d", dns.Id)
		}

		return nil
	}
}

func testAccCheckSoftLayerDnsDomainExists(n string, dns_domain *datatypes.SoftLayer_Dns_Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		dns_id, _ := strconv.Atoi(rs.Primary.ID)

		client := testAccProvider.Meta().(*Client).dnsDomainService
		found_domain, err := client.GetObject(dns_id)

		if err != nil {
			return err
		}

		if strconv.Itoa(int(found_domain.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*dns_domain = found_domain

		return nil
	}
}

var testAccCheckSoftLayerDnsDomainConfig_basic = fmt.Sprintf(`
resource "softlayer_dns_domain" "test_dns_domain-1" {
	name = "%s"
}
`, test_dns_domain_name)

var test_dns_domain_name = "zxczcxzxc.com"