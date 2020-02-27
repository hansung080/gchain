package node

import (
	"bytes"
	"encoding/binary"
	"github.com/hansung080/gchain/encoding/base58"
	"log"
	"os"
)

func IntToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func FileExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}

	return true
}

func GetPkeyHashFromAddress(addr []byte) []byte {
	payload := base58.Decode(addr)
	return payload[addressVersionLen:len(payload) - addressChecksumLen]
}
