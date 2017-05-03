output "global ip"{
	value = "http://${ibmcloud_infra_global_ip.test-global-ip.ip_address}"
}