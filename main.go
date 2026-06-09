package main

import (
	"fmt"
	"log"
	"main/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		SetClientID(fmt.Sprintf("tuyadecoder-%d", time.Now().UnixNano())).
		SetKeepAlive(4 * time.Second).
		SetDefaultPublishHandler(onMessage).
		SetConnectionNotificationHandler(onConnection).
		SetPingTimeout(2 * time.Second)

	c := mqtt.NewClient(opts)
	connect_token := c.Connect()
	connect_token.Wait()
	try.E(connect_token.Error())

	hnds := make([]handlers.Handler, 3)
	hnds[0] = &handlers.DCMessageHandler{DeviceId: device}
	hnds[1] = &handlers.DischargeModeHandler{DeviceId: device}
	hnds[2] = &handlers.PVDCDataHandler{DeviceId: device}

	for _, h := range hnds {
		h.RegisterOn(c)
	}

	log.Println("Subscribed")

	<-stop
}
