# Tuyadecoder

A lightweight Go service to read values from MQTT, decode them and push them back.

Primarily meant to be used with [Rustuya Bridge](https://github.com/3735943886/rustuya-bridge).

## Setup

The program expects two environment variables to be set (or configured in an .env file):

|-------------|-----------------------------------------------------------------------|
| MQTT_BROKER | URL of your MQTT broker, can include username and password, if needed |
| DEVICE_ID   | The id of the Tuya device you're working with                         |

## Config

To keep the project simple, most actual configuration just happens in the `main.go` file for now.
