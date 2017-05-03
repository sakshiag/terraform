provider "ibmcloud" {}

# Create a new ssh key 
resource "ibmcloud_infra_ssh_key" "ssh_key_performance" {
  label      = "${var.ssh-label}"
  notes      = "for scale group"
  public_key = "${file(var.ssh_key_path)}"
}

resource "ibmcloud_infra_lb_local" "local_lb" {
  connections = "${var.lb-connections}"
  datacenter  = "${var.datacenter}"
  ha_enabled  = false
  dedicated   = "${var.lb-dedicated}"
}

resource "ibmcloud_infra_lb_local_service_group" "lb_service_group" {
  port             = "${var.lb-servvice-group-port}"
  routing_method   = "${var.lb-servvice-group-routing-method}"
  routing_type     = "${var.lb-servvice-group-routing-type}"
  load_balancer_id = "${ibmcloud_infra_lb_local.local_lb.id}"
  allocation       = "${var.lb-servvice-group-routing-allocation}"
}

resource "ibmcloud_infra_scale_group" "sample-http-cluster" {
  name                 = "${var.auto-scale-name}"
  regional_group       = "${var.auto-scale-region}"
  cooldown             = "${var.auto-scale-cooldown}"
  minimum_member_count = "${var.auto-scale-minimum-member-count}"
  maximum_member_count = "${var.auto-scale-maximumm-member-count}"
  termination_policy   = "${var.auto-scale-termination-policy}"
  virtual_server_id    = "${ibmcloud_infra_lb_local_service_group.lb_service_group.id}"
  port                 = "${var.auto-scale-lb-service-port}"

  health_check = {
    type = "${var.auto-scale-lb-service-health-check-type}"
  }

  virtual_guest_member_template = {
    hostname                = "${var.vm-hostname}"
    domain                  = "${var.vm-domain}"
    cores                   = "${var.vm-cores}"
    memory                  = "${var.vm-memory}"
    os_reference_code       = "${var.vm-os-reference-code}"
    datacenter              = "${var.datacenter}"
    ssh_key_ids             = ["${ibmcloud_infra_ssh_key.ssh_key_performance.id}"]
    post_install_script_uri = "${var.vm-post-install-script-uri}"
  }
}

resource "ibmcloud_infra_scale_policy" "sample-http-cluster-policy" {
  name           = "${var.scale-policy-name}"
  scale_type     = "${var.scale-policy-type}"
  scale_amount   = "${var.scale-policy-scale-amount}"
  cooldown       = "${var.scale-policy-cooldown}"
  scale_group_id = "${ibmcloud_infra_scale_group.sample-http-cluster.id}"

  triggers = {
    type = "RESOURCE_USE"

    watches = {
      metric   = "host.cpu.percent"
      operator = ">"
      value    = "90"
      period   = 130
    }
  }
}
