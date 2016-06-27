---
layout: "softlayer"
page_title: "SoftLayer: softlayer_lb_local"
sidebar_current: "docs-softlayer-resource-lb-local"
description: |-
  Provides Softlayer's Local Load Balancer
---

# softlayer_lb_local

Create, update, and destroy SoftLayer Local Load Balancers.

_Please Note_: SoftLayer Local Load Balancer's are currently priced on a per-month basis, so please use caution when creating the resource as the cost for an entire month is incurred immediately upon creation. For more information on pricing please see this [link](http://www.softlayer.com/load-balancing).

You can also use this REST URL to get a listing of local load balancer choices:

```
https://{{userName}}:{{apiKey}}@api.softlayer.com/rest/v3/SoftLayer_Product_Package/194/getItems.json?objectMask=id;capacity;description;units;keyName;prices.id;prices.categories.id;prices.categories.name
```


## Example Usage | [SLDN](http://sldn.softlayer.com/reference/services/SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress)

```
resource "softlayer_lb_local" "test_lb_local" {
    connections = 15000
    location    = "tok02"
    ha_enabled  = false
}
```

## Argument Reference

* `connections` | *int*
    * (Required) Specifies the max connections supported by the Local Load Balancer. Accepted values are `150000` and `15000`.
* `location` | *string*
    * (Required) The datacenter location.
* `ha_enabled` | *boolean*
    * (Required) Denotes whether or not the Local Load Balancer will be configured within a high availability cluster.
* `security_certificate_id` | *int*
    * (Optional) The id of the Security Certificate to be used.

## Attributes Reference

* `id` - A Local Load Balancer's internal identifier.
* `subnet_id` - A Local Load Balancer's subnet id.
* `ip_address` - A local Load Balancer's ip address.
