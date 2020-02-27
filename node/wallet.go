package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"bytes"

	"golang.org/x/crypto/ripemd160"
	"github.com/hansung080/gchain/encoding/base58"
)

/**
  @ How to Create Address

                  Public Key
                      V
       RIPEMD160( SHA256( Public Key ) )
                      V
                      V   ---> SHA256( SHA256( Version + Public Key Hash ) )
                      V   |        V
    Version + Public Key Hash + Checksum = Payload (1 + 20 + 4 = 25 bytes)
                      V
             Base58Encode( Payload )
                      V
                   Address

  @ Address Example: First Bitcoin Address Owned by Satoshi Nakamoto
    - hex-encoded address
      0062E907B15CBF27D5425399EBF6F0FB50EBB88F18C29B7D93
      Version  Public Key Hash                           Checksum
      00       62E907B15CBF27D5425399EBF6F0FB50EBB88F18  C29B7D93

    - base58-encoded address
      1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
*/

const (
	addressVersion     = byte(0x00)
	addressVersionLen  = 1
	addressChecksumLen = 4
)

type Wallet struct {
	Skey ecdsa.PrivateKey // Private key is a random value.
	Pkey []byte // Public key is (x, y) on the elliptic curve. `pkey` combines x with y as a byte array.
}

func (w Wallet) GetAddress() []byte {
	pkeyHash := HashPkey(w.Pkey)
	versionedPkeyHash := append([]byte{addressVersion}, pkeyHash...)
	checksum := newChecksum(versionedPkeyHash)
	payload := append(versionedPkeyHash, checksum...)
	return base58.Encode(payload)
}

func NewWallet() *Wallet {
	skey, pkey := newKeyPair()
	return &Wallet{
		Skey: skey,
		Pkey: pkey,
	}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	skey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pkey := append(skey.PublicKey.X.Bytes(), skey.PublicKey.Y.Bytes()...)
	return *skey, pkey
}

func HashPkey(pkey []byte) []byte {
	pkeySHA256 := sha256.Sum256(pkey)

	hasherRIPEMD160 := ripemd160.New()
	if _, err := hasherRIPEMD160.Write(pkeySHA256[:]); err != nil {
		log.Panic(err)
	}

	pkeyRIPEMD160 := hasherRIPEMD160.Sum(nil)
	return pkeyRIPEMD160
}

func newChecksum(data []byte) []byte {
	firstSHA256 := sha256.Sum256(data)
	secondSHA256 := sha256.Sum256(firstSHA256[:])
	return secondSHA256[:addressChecksumLen]
}

func ValidateAddress(addr string) bool {
	payload := base58.Decode([]byte(addr))
	actualChecksum := payload[len(payload) - addressChecksumLen:]
	version := payload[0]
	pkeyHash := payload[1:len(payload) - addressChecksumLen]
	targetChecksum := newChecksum(append([]byte{version}, pkeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
