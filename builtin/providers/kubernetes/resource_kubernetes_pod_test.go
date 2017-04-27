package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKubernetesPod_basic(t *testing.T) {
	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	secretName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfig_basic(secretName, podName),
			},
			{
				Config: testAccKubernetesPodConfig_imageUpdate(secretName, podName),
			},
		},
	})
}

func testAccKubernetesPodConfig_basic(secretName, podName string) string {
	return fmt.Sprintf(`

resource "kubernetes_secret" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			TestLabelOne = "one"
			TestLabelThree = "three"
		}
		name = "%s"
	}
	data {
		one = "first"
		two = "second"
		nine = "ninth"
	}
}

resource "kubernetes_config_map" "test" {
	metadata {
		name = "some-config-map"
	}
	data {
		one = "ONE"
	}
}


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
			name = "containername"
			env = [{
				name = "EXPORTED_VARIBALE_FROM_SECRET"
				value_from {
					secret_key_ref{
						name = "${kubernetes_secret.test.metadata.0.name}"
						key = "one"
					}
				}
			},
			{
				name = "EXPORTED_VARIBALE_FROM_CONFIG_MAP"
				value_from {
					config_map_key_ref{
						name = "${kubernetes_config_map.test.metadata.0.name}"
						key = "one"
					}
				}
			}]
		}
		volumes =  [{
        name = "mycloudant",
          secret =  {
            secret_name =  "${kubernetes_secret.test.metadata.0.name}"
          }
         }]
	}
}
	`, secretName, podName)
}

func testAccKubernetesPodConfig_imageUpdate(secretName, podName string) string {
	return fmt.Sprintf(`

resource "kubernetes_secret" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			TestLabelOne = "one"
			TestLabelThree = "three"
		}
		name = "%s"
	}
	data {
		one = "first"
		two = "second"
		nine = "ninth"
	}
}

resource "kubernetes_config_map" "test" {
	metadata {
		name = "some-config-map"
	}
	data {
		one = "ONE"
	}
}


resource "kubernetes_pod" "test" {
	metadata {
		labels {
			app  = "pod_label"
		}
		name = "%s"
	}
	spec {
		containers {
			image = "nginx:1.11"
			name = "containername"
			env = [{
				name = "EXPORTED_VARIBALE_FROM_SECRET"
				value_from {
					secret_key_ref{
						name = "${kubernetes_secret.test.metadata.0.name}"
						key = "one"
					}
				}
			},
			{
				name = "EXPORTED_VARIBALE_FROM_CONFIG_MAP"
				value_from {
					config_map_key_ref{
						name = "${kubernetes_config_map.test.metadata.0.name}"
						key = "one"
					}
				}
			}]
		}
		volumes =  [{
        name = "mycloudant",
          secret =  {
            secret_name =  "${kubernetes_secret.test.metadata.0.name}"
          }
         }]
	}
}
	`, secretName, podName)
}
