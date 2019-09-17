# PIOT MQTT Adapter

Translates messages from PIOT proprietary format to MQTT and publishes
converted messages to MQTT broker.

Format of the message:
```
{
    "device": "xxx"
    "readings": [{
        "address": "xyz",
        "t": "number",
        "h": "number"
    }]
}
```
