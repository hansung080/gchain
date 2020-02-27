package base58

import (
	"math/big"
	"bytes"
)

// Base58 gets rid of "0OIl+/" from base64, because those characters are ambiguous for human to distinguish.
var b58Characters = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Encode(src []byte) []byte {
	var dest []byte

	x := big.NewInt(0).SetBytes(src)
	base := big.NewInt(int64(len(b58Characters)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		dest = append(dest, b58Characters[mod.Int64()])
	}

	// https://en.bitcoin.it/wiki/Base58Check_encoding#Version_bytes
	if src[0] == 0x00 {
		dest = append(dest, b58Characters[0])
	}

	reverseBytes(dest)
	return dest
}

func Decode(src []byte) []byte {
	dest := big.NewInt(0)

	for _, b := range src {
		index := bytes.IndexByte(b58Characters, b)
		dest.Mul(dest, big.NewInt(58))
		dest.Add(dest, big.NewInt(int64(index)))
	}

	decoded := dest.Bytes()
	if src[0] == b58Characters[0] {
		decoded = append([]byte{0x00}, decoded...)
	}

	return decoded
}
