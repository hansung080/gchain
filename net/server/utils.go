package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

func isNodeKnown(addr string) bool {
	for _, a := range knownAddrs {
		if a == addr {
			return true
		}
	}

	return false
}

func commandToBytes(cmd string) []byte {
	var bytes [commandLen]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x00 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func extractCommand(req []byte) []byte {
	return req[:commandLen]
}

func marshalGob(v interface{}) []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(v); err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func unmarshalGob(data []byte, v interface{}) {
	var buf bytes.Buffer

	buf.Write(data)
	decoder := gob.NewDecoder(&buf)
	if err := decoder.Decode(v); err != nil {
		log.Panic(err)
	}
}
