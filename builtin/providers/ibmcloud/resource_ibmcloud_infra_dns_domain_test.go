package ibmcloud

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func TestAccIBMCloudInfraDNSDomain_Basic(t *testing.T) {
	var dns_domain datatypes.Dns_Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(config, domainName1, target1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainAttributes(&dns_domain),
					saveIBMCloudInfraDNSDomainId(&dns_domain, &firstDnsId),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "name", domainName1),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "target", target1),
				),
				Destroy: false,
			},
			{
				Config: fmt.Sprintf(config, domainName2, target1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainAttributes(&dns_domain),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "name", domainName2),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "target", target1),
					testAccCheckIBMCloudInfraDNSDomainChanged(&dns_domain),
				),
				Destroy: false,
			},
			{
				Config: fmt.Sprintf(config, domainName2, target2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainAttributes(&dns_domain),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "name", domainName2),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_dns_domain.acceptance_test_dns_domain-1", "target", target2),
				),
				Destroy: false,
			},
		},
	})
}

func testAccCheckIBMCloudInfraDNSDomainDestroy(s *terraform.State) error {
	service := services.GetDnsDomainService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_infra_dns_domain" {
			continue
		}

		dnsId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the domain
		_, err := service.Id(dnsId).GetObject()

		if err == nil {
			return fmt.Errorf("Dns Domain with id %d still exists", dnsId)
		}
	}

	return nil
}

func testAccCheckIBMCloudInfraDNSDomainAttributes(dns *datatypes.Dns_Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name := sl.Get(dns.Name); name == "" {
			return errors.New("Empty dns domain name")
		}

		// find a record with host @; that will have the current target.
		foundTarget := false
		for _, record := range dns.ResourceRecords {
			if *record.Type == "a" && *record.Host == "@" {
				foundTarget = true
				break
			}
		}

		if !foundTarget {
			return fmt.Errorf("Target record not found for dns domain %s (%d)", sl.Get(dns.Name), sl.Get(dns.Id))
		}

		if id := sl.Get(dns.Id); id == 0 {
			return fmt.Errorf("Bad dns domain id: %d", id)
		}

		return nil
	}
}

func saveIBMCloudInfraDNSDomainId(dns *datatypes.Dns_Domain, id_holder *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*id_holder = *dns.Id

		return nil
	}
}

func testAccCheckIBMCloudInfraDNSDomainChanged(dns *datatypes.Dns_Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		service := services.GetDnsDomainService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

		_, err := service.Id(firstDnsId).Mask(
			"id,name,updateDate,resourceRecords",
		).GetObject()
		if err == nil {
			return fmt.Errorf("Dns domain with id %d still exists", firstDnsId)
		}

		return nil
	}
}

func testAccCheckIBMCloudInfraDNSDomainExists(n string, dns_domain *datatypes.Dns_Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		dns_id, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetDnsDomainService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		found_domain, err := service.Id(dns_id).Mask(
			"id,name,updateDate,resourceRecords",
		).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*found_domain.Id)) != rs.Primary.ID {
			return errors.New("Record not found")
		}

		*dns_domain = found_domain

		return nil
	}
}

var config = `
resource "ibmcloud_infra_dns_domain" "acceptance_test_dns_domain-1" {
	name = "%s"
	target = "%s"
}
`

var domainName1 = fmt.Sprintf("tfuatdomain%s.com", acctest.RandString(10))
var domainName2 = fmt.Sprintf("tfuatdomain%s.com", acctest.RandString(10))
var target1 = "172.16.0.100"
var target2 = "172.16.0.101"
var firstDnsId = 0
