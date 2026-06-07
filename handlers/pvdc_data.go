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

type PVDCDataHandler struct {
	ListenTopic string
	OutputTopic string
}

type PVDCData struct {
	KaFlag0     int8
	KaFlag1     int8
	Ka1         int16
	Ka2         int16
	SolarPower0 int16
	SolarPower1 int16
	Ka4         int16
	Ka5         int8
}

func (s *PVDCDataHandler) RegisterOn(c mqtt.Client) {
	tok := c.Subscribe(s.ListenTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)

		data := PVDCData{}
		r := bytes.NewBuffer(ParseBase64(m.Payload()))
		try.E(binary.Read(r, binary.LittleEndian, &data))

		log.Printf("Read PVDCData: %+v\n", data)

		s1_topic := path.Join(s.OutputTopic, "solar_power0", "state")
		tok1 := c.Publish(s1_topic, 1, false, strconv.Itoa(int(data.SolarPower0)))
		s2_topic := path.Join(s.OutputTopic, "solar_power1", "state")
		tok2 := c.Publish(s2_topic, 1, false, strconv.Itoa(int(data.SolarPower1)))

		tok1.Wait()
		tok2.Wait()
		try.E(tok1.Error())
		try.E(tok2.Error())
		log.Printf("Reported solar_power")
	})
	tok.Wait()
	if tok.Error() != nil {
		log.Println("Error registering PVDCData listener: ", tok.Error())
	}
}
