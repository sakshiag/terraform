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

func TestAccKubernetesPod_with_container_liveness_probe_using_exec(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/busybox"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithLivenessProbeUsingExec(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.args.#", "3"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.exec.0.command.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.exec.0.command.0", "cat"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.exec.0.command.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.failure_threshold", "3"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.initial_delay_seconds", "5"),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_container_liveness_probe_using_http_get(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithLivenessProbeUsingHTTPGet(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.args.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.0.path", "/healthz"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.0.http_headers.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.0.http_headers.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.http_get.0.http_headers.0.value", "Awesome"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.initial_delay_seconds", "3"),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_container_liveness_probe_using_tcp(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithLivenessProbeUsingTCP(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.args.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.tcp_socket.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.liveness_probe.0.tcp_socket.0.port", "8080"),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_container_lifecycle(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithLifeCycle(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.post_start.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.post_start.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.post_start.0.exec.0.command.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.post_start.0.exec.0.command.0", "ls"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.post_start.0.exec.0.command.1", "-al"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.pre_stop.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.pre_stop.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.pre_stop.0.exec.0.command.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.lifecycle.0.pre_stop.0.exec.0.command.0", "date"),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_container_security_context(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithContainerSecurityContext(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.security_context.#", "1"),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_volume_mounts(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	secretName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithVolumeMounts(secretName, podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.volume_mounts.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.volume_mounts.0.mount_path", "/tmp/my_path"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.volume_mounts.0.name", "db"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.volume_mounts.0.read_only", "false"),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.volume_mounts.0.sub_path", ""),
				),
			},
		},
	})
}

func TestAccKubernetesPod_with_resource_requirements(t *testing.T) {
	var conf api.Pod

	podName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPodConfigWithResourceRequirements(podName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesPodExists("kubernetes_pod.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_pod.test", "spec.0.containers.0.image", imageName),
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
    security_context {
      run_as_non_root     = true
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

func testAccKubernetesPodConfigWithLivenessProbeUsingExec(podName, imageName string) string {
	return fmt.Sprintf(`
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
      args  = ["/bin/sh", "-c", "touch /tmp/healthy; sleep 300; rm -rf /tmp/healthy; sleep 600"]

      liveness_probe {
        exec {
          command = ["cat", "/tmp/healthy"]
        }

        initial_delay_seconds = 5
        period_seconds        = 5
      }
    }
  }
}
	`, podName, imageName)
}

func testAccKubernetesPodConfigWithLivenessProbeUsingHTTPGet(podName, imageName string) string {
	return fmt.Sprintf(`
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
      args  = ["/server"]

      liveness_probe {
        http_get {
          path = "/healthz"
          port = 8080

          http_headers {
            name  = "X-Custom-Header"
            value = "Awesome"
          }
        }
        initial_delay_seconds = 3
        period_seconds        = 3
      }
    }
  }
}
	`, podName, imageName)
}

func testAccKubernetesPodConfigWithLivenessProbeUsingTCP(podName, imageName string) string {
	return fmt.Sprintf(`
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
      args  = ["/server"]

      liveness_probe {
        tcp_socket {
          port = 8080
        }

        initial_delay_seconds = 3
        period_seconds        = 3
      }
    }
  }
}
	`, podName, imageName)
}

func testAccKubernetesPodConfigWithLifeCycle(podName, imageName string) string {
	return fmt.Sprintf(`
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
      args  = ["/server"]

      lifecycle {
        post_start {
          exec {
            command = ["ls", "-al"]
          }
        }
        pre_stop {
          exec {
            command = ["date"]
          }
        }
      }
    }
  }
}

	`, podName, imageName)
}

func testAccKubernetesPodConfigWithContainerSecurityContext(podName, imageName string) string {
	return fmt.Sprintf(`
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

      security_context {
        privileged = true

        se_linux_options {
          level = "s0:c123,c456"
        }
      }
    }
  }
}


	`, podName, imageName)
}

func testAccKubernetesPodConfigWithVolumeMounts(secretName, podName, imageName string) string {
	return fmt.Sprintf(`

resource "kubernetes_secret" "test" {
  metadata {
    name = "%s"
  }

  data {
    one = "first"
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
			volume_mounts = [{
				mount_path = "/tmp/my_path"
				name  = "db"
			}]
    }
    volumes = [{
      name = "db"
      secret = {
        secret_name = "${kubernetes_secret.test.metadata.0.name}"
      }
    }]
  }
}
	`, secretName, podName, imageName)
}

func testAccKubernetesPodConfigWithResourceRequirements(podName, imageName string) string {
	return fmt.Sprintf(`

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
			//TODO
			/*
			resources{
				limits{
					cpu = "1"
					memory = "512Mi"
				}
				requests{
					cpu = "250m"
					memory = "50Mi"
				}
			}
			*/
    }
  }
}
	`, podName, imageName)
}
