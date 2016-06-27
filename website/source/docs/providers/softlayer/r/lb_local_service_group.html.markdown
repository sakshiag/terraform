---
layout: "softlayer"
page_title: "SoftLayer: softlayer_lb_local_service_group"
sidebar_current: "docs-softlayer-resource-lb-local-service-group"
description: |-
  Provides Softlayer's Local Load Balancer Service Group
---

# softlayer_lb_local_service_group

Create, update, and destroy SoftLayer Local Load Balancer Service Groups.

## Example Usage | [SLDN](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_Service_Group)

```
resource "softlayer_lb_local_service_group" "test_service_group" {
    port = 82
    routing_method = "CONSISTENT_HASH_IP"
    routing_type = "HTTP"
    load_balancer_id = "${softlayer_lb_local.test_lb.id}"
    allocation = 100
}
```

## Argument Reference

* `port` | *int*
    * (Required) Specifies the listening port for the Local Load Balancer Service Group.
* `routing_method` | *string*
    * (Required) The routing method for the Local Load Balancer Service Group. Accepted values: `CONSISTENT_HASH_IP`,
    `INSERT_COOKIE`, `LEAST_CONNECTIONS`, `LEAST_CONNECTIONS_INSERT_COOKIE`, `LEAST_CONNECTIONS_PERSISTENT_IP`,
    `PERSISTENT_IP`, `ROUND_ROBIN`, `ROUND_ROBIN_INSERT_COOKIE`, `ROUND_ROBIN_PERSISTENT_IP`, `SHORTEST_RESPONSE`,
    `SHORTEST_RESPONSE_INSERT_COOKIE`, `SHORTEST_RESPONSE_PERSISTENT_IP`.
* `routing_type` | *string*
    * (Required) The routing method for the Local Load Balancer Service Group. Accepted values: `DNS`,
    `FTP`, `HTTP`, `HTTPS`, `TCP`, `UDP`.
* `load_balancer_id` | *string*
    * (Required) The id of the Local Load Balancer the Service Group will be associated with.
* `allocation` | *int*
    * (Required) The allocation for the Local Load Balancer the Service Group.

## Attributes Reference

* `id` - A Local Load Balancer Service Group's internal identifier.
* `virtual_server_id` - A Local Load Balancer Service Group's associated Virtual Server identifier. Note that
    the implementation details of the virtual server are handled internally by Softlayer.
