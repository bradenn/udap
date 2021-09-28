# udap

Universal Device Aggregation Platform

### Workflow

##### Web Browser -> Server -> Modules -> Items to control

Modules can be configured to control computer settings, lights, music, media, or even spaceships.

### Goal

##### Provide a clean, robust, and reliable, centralization for controlling complex or simple items.

Create endpoint from diagnostic -> use join code on new endpoint -> get jwt

## Protocol

### Websocket Request Format

```json
{
  "token": "JWT Token",
  "type": "enroll|action|configuration|diagnostic",
  "payload": {}
}
```

#### Enrollment Payload

```json
{
  "instances": ["instanceId1", "instanceId2"]
}
```

#### Action Payload

```json
{
  "instance": "instanceId1",
  "action": "actionName"
}
```

#### Diagnostic Payload

```json
{
  "instance": "instanceId1",
  "action": "actionName"
}
```