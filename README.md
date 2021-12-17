# udap

##### Universal Device Aggregation Platform

Modules can be configured to control computer settings, lights, music, media, or even spaceships.

### Goal

##### Provide a clean, robust, and reliable, centralization for controlling complex or simple items.

## Components

```

endpoints.endpoint-a1.

```
### Key Details

#### Hostnames
The udap module listens on 0.0.0.0, but more specifically on the udap transient & shadow networks, by the ip `10.0.2.2`

The remaining network is as follows:
```
2.1 vyOS Router
2.2 Udap (host server)
2.3 Access Point 
2.4 Access Point 
2.5 - 2.32 Reserved (Future Network)
2.33 - 2.100 Reserved (Wired Devices)
2.101 - 2.201 Access Point Allocation Pool
2.201 - 2.255 Transient Device DHCP Pool


10.0.2.2

```

#### Port
This udap module runs on the port `:3020` in development mode, and `:8327 (UDAP) ` in deployment

### Modules

Modules are general purpose plugins that can be used to control or interact with other devices or interfaces.

### Daemons

Daemons are unlike modules, a daemon cannot be directly triggered, it is only called when there is a state change from
the UDAP core or from a module.

## Protocols

### Endpoint Authentication

An access key can be generated by creating a new endpoint.

```http request
POST /register/{accessKey}
```

### Endpoint WebSockets Protocol

```json
{
  "target": "endpoint",
  "operation": "enroll",
  "body": {
    "id": "endpointId"
  }
}
```

```json
{
  "target": "endpoint|instance|entity",
  "operation": "create|delete|etc.",
  "body": {
    "instance": "instanceid"
  }
}
```