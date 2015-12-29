---
layout: "brooklyn"
page_title: "Brooklyn: applications"
sidebar_current: "docs-brooklyn-resource-application"
description: |-
  Manages Apache Brooklyn Applications.
---

# brooklyn\application

Provides Applications. This allows Application to be created, updated and deleted via CAMP YAML specification.
For additional details please refer to [documentation](https://brooklyn.apache.org/learnmore/blueprint-tour.html).

## Example Usage

```
resource "brooklyn_application" "application1" {
    application_spec = "/home/user/brooklyn_test.yaml"
}
```

## Argument Reference

The following arguments are supported:

* `application_spec` - (Required) CAMP YAML application specification.

## Attributes Reference

The following attributes are exported:

* `id` - id of the new application.
* `name` - name of the new application, as defined in YAML specification. 
* `status` - application status (e.g. *RUNNING*).
* `type` - application type (e.g. *org.apache.brooklyn.entity.stock.BasicApplication*).
* `links` - application links. 
* `locations` - application locations. 
