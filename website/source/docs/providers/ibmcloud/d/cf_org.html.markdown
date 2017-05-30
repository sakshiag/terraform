---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_org"
sidebar_current: "docs-ibmcloud-datasource-cf-org"
description: |-
  Get information about an IBM Bluemix organization.
---

# ibmcloud\_cf_org

Import the details of an existing IBM Bluemix org as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl
data "ibmcloud_cf_org" "orgdata" {
  org = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `org` - (Required) The name of the Bluemix org.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the org.  
