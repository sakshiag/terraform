provider "ibmcloud" {
}

resource "ibmcloud_infra_ssh_key" "ssh_key_gip" {
    label = "${var.ssh_label}"
    public_key = "${file(var.ssh_key_path)}"
}

resource "ibmcloud_infra_virtual_guest" "vm1" {
    hostname = "terraform-ibm"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "${var.datacenter}"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
    ssh_key_ids = [
        "${ibmcloud_infra_ssh_key.ssh_key_gip.id}"
    ]
    provisioner "remote-exec" {
        script = "gip.sh"
    }
}

resource "ibmcloud_infra_global_ip" "test-global-ip" {
    routes_to = "${ibmcloud_infra_virtual_guest.vm1.ipv4_address}"
}

resource "ibmcloud_infra_fw_hardware_dedicated" "accfw" {
  ha_enabled = false
  public_vlan_id = "${ibmcloud_infra_virtual_guest.vm1.public_vlan_id}"
}

resource "ibmcloud_infra_fw_hardware_dedicated_rules" "rules" {
 firewall_id = "${ibmcloud_infra_fw_hardware_dedicated.accfw.id}"
 rules = {
      "action" = "deny"
      "src_ip_address"= "0.0.0.0"
      "src_ip_cidr"= 0
      "dst_ip_address"= "any"
      "dst_ip_cidr"= 32
      "dst_port_range_start"= 1
      "dst_port_range_end"= 65535
      "notes"= "Deny all"
      "protocol"= "tcp"
 }
 rules = {
      "action" = "permit"
      "src_ip_address"= "0.0.0.0"
      "src_ip_cidr"= 0
      "dst_ip_address"= "any"
      "dst_ip_cidr"= 32
      "dst_port_range_start"= 22
      "dst_port_range_end"= 22
      "notes"= "Allow SSH"
      "protocol"= "tcp"
 }
}
