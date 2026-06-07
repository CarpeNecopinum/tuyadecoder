package handlers

import (
	"bytes"
	"encoding/binary"
	"log"
	"path"
	"strconv"

	"github.com/dsnet/try"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DischargeModeHandler struct {
	ListenTopic string
	OutputTopic string
}

type DischargeMode struct {
	Ka0    int16
	Ka1    int16
	Ka2    int16
	Power0 int16
	Ka3    int16
	Ka4    int16
	Ka5    int16
	Ka6    int16
	Ka7    int16
	Ka8    int16
	Ka9    int16
	Ka10   int16
	Ka11   int16
	Ka12   int16
	Ka13   int16
	Ka14   int16
	Ka15   int16
	Ka16   int16
}

func (s *DischargeModeHandler) RegisterOn(c mqtt.Client) {
	defer try.F(log.Println)
	tok1 := c.Subscribe(s.ListenTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)
		data := DischargeMode{}

		r := bytes.NewBuffer(ParseBase64(m.Payload()))
		try.E(binary.Read(r, binary.BigEndian, &data))

		log.Printf("Read DischargeMode: %+v\n", data)

		or_topic := path.Join(s.OutputTopic, "output_rate", "state")
		tok := c.Publish(or_topic, 1, false, strconv.Itoa(int(data.Power0)))
		tok.Wait()
		try.E(tok.Error())
		log.Printf("Reported output_rate")
	})
	tok1.Wait()
}
