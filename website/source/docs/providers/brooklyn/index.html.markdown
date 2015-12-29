---
layout: "brooklyn"
page_title: "Provider: Apache Brooklyn"
sidebar_current: "docs-brooklyn-index"
description: |-
  The Brooklyn provider is used to interact with Apache Brooklyn application.
---

# Apache Brooklyn Provider

The SoftLayer provider is used to manage SoftLayer resources.

Use the navigation to the left to read about the available resources.

<div class="alert alert-block alert-info">
<strong>Note:</strong> The Brooklyn provider is new as of Terraform 0.X.
It is ready to be used but many features are still being added. If there
is a SoftLayer feature missing, please report it in the GitHub repo.
</div>

## Example Usage

Here is an example that will create an application:

(create this as br.tf and run terraform commands from this directory):

```hcl
provider "brooklyn" {
    username = ""
    api_key = ""
}

# This will create a new application in Apache Brooklyn using YAML template that will show up under the Applications tab in the Apache Brooklyn console.
resource "brooklyn_application" "application1" {
    application_spec = "/home/user/brooklyn_test.yaml"
}
```

You'll need to provide your Brooklyn URL, username and password,
so that Terraform can connect. If you don't want to put
credentials in your configuration file, you can leave them
out:

```
provider "brooklyn" {}
```

...and instead set these environment variables:

- **BROOKLYN_URL**: Your Brooklyn application URL.
- **BROOKLYN_USERNAME**: Your Brooklyn username.
- **BROOKLYN_PASSWORD**: Your Brooklyn password.

