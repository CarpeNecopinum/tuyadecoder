# Tuyadecoder

A lightweight Go service to read values from MQTT, decode them and push them back.

Currently only supports the Lidl TRONIC Solar Battery (Tuya product category `bxsdy`).

Primarily meant to be used with [Rustuya Bridge](https://github.com/3735943886/rustuya-bridge).

## Setup

The program expects two environment variables to be set (or configured in an .env file):

| name        | description                                                           |
|-------------|-----------------------------------------------------------------------|
| MQTT_BROKER | URL of your MQTT broker, can include username and password, if needed |
| DEVICE_ID   | The id of the Tuya device you're working with                         |

## Config

To keep the project simple, most actual configuration just happens in the `main.go` file for now.

### Handlers

Most work happens in "Handlers". They are explicitly constructed in the main file and then `RegisterOn` is called with the MQTT client.
Their naming follows the names registered with Tuya for the respective data points they act on.

They assume the event topic template in Rustuya to be `"tuya/{id}/{dp}/state"` and the command topic to be `"rustuya/command"`.
