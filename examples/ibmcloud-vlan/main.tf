provider "ibmcloud" {
}

# Create a public vlan
resource "ibmcloud_infra_vlan" "test_vlan_public" {
  name            = "${var.vlan_name_public}"
  datacenter      = "${var.datacenter}"
  type            = "PUBLIC"
  subnet_size     = 8
}

# Create a private vlan
resource "ibmcloud_infra_vlan" "test_vlan_private" {
  name        = "${var.vlan_name_private}"
  datacenter  = "${var.datacenter}"
  type        = "PRIVATE"
  subnet_size = 8
}

# Create a new ssh key 
resource "ibmcloud_infra_ssh_key" "ssh_key" {
    label = "${var.ssh_label}"
    notes = "for public vlan test"
    public_key = "${file(var.ssh_public_key)}"
}

# Create a new virtual guest using image "CENTOS_7_64"
resource "ibmcloud_infra_virtual_guest" "vm1" {
    hostname = "${var.vm_hostname}"
    os_reference_code = "${var.vm_os_reference_code}"
    domain = "${var.vm_domain}"
    datacenter = "${var.datacenter}"
    network_speed = "${var.vm_network_speed}"
    hourly_billing = true
    private_network_only = false
    cores = "${var.vm_cores}"
    memory = "${var.vm_memory}"
    disks = "${var.vm_disks}"
    user_metadata = "{\"value\":\"newvalue\"}"
    dedicated_acct_host_only = true
    local_disk = false
    ssh_key_ids = ["${ibmcloud_infra_ssh_key.ssh_key.id}"]
    public_vlan_id  = "${ibmcloud_infra_vlan.test_vlan_public.id}"
    private_vlan_id = "${ibmcloud_infra_vlan.test_vlan_private.id}"
    
}
