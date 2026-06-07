package handlers

import (
	"encoding/base64"
	"log"

	"github.com/dsnet/try"
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
