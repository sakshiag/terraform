---
layout: "ibmcloud"
page_title: "IBM Cloud: infra_ssh_key"
sidebar_current: "docs-ibmcloud-resource-infra-ssh-key"
description: |-
  Manages IBM Cloud infrastructure SSH keys.
---

# ibmcloud\_infra_ssh_key

Provides SSH keys. This allows SSH keys to be created, updated, and deleted.

For additional details, see the [Bluemix Infrastructure (SoftLayer) API docs](http://sldn.softlayer.com/reference/datatypes/SoftLayer_Security_Ssh_Key).

## Example Usage

```
resource "ibmcloud_infra_ssh_key" "test_ssh_key" {
    label = "test_ssh_key_name"
    notes = "test_ssh_key_notes"
    public_key = "ssh-rsa <rsa_public_key>"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) A descriptive name used to identify an SSH key.
* `public_key` - (Required) The public SSH key.
* `notes` - (Optional) Descriptive text about an SSH key to use at your discretion.

The `label` and `notes` fields are editable.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the new SSH key.
* `fingerprint` - Sequence of bytes to authenticate or look up a longer SSH key.
