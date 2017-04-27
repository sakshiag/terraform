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
)

func TestAccIBMCloudInfraDNSDomainRecord_Basic(t *testing.T) {
	var dns_domain datatypes.Dns_Domain
	var dns_domain_record datatypes.Dns_Domain_ResourceRecord

	domainName := fmt.Sprintf("tfuatdomainr%s.ibm.com", acctest.RandString(10))
	host1 := acctest.RandString(10) + "ibm.com"
	host2 := acctest.RandString(10) + "ibm.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCloudInfraDNSDomainRecordConfigBasic(domainName, host1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.test_dns_domain_records", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordA", &dns_domain_record),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "data", "127.0.0.1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "expire", "900"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "minimum_ttl", "90"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "mx_priority", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "refresh", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "host", host1),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "responsible_person", "user@softlayer.com"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "ttl", "900"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "retry", "1"),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "type", "a"),
				),
			},
			{
				Config: testAccCheckIBMCloudInfraDNSDomainRecordConfigBasic(domainName, host2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.test_dns_domain_records", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordA", &dns_domain_record),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordA", "host", host2),
				),
			},
		},
	})
}

func TestAccIBMCloudInfraDNSDomainRecord_Types(t *testing.T) {
	var dns_domain datatypes.Dns_Domain
	var dns_domain_record datatypes.Dns_Domain_ResourceRecord

	domainName := acctest.RandString(10) + "dnstest.ibm.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckIBMCloudInfraDNSDomainRecordConfig_all_types, domainName, "_tcp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.test_dns_domain_record_types", &dns_domain),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordA", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordAAAA", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordCNAME", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordMX", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordSPF", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordTXT", &dns_domain_record),
					testAccCheckIBMCloudInfraDNSDomainRecordExists("ibmcloud_infra_dns_domain_record.recordSRV", &dns_domain_record),
				),
			},

			{
				Config: fmt.Sprintf(testAccCheckIBMCloudInfraDNSDomainRecordConfig_all_types, domainName, "_udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraDNSDomainExists("ibmcloud_infra_dns_domain.test_dns_domain_record_types", &dns_domain),
					resource.TestCheckResourceAttr("ibmcloud_infra_dns_domain_record.recordSRV", "protocol", "_udp"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraDNSDomainRecordExists(n string, dns_domain_record *datatypes.Dns_Domain_ResourceRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		dns_id, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetDnsDomainResourceRecordService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		found_domain_record, err := service.Id(dns_id).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*found_domain_record.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record %d not found", dns_id)
		}

		*dns_domain_record = found_domain_record

		return nil
	}
}

func testAccCheckIBMCloudInfraDNSDomainRecordConfigBasic(domainName, hostname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_dns_domain" "test_dns_domain_records" {
	name = "%s"
	target = "172.16.0.100"
}

resource "ibmcloud_infra_dns_domain_record" "recordA" {
    data = "127.0.0.1"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_records.id}"
    expire = 900
    minimum_ttl = 90
    mx_priority = 1
    refresh = 1
    host = "%s"
    responsible_person = "user@softlayer.com"
    ttl = 900
    retry = 1
    type = "a"
}`, domainName, hostname)
}

var testAccCheckIBMCloudInfraDNSDomainRecordConfig_all_types = `
resource "ibmcloud_infra_dns_domain" "test_dns_domain_record_types" {
	name = "%s"
	target = "172.16.12.100"
}

resource "ibmcloud_infra_dns_domain_record" "recordA" {
    data = "127.0.0.1"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta.com"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "a"
}

resource "ibmcloud_infra_dns_domain_record" "recordAAAA" {
    data = "fe80:0000:0000:0000:0202:b3ff:fe1e:8329"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-2.com"
    responsible_person = "user2changed@softlayer.com"
    ttl = 1000
    type = "aaaa"
}

resource "ibmcloud_infra_dns_domain_record" "recordCNAME" {
    data = "testsssaaaass.com"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-cname.com"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "cname"
}

resource "ibmcloud_infra_dns_domain_record" "recordMX" {
    data = "email.example.com"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-mx.com"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "mx"
}

resource "ibmcloud_infra_dns_domain_record" "recordSPF" {
    data = "v=spf1 mx:mail.example.org ~all"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-spf"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "spf"
}

resource "ibmcloud_infra_dns_domain_record" "recordTXT" {
    data = "127.0.0.1"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-txt.com"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "txt"
}

resource "ibmcloud_infra_dns_domain_record" "recordSRV" {
    data = "ns1.example.org"
    domain_id = "${ibmcloud_infra_dns_domain.test_dns_domain_record_types.id}"
    host = "hosta-srv.com"
    responsible_person = "user@softlayer.com"
    ttl = 900
    type = "srv"
	port = 8080
	priority = 3
	protocol = "%s"
	weight = 3
	service = "_mail"
}
`
