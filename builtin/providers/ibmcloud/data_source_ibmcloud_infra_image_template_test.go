package ibmcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudInfraImageTemplateDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Tests looking up private or shared images
			{
				Config: testAccCheckIBMCloudInfraImageTemplateDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibmcloud_infra_image_template.tfacc_img_tmpl",
						"name",
						"jumpbox",
					),
					resource.TestMatchResourceAttr(
						"data.ibmcloud_infra_image_template.tfacc_img_tmpl",
						"id",
						regexp.MustCompile("^[0-9]+$"),
					),
				),
			},
			// Tests looking up a public image
			{
				Config: testAccCheckIBMCloudInfraImageTemplateDataSourceConfig_basic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibmcloud_infra_image_template.tfacc_img_tmpl",
						"name",
						"RightImage_Ubuntu_12.04_amd64_v13.5",
					),
					resource.TestMatchResourceAttr(
						"data.ibmcloud_infra_image_template.tfacc_img_tmpl",
						"id",
						regexp.MustCompile("^[0-9]+$"),
					),
				),
			},
		},
	})
}

const testAccCheckIBMCloudInfraImageTemplateDataSourceConfig_basic = `
data "ibmcloud_infra_image_template" "tfacc_img_tmpl" {
    name = "jumpbox"
}
`

const testAccCheckIBMCloudInfraImageTemplateDataSourceConfig_basic2 = `
data "ibmcloud_infra_image_template" "tfacc_img_tmpl" {
    name = "RightImage_Ubuntu_12.04_amd64_v13.5"
}
`
