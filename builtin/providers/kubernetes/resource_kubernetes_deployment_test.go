package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func TestAccKubernetesDeployment_basic(t *testing.T) {
	var conf v1beta1.Deployment
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.app", "test-for-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", "test-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.replicas", "3"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.name", "test1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.image", "nginx:1.7.9"),
				),
			},
			{
				Config: testAccKubernetesDeploymentConfig_updatedReplica(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.app", "test-for-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", "test-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.replicas", "4"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.name", "test1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.image", "nginx:1.7.9"),
				),
			},
			{
				Config: testAccKubernetesDeploymentConfig_updatedImage(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.app", "test-for-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", "test-deployment"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.replicas", "4"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.name", "test1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.containers.0.image", "nginx:1.11"),
				),
			},
		},
	})
}

func testAccCheckKubernetesDeploymentDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*kubernetes.Clientset)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_deployment" {
			continue
		}
		namespace, name := idParts(rs.Primary.ID)
		resp, err := conn.ExtensionsV1beta1().Deployments(namespace).Get(name)
		if err == nil {
			if resp.Name == name {
				return fmt.Errorf("Depolyment still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil

}

func testAccCheckKubernetesDeploymentExists(n string, obj *v1beta1.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)
		namespace, name := idParts(rs.Primary.ID)
		out, err := conn.ExtensionsV1beta1().Deployments(namespace).Get(name)
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesDeploymentConfig_basic() string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			app  = "test-for-deployment"
		}
		name = "test-deployment"
	}
	spec {
		replicas =  3
		template {
			metadata {
					labels {
						app  = "tempalteapp"
					}
			}
			spec {
				containers = [{
					image = "nginx:1.7.9"
					name = "test1"
				}]
			}
		}
	}
	delete_options {
   	orphan_dependents = false
   }
	
}
	`)
}

func testAccKubernetesDeploymentConfig_updatedReplica() string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			app  = "test-for-deployment"
		}
		name = "test-deployment"
	}
	spec {
		replicas =  4
		template {
			metadata {
					labels {
						app  = "tempalteapp"
						
					}
			}
			spec {
			
				containers = [{
					image = "nginx:1.7.9"
					name = "test1"
				}]
			}
		}
	}
	delete_options {
   	orphan_dependents = false
   }
	
}
	`)
}

func testAccKubernetesDeploymentConfig_updatedImage() string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
	metadata {
		annotations {
			TestAnnotationOne = "one"
			Different = "1234"
		}
		labels {
			app  = "test-for-deployment"
		}
		name = "test-deployment"
	}
	spec {
		replicas =  4
		template {
			metadata {
					labels {
						app  = "tempalteapp"
						
					}
			}
			spec {
			
				containers = [{
					image = "nginx:1.11"
					name = "test1"
					ports {
        				container_port = 80
        			}
				}]
			}
		}
	}
	delete_options {
   	orphan_dependents = false
   }
	
}
	`)
}
