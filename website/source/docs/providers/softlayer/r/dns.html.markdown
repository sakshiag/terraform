
Provides a Softlayer's DNS Domain
---

# SoftLayer_Dns_Domain

The SoftLayer_Dns_Domain data type represents a single DNS domain record hosted on the SoftLayer nameservers.
Domains contain general information about the domain name such as name and serial. Individual records such as _A_, _AAAA_,
_CTYPE_, and _MX_ records are stored in the domain's associated [`SoftLayer_Dns_Domain_ResourceRecord`](/docs/providers/softlayer/r/dns_records.html) records.

## Example Usage

###Create example

```
resource "softlayer_dns_domain" "dns-domain-test" {
	name = "test_domain_qwerqq.com"
	records = {
		record_data = "127.0.0.1"
		domain_id = 1857408
		host = "hostazzzz.com"
		contact_email = "user@softlaer.com"
		ttl = 900
		record_type = "a"
	}
}
```

Create Dns domain. _Name_ is mandatory field. Currently only one record may be specified in the records field due to Terraform limitations.
Other records may be linked via entity id of the created domain. You can omit the records field at all:

```
resource "softlayer_dns_domain" "dns-domain-test" {
	name = "test_domain_qwerqq.com"
}
```

During creation two Dns Domain Records are created automatically, of NS and SOA type, which breaks the **plan** function of terraform in case you specify the records fields. 


###Update example

```
resource "softlayer_dns_domain" "dns-domain-test" {
	name = "test_domain_qwerqq_changed.com"
}
```

```
resource "softlayer_dns_domain_record" "recordA" {
    record_data = "127.0.0.1"
    domain_id = "${softlayer_dns_domain.dns-domain-test.id}"
    host = "hosta.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "a"
}
```

```
resource "softlayer_dns_domain_record" "recordAAAA" {
    record_data = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329"
    domain_id = "${softlayer_dns_domain.dns-domain-test.id}"
    host = "hosta-2.com"
    contact_email = "user2@softlaer.com"
    ttl = 1000
    record_type = "aaaa"
}
```

Name field may be changed, but it will **force** the ```delete->create``` process for the _dns_doamin_ and all records,
linked to it via _domain_id_ field. All these changes are reflected in the terraform **plan** command output. 

###Delete example

Simple removal of the resource from the .tf file will cause the _dns_domain_ to be deleted.

## Argument Reference

The following arguments are supported:

* `name` - (Required) A domain's name including top-level domain, for example "example.com".

## Attributes Reference

* `id` - A domain record's internal identifier.
* `serial` - A unique number denoting the latest revision of a domain.
* `update_date` - The date that this domain record was last updated.