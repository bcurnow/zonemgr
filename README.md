<!-- markdownlint-disable MD007 -->
<!-- markdownlint-disable MD010 -->
<!-- markdownlint-disable MD033 -->
<!-- markdownlint-disable MD041 -->
<!-- vscode-markdown-toc -->
* [Design](#Design)
	* [Zone File Format](#ZoneFileFormat)
		* [Line Types](#LineTypes)
		* [Control Entries](#ControlEntries)
		* [Resource Records](#ResourceRecords)
		* [Classes](#Classes)
		* [Special Values and Escapes](#SpecialValuesandEscapes)
	* [YAML Format](#YAMLFormat)
		* [Examples](#Examples)
	* [Built-In Plugins](#Built-InPlugins)
		* [Plugin Behavior](#PluginBehavior)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->
<!-- markdownlint-enable MD007 -->
<!-- markdownlint-enable MD010 -->
<!-- markdownlint-enable MD041 -->

# zonemgr

A command-line utility for generating BIND zone files from a YAML input

## <a name='Design'></a>Design

### <a name='ZoneFileFormat'></a>Zone File Format

The format of a zone file is largely contained in [RFC1035](https://datatracker.ietf.org/doc/html/rfc1035). Clarification of the 'minimum' value on the SOA record and the introduction of the $TTL line is included in [RFC2308](https://datatracker.ietf.org/doc/html/rfc2308).

#### <a name='LineTypes'></a>Line Types

The zone file is a line oriented file with parentheses used to continue a list of items across a line boundary.

The following basic lines exist:

* \<blank\>[\<comment\>]
* $TTL \<TTL\> [\<comment\>]
* $ORIGIN \<domain name\> [\<comment\>]
* $INCLUDE \<file-name\> [\<domain name\>] [\<comment\>]
* \<domain-name\>\<rr\> [\<comment\>]
* \<blank\>\<rr\> [\<comment\>]

Any number of blank lines are allowed in the file.

#### <a name='ControlEntries'></a>Control Entries

There are several control entries (starting with '$'):

* $TTL - Used to define the TTL for resource records without an explicit TTL
* $ORIGIN - Used to reset the the current origina for relative domain names
* $INCLUDE - Inserts the named file into the current file and optionally includes a domain name that will relative domain name origin for the included file

NOTE: Currently $INCLUDE is not used in Zonemgr

#### <a name='ResourceRecords'></a>Resource Records

Resource records (A, CNAME, MX, etc.) have the following two formats:

* [\<TTL\>] [\<class\>] \<type\> \<RDATA\>
* [\<class\>] [\<TTL\>] \<type\> \<RDATA\>

##### Resource Record Types

The following resource record types are defined by the RFC:

* A - Host address
* NS - An authoritative name server
* MD - Deprecated, use MX
* MF - Deprecated, use MX
* CNAME - Canonical name
* SOA - Start of authority
* MB - Mailbox domain (experimental)
* MG - Mail group  (experimental)
* MR - Mail renname domain (experimental)
* NULL - Null resource record (experimental)
* WKS - Well known service
* PTR - Domain name pointer
* HINFO - Host information
* MINFO - Mailbox or mail list information
* MX  - Mail Exchange
* TXT - Text strings

#### <a name='Classes'></a>Classes

The following classes are defined by the RFC:

* IN - Internet
* CS - CSNET (Obsolete and used only for examples in obsolete RFCs)
* CH - Chaos
* HS - Hesiod

#### <a name='SpecialValuesandEscapes'></a>Special Values and Escapes

There are some additional values that can be used to express arbitrary data:

* @ - When free standing, used to denote the current origin
* \X - X is any character other than a digit (0-9) and is used to escape the character (e.g. \@)
* \DDD - Where each D is a digit in the octal representation of the decimal number
* ( ) - Used to group data that crosses line boundaries.
* ; - Indicates a comment

### <a name='YAMLFormat'></a>YAML Format

Zonemgr parses a YAML with the following format:

```yaml
\<domain name\>: # The origin
  config:
    generate_serial: yes|no|true|false # If true, a serial number will be generated for you and any serial number specified will be ignored
    serial_change_index: integer # This value is only used if generate_serial is set to true, this value will be added to the end of the generated serial number to allow for multiple changes in the same day
    generate_reverse_lookup_zones: true # If true, any necessary reverse lookup zones x.x.x.in-addr.arpa will be created automatically
  ttl:
    value: 14400
    comment: Optional 32 bit time interval in seconds, the default TTL for each resource record that doesn't explicitly define one
  resource_records: # The full collection of resource records
    \<identifier\>: \<string\> # A unique name for the resource record. Some plugins may use this as the name field if 'name' is not present.
      name: \<string\> # The name of the record
      type: \<type\> # The resource record type, e.g. A, CNAME, SOA, NS, etc.
      class: \<class\> # Typically IN, for the default plugins, this will default to IN if not specified
      ttl: \<integer\> # optional 32 bit time intervale in seconds before this record should be refreshed
      values: # an arbibrary length set of values for the record, most resource records have a single value (e.g. for an A record it is the IP address of the host) but some, notably the SOA record, have a set of values
       - value: \<string\> # The value for the record, some plugins can leverage the identifiedr if this is missing
         comment: <\string\> # Optional comment for the value
```

#### <a name='Examples'></a>Examples

The following examples leverage the builtin plugins for the resource record types, please see the plugin documentation if using an alternative plugin.

##### NS record

Full example:

```yaml
ns1.example.com.:
  name: '@'
  type: NS
  class: IN
  TTL: 21600
  values:
    - value: ns1.example.com.
      comment: This is the primary nameserver
```

Minimal Example:

```yaml
ns1.example.com.:
  type: NS
```

##### A record

Full example:

```yaml
www:
  name: www
  type: A
  class: IN
  ttl: 14400
  values:
    - value: 1.2.3.4
      comment: The web server
```

Minimal Example:

```yaml
www:
  type: A
  value: 1.2.3.4
```

##### CNAME Record

Full Example:

```yaml
base:
  name: base
  type: CNAME
  class: IN
  ttl: 14400
  values:
    - value: www
      comment: alias for the www name
```

Minimal Example:

```yaml
base:
  type: CNAME
  value: www
```

##### SOA Record

Full Example:

```yaml
example.com:
  name: example.com.
  type: SOA
  class: IN
  TTL: 14400
  values:
    - value: n1.example.com.
      comment: Primary nameserver for the zone
    - value: admin@example.com OR admin.example.com
      comment: Mailbox of the person responsible for the zone
    - value: 20250803
      comment: unsigned 32 bit number representing the serial number (ignored if 'generate_serial' is true)
    - value: 7200
      comment: 32 bit time interval in seconds before the zone should be refreshed
    - value: 600
      comment: 32 bit time interval in seconds before a failed refresh should be retried
    - value: 3600000
      comment: 32 bit time intervale in seconds before a zone is no longer authoritative
    - value: 172800
      comment: unsinged 32 bit time interval in seconds before a negative cached response should be retired
```

Minimal Example:

```yaml
example.com:
  type: SOA
  values: 
    - value: n1.example.com.
    - value: admin@example.com OR admin.example.com
    - value: 7200
    - value: 600
    - value: 3600000
    - value: 172800
```

##### PTR Record

Full example:

```yaml
"1.2.3.4":
  name: "1.2.3.4"
  type: PTR
  class: IN
  ttl: 14400
  values:
   - value: www
     comment: Typically only used in reverse lookup zones (e.g. 3.2.1.in-addr.arpa) to provide IP to name lookup.
```

Minimal Example:

```yaml
"1.2.3.4":
  type: PTR
  value: www
```

### <a name='Built-InPlugins'></a>Built-In Plugins

The following are the built-in plugins which are provided out of the box, these plugins may be overridden:

* A
* NS
* CNAME
* SOA
* PTR

#### <a name='PluginBehavior'></a>Plugin Behavior

* All 'comment's are optional.
* All 'ttl's are optional.
* All 'class'es are optional and, where necessary, will default to IN if not specified.
* All dns name must be fully qualified, for example 'example.com.' and not just 'example.com'
* Any resource record with a single value can use 'value' and 'comment' as a short cut

##### NS

* The 'name' is optional, will default to "@" if not specified

##### SOA

* The value elements are printed in order, so they must follow one of the following patterns:
  * With explicit serial number:
    * MNAME
    * RNAME
    * SERIAL
    * REFRESH
    * RETRY
    * EXPIRE
    * NCACHE
  * With generated serial number:
    * MNAME
    * RNAME
    * REFRESH
    * RETRY
    * EXPIRE
    * NCACHE
* If 'generate_serial' is true but the explicit serial number is provided, it will be ignored.
* The primary name server (MNAME) is a DNS name and therefore must be fully qualified (see above)
* The administrator (RNAME) can either be specified as a valid email address (e.g. <admin@example.com>) or as the zone file specific format where the '@' is replaced by a dot ('.') (e.g. admin.example.com.). If using the latter, that's a specific name and needs to be fully qualified (see above)
