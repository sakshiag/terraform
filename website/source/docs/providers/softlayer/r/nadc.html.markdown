---
layout: "softlayer"
page_title: "SoftLayer: softlayer_network_application_delivery_controller"
sidebar_current: "docs-softlayer-resource-softlayer-network_application_delivery_controller"
description: |-
  Provides Softlayer's Network Application Delivery Controller
-------------------------------------------

#softlayer_network_application_delivery_controller_resource

Create, update, and delete Softlayer Network Application Delivery Controllers. For additional details please refer to the [API documentation](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Application_Delivery_Controller).
## Example Usage

```
resource "softlayer_network_application_delivery_controller" "test_nadc" {
    datacenter = "DALLAS06"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
}
```

## Argument Reference

* `datacenter` | *string*
    * (Required) Specifies which datacenter the Network Application Delivery Controller is to be provisioned in. Accepted values can be found [here](http://www.softlayer.com/data-centers).
* `speed` | *int*
    * (Required) The speed in Mbps. Accepted values are `10`, `200`, and `1000`.
* `version` | *string*
    * (Required) The Network Application Delivery Controller version. Accepted values are `10.1` and `10.5`.
* `plan` | *string*
    * (Required) The Network Application Delivery Controller plan. Accepted values are `Standard` and `Platinum`.
* `ip_count` | *int*
    * (Required) The number of static public IP addresses assigned to the Network Application Delivery Controller. Accepted values are `2`, `4`, `8`, and `16`.

## Attributes Reference

* `id` - A Network Application Delivery Controller's internal identifier.
* `name` - A Network Application Delivery Controller's internal name.
