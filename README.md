# Interfaces plugin

The interfaces plugin can be used to get information about interfaces and local IP-addresses on the system.

## Installation

Follow the [instructions](https://docs.halon.io/manual/comp_install.html#installation) in our manual to add our package repository and then run the below command.

### Ubuntu

```
apt-get install halon-extras-interfaces
```

### RHEL

```
yum install halon-extras-interfaces
```

## Exported functions

These functions needs to be [imported](https://docs.halon.io/hsl/structures.html#import) from the `extras://interfaces` module path.

### interfaces()

Get all the interfaces.

**Returns**

An array of the interfaces.

**Example**

```
import { interfaces } from "extras://interfaces";

echo interfaces(); // [0=>["name"=>"lo"],1=>["name"=>"tunl0"],2=>["name"=>"ip6tnl0"],3=>["name"=>"eth0"]]
```

### local_ips([interface])

Get all the local IP-addresses or only those for a specific interface.

**Params**

- interface `string` - The interface

**Returns**

An array of the local IP-addresses.

**Example**

```
import { local_ips } from "extras://interfaces";

echo local_ips(); // [0=>["address"=>"172.17.0.2"]]
echo local_ips("eth0"); // [0=>["address"=>"172.17.0.2"]]
```