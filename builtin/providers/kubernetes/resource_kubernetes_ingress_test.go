package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func TestAccKubernetesIngress_basic(t *testing.T) {
	var conf v1beta1.Ingress
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesIngressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesIngressConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesIngressExists("kubernetes_ingress.ingress", &conf),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.labels.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.labels.provider", "kubernetes"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.name", "echo-map"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.host", "foo.bar.com"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.path", "/foo"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.backend.0.service_name", "echoheaders-x"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.backend.0.service_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.host", "bar.baz.com"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.path", "/bar"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.backend.0.service_name", "echoheaders-y"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.backend.0.service_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.1.path", "/foo"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.1.backend.0.service_name", "echoheaders-x"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.1.backend.0.service_port", "80"),
				),
			},
			{
				Config: testAccKubernetesIngressConfig_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesIngressExists("kubernetes_ingress.ingress", &conf),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.labels.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.labels.provider", "kubernetes"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.labels.app", "echoMap"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "metadata.0.name", "echo-map"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.host", "foo1.bar.com"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.path", "/foo1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.backend.0.service_name", "echoheaders-x1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.0.http.0.paths.0.backend.0.service_port", "90"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.host", "bar1.baz.com"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.path", "/bar1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.backend.0.service_name", "echoheaders-y1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.rules.1.http.0.paths.0.backend.0.service_port", "90"),
				),
			},
			{
				Config: testAccKubernetesIngressConfig_tlsSpecified(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesIngressExists("kubernetes_ingress.ingress", &conf),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.tls.0.hosts.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.tls.0.secret_name", "foobar1-ssl"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.tls.1.hosts.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_ingress.ingress", "spec.0.tls.1.secret_name", "bar1foo-ssl"),
				),
			},
		},
	})
}

func testAccCheckKubernetesIngressDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*kubernetes.Clientset)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_ingress" {
			continue
		}
		namespace, name := idParts(rs.Primary.ID)
		resp, err := conn.ExtensionsV1beta1().Ingresses(namespace).Get(name)
		if err == nil {
			if resp.Name == name {
				return fmt.Errorf("Ingress still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil

}

func testAccCheckKubernetesIngressExists(n string, obj *v1beta1.Ingress) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)
		namespace, name := idParts(rs.Primary.ID)
		out, err := conn.ExtensionsV1beta1().Ingresses(namespace).Get(name)
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesIngressConfig_basic() string {
	return fmt.Sprintf(`
resource "kubernetes_ingress" "ingress" {
  metadata {
    labels {
      provider = "kubernetes"
    }
    name = "echo-map"
  }
  spec {
    rules =[
      {
        host = "foo.bar.com"
        http {
          paths = [
            {
              path = "/foo"
              backend {
                service_name = "echoheaders-x"
                service_port = 80
              }
            }
          ]
        }
      },
      {
        host= "bar.baz.com"
        http {
          paths = [
            {
              path = "/bar"
              backend {
                service_name = "echoheaders-y"
                service_port = 80
              }
            },
            {
              path = "/foo"
              backend {
                service_name = "echoheaders-x"
                service_port = 80
              }
            }
          ]
        }
      }
    ]
  }
}
	`)
}

func testAccKubernetesIngressConfig_update() string {
	return fmt.Sprintf(`
resource "kubernetes_ingress" "ingress" {
  metadata {
    labels {
      provider = "kubernetes"
	  app = "echoMap"
    }
    name = "echo-map"
  }
  spec {
    rules =[
      {
        host = "foo1.bar.com"
        http {
          paths = [
            {
              path = "/foo1"
              backend {
                service_name = "echoheaders-x1"
                service_port = 90
              }
            }
          ]
        }
      },
      {
        host= "bar1.baz.com"
        http {
          paths = [
            {
              path = "/bar1"
              backend {
                service_name = "echoheaders-y1"
                service_port = 90
              }
            }
          ]
        }
      }
    ]
  }
}
	`)
}

func testAccKubernetesIngressConfig_tlsSpecified() string {
	return fmt.Sprintf(`
resource "kubernetes_ingress" "ingress" {
  metadata {
    labels {
      provider = "kubernetes"
	  app = "echoMap"
    }
    name = "echo-map"
  }
  spec {
    rules =[
      {
        host = "foo1.bar.com"
        http {
          paths = [
            {
              path = "/foo1"
              backend {
                service_name = "echoheaders-x1"
                service_port = 90
              }
            }
          ]
        }
      },
      {
        host= "bar1.baz.com"
        http {
          paths = [
            {
              path = "/bar1"
              backend {
                service_name = "echoheaders-y1"
                service_port = 90
              }
            }
          ]
        }
      }
    ]
	
	tls = [
		{
			 hosts = ["foo1.bar.com"]
			 secret_name = "foobar1-ssl"
		},		
		{
			 hosts = ["bar1.bar.com"]
			 secret_name = "bar1foo-ssl"
		}
	
	]
  }
}
	`)
}
