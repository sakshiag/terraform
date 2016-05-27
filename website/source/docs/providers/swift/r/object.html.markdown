---
layout: "swift"
page_title: "Swift: object"
sidebar_current: "docs-swift-resource-object"
description: |-
  For managing objects in a Swift object store.
---

# swift_object

Provides an object resource. Allows the creation of an object in a Swift object store.

## Example Usage

```
# Create a new object in Swift with contents extracted from a local file
resource "swift_object" "test_object" {
    name = "foo.txt" # Object name
    container_name = "${swift_container.test_container_1.name}"
    source_file = "foo.txt" # Where to read the contents of the new object from
}
```

```
# Create a new object in Swift with contents specified as a variable.
# NOTE: the path specified will automatically be created
variable "secrets" {
    type = "string"
}

resource "swift_object" "test_object2" {
    name = "path/bar.txt" # Path name/Object name
    container_name = "${swift_container.test_container_1.name}"
    contents = "${var.secrets}"
}
```

## Argument Reference

The following arguments are supported:

* `name` | *string*
	* Name of the object. This name can also have forward slashes, which will act as a pseudo file path identifying the object location within the container (e.g. _path/to/foo.txt_).
	* **Required**
* `container_name` | *string*
	* Name of the container to put this object in.
	* **Required**
* `source_file` | *string*
	* The source file containing the desired contents of the object.
	* **Optional**
* `contents` | *string*
	* The desired contents of the object. If *source_file* is specified, this argument will be ignored.
	* **Optional**

If neither `source_file` nor `contents` are specified, an empty object will be created.
