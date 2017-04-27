---
layout: "ibmcloud"
page_title: "IBM Cloud: objectstorage_account"
sidebar_current: "docs-ibmcloud-resource-infra-objectstorage-account"
description: |-
  Manages IBM Cloud infrastructure object storage account.
---

# ibmcloud\_infra_objectstorage_account

Retrieve the account name for an existing Object Storage instance within your IBM Cloud account. If there is no Object Storage instance, you can use this resources to order one for you and remember the account name. 

This resource is not intended for managing the lifecycle (e.g. update, delete) of an Object Storage instance in IBM Cloud. For lifecycle management, see the Swift API or Swift resources. 

## Example Usage

```hcl
resource "ibmcloud_infra_objectstorage_account" "foo" {
}
```

## Argument Reference

No additional arguments needed.

## Computed Fields

The following attributes are exported:

* `id` - The Object Storage account name, which you can use with Swift resources.
