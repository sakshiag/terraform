---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_service_key"
sidebar_current: "docs-ibmcloud-resource-cf-service-key"
description: |-
  Manages IBM Cloud Cloud Foundry service key.
---

# ibmcloud\_cf_service_key

Create, update, or delete CF service keys for IBM Bluemix.

## Example Usage

```hcl
data "ibmcloud_cf_service_instance" "service_instance" {
  name = "mycloudant"
}

resource "ibmcloud_cf_service_key" "serviceKey" {
  name                  = "mycloudantkey"
  service_instance_guid = "${data.ibmcloud_cf_service_instance.service_instance.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A descriptive name used to identify a service key.
* `parameters` - (Optional, map) Arbitrary parameters to pass along to the service broker. Must be a JSON object.
* `service_instance_guid` - (Required, string) The GUID of the service instance that the service key needs to be associated with.



## Attributes Reference

The following attributes are exported:

* `credentials` - Credentials associated with the key.
