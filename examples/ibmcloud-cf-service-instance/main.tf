provider "ibmcloud" {}

data "ibmcloud_cf_space" "spacedata" {
  space = "${var.space}"
  org   = "${var.org}"
}

resource "ibmcloud_cf_service_instance" "service-instance" {
  name       = "${var.service_instance_name}"
  space_guid = "${data.ibmcloud_cf_space.spacedata.id}"
  service    = "${var.service}"
  plan       = "${var.plan}"
  tags       = ["cluster-service", "cluster-bind"]
}

resource "ibmcloud_cf_service_key" "serviceKey" {
  name                  = "${var.service_key_name}"
  service_instance_guid = "${ibmcloud_cf_service_instance.service-instance.id}"
}
