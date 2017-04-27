---
layout: "ibmcloud"
page_title: "IBM Cloud: infra_provisioning_hook"
sidebar_current: "docs-ibmcloud-resource-infra-provisioning-hook"
description: |-
  Manages IBM Cloud infrastructure provisioning hook.
---


# ibmcloud\_infra_provisioning_hook

Provides provisioning hooks containing all the information needed to add a hook into a server or virtual provision and OS reload. This allows provisioning hooks to be created, updated, and deleted.

For additional details, see the [Bluemix Infrastructure (SoftLayer) API docs](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Provisioning_Hook).

## Example Usage

```hcl
resource "ibmcloud_infra_provisioning_hook" "test_provisioning_hook" {
    name = "test_provisioning_hook_name"
    uri = "https://raw.githubusercontent.com/test/slvm/master/test-script.sh"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A descriptive name used to identify a provisioning hook.
* `uri` - (Required, string) The endpoint that the script is downloaded or downloaded and executed from.

The `name` and `uri` field are editable.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the new provisioning hook.
