---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_app"
sidebar_current: "docs-ibmcloud-resource-cf-app"
description: |-
  Manages IBM Cloud Cloud Foundry app.
---

# ibmcloud\_cf_app

Create, update, or delete CF app on IBM Bluemix.

## Example Usage

```hcl
	
data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "hello.zip"
	ports = [9080]
	wait_timeout = 90
	buildpack = "sdk-for-nodejs"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The GUID of the associated domain. The values can be retrieved from data source `ibmcloud_cf_domain`.
* `memory` - (Optional, int) The amount of memory each instance should have. In megabytes.
* `instances` - (Optional, int) The number of instances.
* `disk_quota` - (Optional, int) The maximum amount of disk available to an instance of an app. In megabytes.
* `space_guid` - (Required, string) Define space guid to which app belongs. The values can be retrieved from data source `ibmcloud_cf_space`.
* `command` - (Optional, string) The initial command for the app.
* `buildpack` - (Optional, string) Buildpack to build the app. 3 options: a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack. 
* `diego` - (Optional, bool) Use diego to stage and to run when available.
* `docker_image` - (Optional, string) Name of the Docker image containing the app.
* `docker_credentials_json` - (Optional, map) Docker credentials for pulling docker image.
* `environment_json` - (Optional, map) Key/value pairs of all the environment variables to run in your app. Does not include any system or service variables.
* `ports` - (Optional, list) Ports on which application may listen. Overwrites previously configured ports. Ports must be in range 1024-65535. Supported for Diego only.			
* `route_guid` - (Optional, set) Define the route guid needs to be attached to application.
* `service_instance_guid` - (Optional, set) Define the service guid needs to be attached to application.
* `wait_timeout` - (Optional, int) Define timeout to wait for the app to start.
* `app_path` - (Optional, string) Define the path of the zip file of the application

## Attributes Reference6

The following attributes are exported:

* `id` - The ID of the app.

