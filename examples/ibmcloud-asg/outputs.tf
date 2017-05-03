#ip_address - cluster address
output "cluster_address" {
  value = "http://${ibmcloud_infra_lb_local.local_lb.ip_address}"
}
