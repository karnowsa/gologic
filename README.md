# gologic
## Features

  * Get status of Admin and Managed Server
  * Simple commands to start and stop Managed Server
  * Specify a list of Servers
  * Print info about Weblogic Servers


#### Supported Go Versions

Gologic was build on Version:

- 1.14.3

## Build

With this command you can build this project:

```bash
go build gologic.go AdminServer.go ManagedServer.go
```

## Installation

You need to create a config file at /etc/gologic.conf.  
With the followed content:
```json
{
    "ip": "127.0.0.1",
    "port": 7001,
    "username": "admin",
    "password": "securepassword",
}
```
