package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/dsnet/try"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func removeFrom(arr []byte, val byte) []byte {
	res := []byte{}
	for _, c := range arr {
		if c != val {
			res = append(res, c)
		}
	}
	return res
}

func ParseBase64(msg []byte) []byte {
	stripped := removeFrom(msg, '"') // comes with extra double quotes

	log.Printf("Base64: %s", stripped)
	buf := make([]byte, 64)
	n := try.E1(base64.StdEncoding.Decode(buf, stripped))
	return buf[:n]
}

func EncodeBase64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	// quoted := "\"" + string(encoded) + "\""
	return encoded
}

type setDischargeModeStruct struct {
	Action string            `json:"action"`
	Id     string            `json:"id"`
	Dps    map[string]string `json:"dps"`
}

func SetDpJson(device, dp, value string) []byte {
	dps := map[string]string{}
	dps[dp] = value
	bs, err := json.Marshal(setDischargeModeStruct{
		Action: "set",
		Id:     device,
		Dps:    dps,
	})
	if err != nil {
		log.Println("Error Marshalling set command: ", err)
	}
	return bs
}

const InputPrefix = "tuya"
const OutputPrefix = "tuyadecoder"

type Handler interface {
	RegisterOn(c mqtt.Client)
}
