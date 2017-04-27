---
layout: "ibmcloud"
page_title: "IBM Cloud: global_ip"
sidebar_current: "docs-ibmcloud-resource-global-ip"
description: |-
  Manages IBM Cloud Infrastructure global IP.
---

# ibmcloud\_infra_global_ip

Provides a global IP resource to route between servers. This allows global IPs to be created, updated, and deleted. Global IPs are not restricted to routing within the same data center.

For additional details, see the [Bluemix Infrastructure (SoftLayer) API docs](http://sldn.softlayer.com/reference/services/SoftLayer_Network_Subnet_IpAddress_Global) and [global IP address overview](https://knowledgelayer.softlayer.com/learning/global-ip-addresses).

## Example Usage

```hcl
resource "ibmcloud_infra_global_ip" "test_global_ip " {
    routes_to = "119.81.82.163"
}
```

## Argument Reference

The following arguments are supported:

* `routes_to` - (Required, string) Destination IP address that the global IP routes traffic through. The destination IP address can be a public IP address of IBM Cloud resources in the same account, such as a public IP address of virtual guests and public virtual IP address of NetScaler VPXs. 

## Attributes Reference

The following attributes are exported:

* `id` - ID of the global IP.
* `ip_address` - Address of the global IP.
