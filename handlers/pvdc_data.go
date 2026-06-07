package handlers

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/dsnet/try"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PVDCDataHandler struct {
	ListenTopic string
}

type PVDCData struct {
	KaFlag0      int8
	KaFlag1      int8
	Ka1          int16
	Ka2          int16
	CurrentPower int16
	Ka3          int16
	Ka4          int16
	Ka5          int8
}

func (s *PVDCDataHandler) RegisterOn(c mqtt.Client) {
	tok := c.Subscribe(s.ListenTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)

		data := PVDCData{}
		r := bytes.NewBuffer(ParseBase64(m.Payload()))
		try.E(binary.Read(r, binary.BigEndian, &data))

		log.Printf("Read PVDCData: %+v\n", data)
	})
	tok.Wait()
	if tok.Error() != nil {
		log.Println("Error registering PVDCData listener: ", tok.Error())
	}
}
