# udap

Universal Device Aggregation Platform

### Workflow

##### Web Browser -> Server -> Modules -> Items to control

Modules can be configured to control computer settings, lights, music, media, or even spaceships.

### Goal

##### Provide a clean, robust, and reliable, centralization for controlling complex or simple items.

Create endpoint from diagnostic -> use join code on new endpoint -> get jwt

#### Websocket Request

```json
{
  "select": {
    "analytic": "temperature"
  },
  "from": {
    "instance": "InstanceId"
  },
  "at": {
    "frequency": "1/60"
  }
}
```

#### Returns

```json
{
  "instances": [
    {
      "name": "c",
      "id": "1234",
      "components": [
        {
          "name": "Temperature Inside",
          "type": "analytic",
          "frequency": "watch",
          "data": "78"
        },
        {
          "name": "Mode",
          "type": "state",
          "frequency": "watch",
          "data": "cool|off|heat|fan"
        }
      ]
    }
  ]
}
```

#### No... Lets try something more ✨dyanamic✨
```json

{
  "requests": ["thermostat.temp", "thermostat.fan"] 
}

```

