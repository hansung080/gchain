package base58

import (
	"testing"
	"encoding/hex"
	"log"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestBase58(t *testing.T) {
	data, err := hex.DecodeString("0062E907B15CBF27D5425399EBF6F0FB50EBB88F18C29B7D93")
	if err != nil {
		log.Fatal(err)
	}

	encoded := Encode(data)
	assert.Equal(t, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", string(encoded))

	decoded := Decode([]byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
	assert.Equal(t, strings.ToLower("0062E907B15CBF27D5425399EBF6F0FB50EBB88F18C29B7D93"), hex.EncodeToString(decoded))
}
