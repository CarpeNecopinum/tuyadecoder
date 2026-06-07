package handlers

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/dsnet/try"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DCMessageHandler struct {
	ListenTopic string
}

type DCMessage struct {
	Ka0 int16
	Ka1 int16
	Ka2 int16
	Ka3 int16
	Ka4 int16
	Ka5 int16
	Ka6 int8
}

func (s *DCMessageHandler) RegisterOn(c mqtt.Client) {
	tok := c.Subscribe(s.ListenTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)

		data := DCMessage{}
		r := bytes.NewBuffer(ParseBase64(m.Payload()))
		try.E(binary.Read(r, binary.BigEndian, &data))

		log.Printf("Read DCMessage: %+v\n", data)
	})
	tok.Wait()
	if tok.Error() != nil {
		log.Println("Error registering DCMessage listener: ", tok.Error())
	}
}
