package ibmcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAccIBMCloudInfraScalePolicy_Basic(t *testing.T) {
	var scalepolicy datatypes.Scale_Policy
	groupname := fmt.Sprintf("terraformuat_%d", acctest.RandInt())
	hostname := acctest.RandString(16)
	policyname := acctest.RandString(16)
	updatedpolicyname := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraScalePolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudInfraScalePolicyConfig_basic(groupname, hostname, policyname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraScalePolicyExists("ibmcloud_infra_scale_policy.sample-http-cluster-policy", &scalepolicy),
					testAccCheckIBMCloudInfraScalePolicyAttributes(&scalepolicy, policyname),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "name", policyname),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "scale_type", "RELATIVE"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "scale_amount", "1"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "cooldown", "30"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "triggers.#", "3"),
					testAccCheckIBMCloudInfraScalePolicyContainsRepeatingTriggers(&scalepolicy, 2, "0 1 ? * MON,WED *"),
					testAccCheckIBMCloudInfraScalePolicyContainsResourceUseTriggers(&scalepolicy, 120, "80"),
					testAccCheckIBMCloudInfraScalePolicyContainsOneTimeTriggers(&scalepolicy, testOnetimeTriggerDate),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudInfraScalePolicyConfig_updated(groupname, hostname, updatedpolicyname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraScalePolicyExists("ibmcloud_infra_scale_policy.sample-http-cluster-policy", &scalepolicy),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "name", updatedpolicyname),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "scale_type", "ABSOLUTE"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "scale_amount", "2"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "cooldown", "35"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_scale_policy.sample-http-cluster-policy", "triggers.#", "3"),
					testAccCheckIBMCloudInfraScalePolicyContainsRepeatingTriggers(&scalepolicy, 2, "0 1 ? * MON,WED,SAT *"),
					testAccCheckIBMCloudInfraScalePolicyContainsResourceUseTriggers(&scalepolicy, 130, "90"),
					testAccCheckIBMCloudInfraScalePolicyContainsOneTimeTriggers(&scalepolicy, testOnetimeTriggerUpdatedDate),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraScalePolicyDestroy(s *terraform.State) error {
	service := services.GetScalePolicyService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_infra_scale_policy" {
			continue
		}

		scalepolicyId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the key
		_, err := service.Id(scalepolicyId).GetObject()

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for Auto Scale Policy (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudInfraScalePolicyContainsResourceUseTriggers(scalePolicy *datatypes.Scale_Policy, period int, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false

		for _, scaleResourceUseTrigger := range scalePolicy.ResourceUseTriggers {
			for _, scaleResourceUseWatch := range scaleResourceUseTrigger.Watches {
				if *scaleResourceUseWatch.Metric == "host.cpu.percent" && *scaleResourceUseWatch.Operator == ">" &&
					*scaleResourceUseWatch.Period == period && *scaleResourceUseWatch.Value == value {
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("Resource use trigger not found in scale policy")

		}

		return nil
	}
}

func testAccCheckIBMCloudInfraScalePolicyContainsRepeatingTriggers(scalePolicy *datatypes.Scale_Policy, typeId int, schedule string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false

		for _, scaleRepeatingTrigger := range scalePolicy.RepeatingTriggers {
			if *scaleRepeatingTrigger.TypeId == typeId && *scaleRepeatingTrigger.Schedule == schedule {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Repeating trigger %d with schedule %s not found in scale policy", typeId, schedule)

		}

		return nil
	}
}

func testAccCheckIBMCloudInfraScalePolicyContainsOneTimeTriggers(scalePolicy *datatypes.Scale_Policy, testOnetimeTriggerDate string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false
		const IBMCloudInfraTimeFormat = "2006-01-02T15:04:05-07:00"
		utcLoc, _ := time.LoadLocation("UTC")

		for _, scaleOneTimeTrigger := range scalePolicy.OneTimeTriggers {
			if scaleOneTimeTrigger.Date.In(utcLoc).Format(IBMCloudInfraTimeFormat) == testOnetimeTriggerDate {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("One time trigger with date %s not found in scale policy", testOnetimeTriggerDate)
		}

		return nil

	}
}

func testAccCheckIBMCloudInfraScalePolicyAttributes(scalepolicy *datatypes.Scale_Policy, policyname string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *scalepolicy.Name != policyname {
			return fmt.Errorf("Bad name: %s", *scalepolicy.Name)
		}

		return nil
	}
}

func testAccCheckIBMCloudInfraScalePolicyExists(n string, scalepolicy *datatypes.Scale_Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		scalepolicyId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetScalePolicyService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		foundScalePolicy, err := service.Id(scalepolicyId).Mask(strings.Join(IBMCloudInfraScalePolicyObjectMask, ",")).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundScalePolicy.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*scalepolicy = foundScalePolicy
		return nil
	}
}

func testAccCheckIBMCloudInfraScalePolicyConfig_basic(groupname, hostname, policyname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_scale_group" "sample-http-cluster-with-policy" {
    name = "%s"
    regional_group = "na-usa-central-1"
    cooldown = 30
    minimum_member_count = 1
    maximum_member_count = 10
    termination_policy = "CLOSEST_TO_NEXT_CHARGE"
    virtual_guest_member_template = {
        hostname = "%s"
        domain = "terraformuat.ibm.com"
        cores = 1
        memory = 4096
        network_speed = 1000
        hourly_billing = true
        os_reference_code = "DEBIAN_7_64"
        local_disk = false
        datacenter = "dal09"
    }
}

resource "ibmcloud_infra_scale_policy" "sample-http-cluster-policy" {
    name = "%s"
    scale_type = "RELATIVE"
    scale_amount = 1
    cooldown = 30
    scale_group_id = "${ibmcloud_infra_scale_group.sample-http-cluster-with-policy.id}"
    triggers = {
        type = "RESOURCE_USE"
        watches = {

                    metric = "host.cpu.percent"
                    operator = ">"
                    value = "80"
                    period = 120
        }
    }
    triggers = {
        type = "ONE_TIME"
        date = "%s"
    }
    triggers = {
        type = "REPEATING"
        schedule = "0 1 ? * MON,WED *"
    }

}`, groupname, hostname, policyname, testOnetimeTriggerDate)
}

const IBMCloudInfraTestTimeFormat = string("2006-01-02T15:04:05-07:00")

var utcLoc, _ = time.LoadLocation("UTC")

var testOnetimeTriggerDate = time.Now().In(utcLoc).AddDate(0, 0, 1).Format(IBMCloudInfraTestTimeFormat)

func testAccCheckIBMCloudInfraScalePolicyConfig_updated(groupname, hostname, updatedpolicyname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_scale_group" "sample-http-cluster-with-policy" {
    name = "%s"
    regional_group = "na-usa-central-1"
    cooldown = 30
    minimum_member_count = 1
    maximum_member_count = 10
    termination_policy = "CLOSEST_TO_NEXT_CHARGE"
    virtual_guest_member_template = {
        hostname = "%s"
        domain = "terraformuat.ibm.com"
        cores = 1
        memory = 4096
        network_speed = 1000
        hourly_billing = true
        os_reference_code = "DEBIAN_7_64"
        local_disk = false
        datacenter = "dal09"
    }
}
resource "ibmcloud_infra_scale_policy" "sample-http-cluster-policy" {
    name = "%s"
    scale_type = "ABSOLUTE"
    scale_amount = 2
    cooldown = 35
    scale_group_id = "${ibmcloud_infra_scale_group.sample-http-cluster-with-policy.id}"
    triggers = {
        type = "RESOURCE_USE"
        watches = {

                    metric = "host.cpu.percent"
                    operator = ">"
                    value = "90"
                    period = 130
        }
    }
    triggers = {
        type = "REPEATING"
        schedule = "0 1 ? * MON,WED,SAT *"
    }
    triggers = {
        type = "ONE_TIME"
        date = "%s"
    }
}`, groupname, hostname, updatedpolicyname, testOnetimeTriggerUpdatedDate)
}

var testOnetimeTriggerUpdatedDate = time.Now().In(utcLoc).AddDate(0, 0, 2).Format(IBMCloudInfraTestTimeFormat)
