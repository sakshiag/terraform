

  Provides a Softlayer's DNS Domain Records
---

# SoftLayer_Dns_Domain_ResourceRecord

The _SoftLayer_Dns_Domain_ResourceRecord_ data type represents a single resource record entry in a SoftLayer hosted [`domain`](/docs/providers/softlayer/r/dns.html).
Each resource record contains a _host_ and _data_ property, defining a resource's name and it's target data.
Domains contain multiple types of resource records. The _type_ property separates out resource records by _type_.
_Type_ can take one of the following values:

* "a" for address records
* "aaaa" for address records
* "cname" for canonical name records
* "mx" for mail exchanger records
* "ns" for name server records
* "ptr" for pointer records in reverse domains
* "soa" for a domain's start of authority record
* "spf" for sender policy framework records
* "srv" for service records
* "txt" for text records

As _SoftLayer_Dns_Domain_ResourceRecord_ objects are created and loaded, the API verifies the _type_ property and casts
the object as the appropriate type.

## Example Usage

We are using [SoftLayer_Dns_Domain_ResourceRecord](https://sldn.softlayer.com/reference/datatypes/SoftLayer_Dns_Domain_ResourceRecord)
SL’s object for most of CRUD operations. Only for SRV record type we are using [SoftLayer_Dns_Domain_ResourceRecord_SrvType](https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord_SrvType) SL’s object.

Currently we can CRUD almost all record types except _SOA_ type which is initially created on DNS create action. 

###Create example:

```
provider "softlayer" {
    username = "your_user_name"
    api_key = "your_api_key"
}
```

```
# Create a DNS for records creation
resource "softlayer_dns_domain" "test_dns_domain_records" {
	name = "records.test.domain-3.com"
}
```

```
# Creates a record of A _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordA" {
    record_data = "127.0.0.1"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "a"
}
```

```
# Creates a record of AAAA _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordAAAA" {
    record_data = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-2.com"
    contact_email = "user2changed@softlaer.com"
    ttl = 1000
    record_type = "aaaa"
}
```

```
# Creates a record of CNAME _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordCNAME" {
    record_data = "testcname.com"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-cname.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "cname"
}
```

```
# Creates a record of MX _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordMX" {
    record_data = "testmx.com"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-mx.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "mx"
}
```

```
# Creates a record of NS _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordNS" {
    record_data = "ns1.example.org"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-ns.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "ns"
}
```

```
# Creates a record of SPF _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordSPF" {
    record_data = "v=spf1 mx:mail.example.org ~all"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-spf"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "spf"
}  
```

```
# Creates a record of TXT _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordTXT" {
    record_data = "127.0.0.1"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-txt.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "txt"
}
```

```
# Creates a record of SRV _type_ using DNS id created before
resource "softlayer_dns_domain_record" "recordSRV" {
    record_data = "ns1.example.org"
    domain_id = "${softlayer_dns_domain.test_dns_domain_records.id}"
    host = "hosta-srv.com"
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "srv"
	port = 8080
	priority = 3
	protocol = "_tcp"
	weight = 3
	service = "_mail"
}
```

####Create for PTR record type example: 

The issue we faced was to determine a _domainId_ to use on this type of record creation. If a user will use _domainId_ 
which is not in any available subnets he get an error that record cannot be found. For tests we made following steps:
Created a PTR record via UI (Network->DNS->Reverse Records) assigning IPAddress in the range of existing and available subnets (Network->IP Management->Subnets).
This IPAddress we will use on PTR record creation via terraform. 
Listed Subnets via curl, rest or python client to retrieve subnetId which we used at first step
Get Subnet’s details (not python client as it doesn’t provide info about reverse records) via rest tool 
(endpoint: /rest/v3.1/SoftLayer_Network_Subnet/subnetId/getReverseDomainRecords).
As a result we got: "resourceRecords": [{"data": "some.host.name.com.","domainId": 1653916,"expire": null,...} 

So having _IPAddress_ and _domainId_ we can create valid and available PRT record

```
resource "softlayer_dns_domain_record" "recordPTR" {
    record_data = "ptr.domain.com"
    domain_id = 1653916
    host = "45"  ?============= this is the last octet of IPAddress in the range of the subnet
    contact_email = "user@softlaer.com"
    ttl = 900
    record_type = "ptr"
}
```

##Delete example: (removes all records from DNS in the configuration) 

Simple removal of all existing resources from the .tf file will cause them to be deleted.

##Edit example:

You just need to edit any editable field and apply. Note that some fields cause to re-create the resource.
For two properties _record_data_ and _record_type_ you need to update both props with correct values if one of it need to be changed.
Otherwise an error will be retrieved from SL.

## Argument Reference

* `record_data` - (Required) The value of a domain's resource record. This can be an IP address or a hostname. Fully qualified host and domain name data must end with the "." character.
* `domain_id` - (Required) An identifier belonging to the domain that a resource record is associated with.
* `expire` - The amount of time in seconds that a secondary name server (or servers) will hold a zone before it is no longer considered authoritative.
* `host` - (Required)
* `minimum_ttl` - The amount of time in seconds that a domain's resource records are valid. This is also known as a minimum TTL, and can be overridden by an individual resource record's TTL.
* `mx_priority` - Useful in cases where a domain has more than one mail exchanger, the priority property is the priority of the MTA that delivers mail for a domain. A lower number denotes a higher priority, 
and mail will attempt to deliver through that MTA before moving to lower priority mail servers. Priority is defaulted to 10 upon resource record creation.
* `refresh` - The amount of time in seconds that a secondary name server should wait to check for a new copy of a DNS zone from the domain's primary name server. 
If a zone file has changed then the secondary DNS server will update it's copy of the zone to match the primary DNS server's zone.
* `contact_email` - (Required) The email address of the person responsible for a domain, with the "@" replaced with a ".".
 For instance, if root@example.org is responsible for example.org, then example.org's SOA responsibility is "root.example.org.".
* `retry` - The amount of time in seconds that a domain's primary name server (or servers) should wait if an attempt to refresh
 by a secondary name server failed before attempting to refresh a domain's zone with that secondary name server again.
* `ttl` - (Required) The Time To Live value of a resource record, measured in seconds.
 TTL is used by a name server to determine how long to cache a resource record. An SOA record's TTL value defines the domain's overall TTL.
* `record_type` - (Required) A domain resource record's type
* `service` - The symbolic name of the desired service
* `protocol` - The protocol of the desired service; this is usually either TCP or UDP.
* `port` - The TCP or UDP port on which the service is to be found.
* `priority` - The priority of the target host, lower value means more preferred.
* `weight` - A relative weight for records with the same priority.

## Attributes Reference

* `id` - A domain resource record's internal identifier.