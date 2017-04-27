package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func TestAccKubernetesService_basic(t *testing.T) {
	var conf api.Service

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKubernetesServiceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.service", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.provider", "ibm"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.name", "service"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.name", "http"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.target_port", "8989"),
				),
			},
			resource.TestStep{
				Config: testAccKubernetesServiceConfig_updateToLoadBalancerType(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.service", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.name", "service"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.type", "LoadBalancer"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.load_balancer_ip", "130.211.204.1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.selector.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.target_port", "8900"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.session_affinity", "ClientIP"),
				),
			},
			resource.TestStep{
				Config: testAccKubernetesServiceConfig_updateToNodePortType(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.service", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.name", "service"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.type", "NodePort"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.selector.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.name", "http"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.target_port", "8900"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.ports.0.node_port", "30003"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.session_affinity", "None"),
				),
			},
			resource.TestStep{
				Config: testAccKubernetesServiceConfig_externalNametType(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.service", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.name", "service-external"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.type", "ExternalName"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.external_name", "my.database.example.com"),
				),
			},
			resource.TestStep{
				Config: testAccKubernetesServiceConfig_loadbalacersourcerange(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.service", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.labels.app", "ngix"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "metadata.0.name", "service-lbs"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.type", "LoadBalancer"),
					resource.TestCheckResourceAttr("kubernetes_service.service", "spec.0.load_balancer_source_ranges.#", "2"),
				),
			},
		},
	})
}

func testAccCheckKubernetesServiceDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*kubernetes.Clientset)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_service" {
			continue
		}
		namespace, name := idParts(rs.Primary.ID)
		resp, err := conn.Services(namespace).Get(name)
		if err == nil {
			if resp.Name == name {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil

}

func testAccCheckKubernetesServiceExists(n string, obj *api.Service) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)
		namespace, name := idParts(rs.Primary.ID)
		out, err := conn.Services(namespace).Get(name)
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesServiceConfig_basic() string {
	return fmt.Sprintf(`
resource "kubernetes_service" "service" {
	metadata {
		labels {
			provider  = "ibm"
		}
		name = "service"
	}
	spec {
		ports {
			name = "http"
			port = 80
			protocol = "TCP"
			target_port = 8989
		}
		external_ips = ["172.16.127.2", "172.16.127.1"]
	}
}
	`)
}

func testAccKubernetesServiceConfig_updateToLoadBalancerType() string {
	return fmt.Sprintf(`
resource "kubernetes_service" "service" {
	metadata {
		labels {
			app  = "ngix"
		}
		name = "service"
	}
	spec {

        type = "LoadBalancer"
		load_balancer_ip = "130.211.204.1"
		selector {
			app =  "ngix"
		
		}
		session_affinity = "ClientIP"
		ports {
			port = 8080
			target_port = 8900
		}
		external_ips = ["172.16.127.2", "172.16.127.1", "172.16.127.0"]
	}
}
	`)
}

func testAccKubernetesServiceConfig_updateToNodePortType() string {
	return fmt.Sprintf(`
resource "kubernetes_service" "service" {
	metadata {
		labels {
			app  = "ngix"
		}
		name = "service"
	}
	spec {

        type = "NodePort"
		selector {
			app =  "ngix"
		
		}
		session_affinity = "None"
		ports {
			name = "http"
			node_port = 30003
			port = 8080
			protocol = "TCP"
			target_port = 8900
		}
	}
}
	`)
}

func testAccKubernetesServiceConfig_externalNametType() string {
	return fmt.Sprintf(`
resource "kubernetes_service" "service" {
	metadata {
		labels {
			app  = "ngix"
		}
		name = "service-external"
	}
	spec {
        type = "ExternalName"
		external_name = "my.database.example.com"
	}
}
	`)
}

func testAccKubernetesServiceConfig_loadbalacersourcerange() string {
	return fmt.Sprintf(`
resource "kubernetes_service" "service" {
	metadata {
		labels {
			app  = "ngix"
		}
		name = "service-lbs"
	}
	spec {

        type = "LoadBalancer"
		load_balancer_source_ranges = ["130.211.204.1/32", "130.211.204.2/32"]
		selector {
			app =  "ngix"
		
		}
		ports {
			port = 8080
			target_port = 8900
		}
	}
}
	`)
}
