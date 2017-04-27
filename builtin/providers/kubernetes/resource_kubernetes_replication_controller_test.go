package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func TestAccKubernetesReplicationController_basic(t *testing.T) {
	var conf api.ReplicationController
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.replication-controller", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "metadata.0.labels.app", "RC_for_UAT_test"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "metadata.0.name", "replication-controller"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.min_ready_seconds", "60"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.replicas", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.template.0.spec.0.containers.0.name", "uattest"),
				),
			},
			{
				Config: testAccKubernetesReplicationControllerConfig_updated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.replicas", "3"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.min_ready_seconds", "360"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.replication-controller", "spec.0.template.0.spec.0.active_deadline_seconds", "60"),
				),
			},
		},
	})
}

func testAccCheckKubernetesReplicationControllerDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*kubernetes.Clientset)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_replication_controller" {
			continue
		}
		namespace, name := idParts(rs.Primary.ID)
		resp, err := conn.CoreV1().ReplicationControllers(namespace).Get(name)
		if err == nil {
			if resp.Name == name {
				return fmt.Errorf("Replication Controller still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil

}

func testAccCheckKubernetesReplicationControllerExists(n string, obj *api.ReplicationController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)
		namespace, name := idParts(rs.Primary.ID)
		out, err := conn.CoreV1().ReplicationControllers(namespace).Get(name)
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesReplicationControllerConfig_basic() string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "replication-controller" {
	metadata {
	 labels {
		app  = "RC_for_UAT_test"
	 }
	 name = "replication-controller"
	}
	spec {
	 min_ready_seconds = 60
	 replicas = 2
	 template {
	  metadata {
		labels {
		  app  = "replicationUATAPP"
		}
	  }
	  spec {
		containers {
		  image = "nginx:1.7.9"
		  name = "uattest"
			resources{
				limits {
					cpu = "500m"
					memory = "128Mi"
				}
				requests {
					   memory = "64Mi"
        		cpu =  "250m"
				}
			}
		}
	  }
	}
   }
   delete_options {
   	orphan_dependents = false
   }
}
`)
}

func testAccKubernetesReplicationControllerConfig_updated() string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "replication-controller" {
  metadata {
    labels {
      app = "RC_for_UAT_test"
    }

    name = "replication-controller"
  }

  spec {
    min_ready_seconds = 360
    replicas          = 3

    template {
      metadata {
        labels {
          app = "replicationUATAPP"
        }
      }

      spec {
        active_deadline_seconds = 60

        containers {
          image = "nginx:1.7.9"
          name  = "uattest"

          resources {
            limits {
              cpu    = "500m"
              memory = "128Mi"
            }

            requests {
              memory = "64Mi"
              cpu    = "250m"
            }
          }
        }
      }
    }

  }
	delete_options {
      orphan_dependents = false
    }
}
`)
}
