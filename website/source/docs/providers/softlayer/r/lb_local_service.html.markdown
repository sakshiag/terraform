---
layout: "softlayer"
page_title: "SoftLayer: softlayer_lb_local_service"
sidebar_current: "docs-softlayer-resource-lb-local-service"
description: |-
  Provides Softlayer's Local Load Balancer Service
---

# softlayer_lb_local_service

Create, update, and destroy SoftLayer Local Load Balancer Service's.

## Example Usage | [SLDN](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_Service)

```
resource "softlayer_lb_local_service" "test_service" {
    port = 80
    service_group_id = "${softlayer_lb_local_service_group.test_service_group.id}"
    weight = 1
    health_check_type = "DNS"
    ip_address_id = "${softlayer_virtual_guest.test_server_1.ip_address_id}"
}
```

## Argument Reference

* `port` | *int*
    * (Required) Specifies the destination port for the Local Load Balancer Service.
* `service_group_id` | *string*
    * (Required) The id of the Local Load Balancer Service Group that this Service will be associated with.
* `weight` | *int*
    * (Required) The weight of the Local Load Balancer Service Group.
* `health_check_type` | *string*
    * (Required) The health check type of the Local Load Balancer Service. Accepted values are `DEFAULT`,
        `DNS`, `HTTP`, `HTTP-CUSTOM`, `ICMP`, and `TCP`.
* `ip_address_id` | *int*
    * (Required) The IP Address Id of the destination virtual guest.

## Attributes Reference

* `id` - A Local Load Balancer Service's internal identifier.
