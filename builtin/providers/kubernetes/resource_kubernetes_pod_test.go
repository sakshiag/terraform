package kubernetes

import (
	"fmt"
	"testing"

	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKubernetesPod_basic(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	secretName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	configMapName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	imageName1 := "nginx:1.7.9"
	imageName2 := "nginx:1.11"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigBasic(secretName, configMapName, podName, imageName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "metadata.0.annotations.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "metadata.0.labels.app", "pod_label"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "metadata.0.name", podName),
					resource.TestCheckResourceAttrSet("kubernetes_pod.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_pod.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_pod.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_pod.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.env.0.value_from.0.secret_key_ref.0.name", secretName),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.env.1.value_from.0.config_map_key_ref.0.name", configMapName),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.image", imageName1),
				),
			},
			{
				Config: testAccKubernetesPodConfigBasic(secretName, configMapName, podName, imageName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.image", imageName2),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_pod_security_context(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithSecurityContext(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.security_context.0.run_as_non_root", "true"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.security_context.0.supplemental_groups.#", "1"),
				),
			},
		},
	})
}

func testAccCheckKubernetesPodExists(n string, obj *api.Pod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)

		namespace, name := idParts(rs.Primary.ID)
		out, err := conn.CoreV1().Pods(namespace).Get(name)
		if err != nil {
			return err
		}
		*obj = *out
		return nil
	}
}

func testAccKubernetesPodConfigBasic(secretName, configMapName, podName, imageName string) string {
	return fmt.Sprintf(`

resource "kubernetes_secret" "test" {
  metadata {
    name = "%s"
  }

  data {
    one = "first"
  }
}

resource "kubernetes_config_map" "test" {
  metadata {
    name = "%s"
  }

  data {
    one = "ONE"
  }
}

resource "kubernetes_pod" "test" {
  metadata {
    labels {
      app = "pod_label"
    }

    name = "%s"
  }

  spec {
    containers {
      image = "%s"
      name  = "containername"

      env = [{
        name = "EXPORTED_VARIBALE_FROM_SECRET"

        value_from {
          secret_key_ref {
            name = "${kubernetes_secret.test.metadata.0.name}"
            key  = "one"
          }
        }
      },
        {
          name = "EXPORTED_VARIBALE_FROM_CONFIG_MAP"

          value_from {
            config_map_key_ref {
              name = "${kubernetes_config_map.test.metadata.0.name}"
              key  = "one"
            }
          }
        },
      ]
    }

    volumes = [{
      name = "db"

      secret = {
        secret_name = "${kubernetes_secret.test.metadata.0.name}"
      }
    }]
  }
}
	`, secretName, configMapName, podName, imageName)
}

func testAccKubernetesPodConfigWithSecurityContext(podName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_pod" "test" {
  metadata {
    labels {
      app = "pod_label"
    }

    name = "%s"
  }

  spec {
		security_context{
			run_as_non_root = true
			supplemental_groups = [101]
		}
    containers {
      image = "%s"
      name  = "containername"
    }
  }
}
	`, podName, imageName)
}
