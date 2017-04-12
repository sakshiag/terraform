package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKubernetesPod_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-pod-test-%s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testAccKubernetesPodConfigBasic(name),
				Destroy: false,
			},
		},
	})
}

func testAccKubernetesPodConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_pod" "test" {
	metadata {
		labels {
			app  = "pod_label"
		}
		name = "%s"
	}
	spec {
		containers {
		    image = "nginx:1.7.9"
			name = "cname"

			ports {
				container_port = 9080
			}
		}
	}
}
`, name)
}
