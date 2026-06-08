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
	DeviceId    string

	lastState []byte
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

var dischargeModeEndian = binary.BigEndian

func bin2struct(src []byte, dst any) error {
	r := bytes.NewBuffer(src)
	return binary.Read(r, dischargeModeEndian, dst)
}
func struct2bin(src any) ([]byte, error) {
	r := bytes.Buffer{}
	err := binary.Write(&r, dischargeModeEndian, src)
	if err != nil {
		return nil, err
	}
	return r.Bytes(), nil

}

func (s *DischargeModeHandler) RegisterOn(c mqtt.Client) {
	defer try.F(log.Println)
	tok1 := c.Subscribe(s.ListenTopic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)
		data := DischargeMode{}

		s.lastState = ParseBase64(m.Payload())
		try.E(bin2struct(s.lastState, &data))

		log.Printf("Read DischargeMode: %+v\n", data)

		or_topic := path.Join(s.OutputTopic, "output_rate", "state")
		tok := c.Publish(or_topic, 1, false, strconv.Itoa(int(data.Power0)))
		tok.Wait()
		try.E(tok.Error())
		log.Printf("Reported output_rate")
	})
	cmd_topic := path.Join(s.OutputTopic, "output_rate", "set")
	tok2 := c.Subscribe(cmd_topic, 1, func(c mqtt.Client, m mqtt.Message) {
		defer try.F(log.Println)
		if s.lastState == nil {
			return
		}

		data := DischargeMode{}
		try.E(bin2struct(s.lastState, &data))
		data.Power0 = int16(try.E1(strconv.Atoi(string(m.Payload()))))
		log.Println("Pushing new output rate: ", string(m.Payload()), data)
		s.lastState = try.E1(struct2bin(&data))

		stateStr := EncodeBase64(s.lastState)
		json := SetDpJson(s.DeviceId, "106", stateStr)
		log.Println("Will send json: ", string(json))
		tok := c.Publish(path.Join("rustuya", "command"), 1, false, json)
		tok.Wait()
		try.E(tok.Error())
	})
	tok1.Wait()
	try.E(tok1.Error())
	tok2.Wait()
	try.E(tok2.Error())
}
