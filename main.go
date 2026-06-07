package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"main/handlers"

	"github.com/dsnet/try"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const clientid = "tuyadecoder"

func onMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func onConnection(client mqtt.Client, notification mqtt.ConnectionNotification) {
	switch n := notification.(type) {
	case mqtt.ConnectionNotificationConnected:
		log.Printf("[NOTIFICATION] connected\n")
	case mqtt.ConnectionNotificationConnecting:
		log.Printf("[NOTIFICATION] connecting (isReconnect=%t) [%d]\n", n.IsReconnect, n.Attempt)
	case mqtt.ConnectionNotificationFailed:
		log.Printf("[NOTIFICATION] connection failed: %v\n", n.Reason)
	case mqtt.ConnectionNotificationLost:
		log.Printf("[NOTIFICATION] connection lost: %v\n", n.Reason)
		client.Connect()
	case mqtt.ConnectionNotificationBroker:
		log.Printf("[NOTIFICATION] broker connection: %s\n", n.Broker.String())
	case mqtt.ConnectionNotificationBrokerFailed:
		log.Printf("[NOTIFICATION] broker connection failed: %v [%s]\n", n.Reason, n.Broker.String())
	}
}

func main() {
	defer try.F(log.Fatal)

	env := ParseEnv()
	broker := env["MQTT_BROKER"]
	device := env["DEVICE_ID"]

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientid).
		SetKeepAlive(2 * time.Second).
		SetDefaultPublishHandler(onMessage).
		SetConnectionNotificationHandler(onConnection).
		SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	connect_token := c.Connect()
	connect_token.Wait()
	try.E(connect_token.Error())

	dcm := handlers.DCMessageHandler{ListenTopic: path.Join("tuya", device, "33/state")}
	dcm.RegisterOn(c)

	dm := handlers.DischargeModeHandler{
		ListenTopic: path.Join("tuya", device, "106/state"),
		OutputTopic: path.Join("tuyadecoder", device),
	}
	dm.RegisterOn(c)

	pvdc := handlers.PVDCDataHandler{
		ListenTopic: path.Join("tuya", device, "101/state"),
		OutputTopic: path.Join("tuyadecoder", device),
	}
	pvdc.RegisterOn(c)

	log.Println("Subscribed")

	<-stop
}
