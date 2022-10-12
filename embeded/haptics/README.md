### ESP WiFi DMX

Control DMX devices capable of RS485 with an ESP32. (MAX485 or generic alternative required)

### Endpoints

##### Request

Change the current value of a channel, this will immediately change the dmx state.

`POST /channel`

```json
{
  "channel": 1,
  "value": 0
}
```

##### Response

`200: okay`

##### Request

Set the default value when the program starts. When a loss of power occurs, this value will be read from NVS and will be
applied to the DMX buffer. This change is not immediatly apparent.

`POST /default`

```json
{
  "channel": 1,
  "value": 0
}
```

##### Response

`200: okay`

##### Request

Return an array of the channels' current values.

`GET /status`

##### Response

```json
{
    "channels": [
        {
            "channel": 1,
            "value": 0
        },
        {
            "channel": 2,
            "value": 0
        },
        {
            "channel": 3,
            "value": 0
        },
        {
            "channel": 4,
            "value": 0
        },
        {
            "channel": 5,
            "value": 0
        },
        {
            "channel": 6,
            "value": 0
        },
        {
            "channel": 7,
            "value": 0
        },
        {
            "channel": 8,
            "value": 0
        },
        {
            "channel": 9,
            "value": 0
        },
        {
            "channel": 10,
            "value": 0
        },
        {
            "channel": 11,
            "value": 0
        },
        {
            "channel": 12,
            "value": 0
        },
        {
            "channel": 13,
            "value": 0
        },
        {
            "channel": 14,
            "value": 0
        },
        {
            "channel": 15,
            "value": 0
        },
        {
            "channel": 16,
            "value": 0
        }
    ]
}
```

This program contains code from (for DMX over RS485):

https://github.com/luksal/ESP32-DMX