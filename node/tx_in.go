package node

import "bytes"

type TxIn struct {
	Txid []byte // previous transaction ID connected with input
	Vout int    // previous transaction output index connected with input
	Sig  []byte // signature of transaction trimmed copy signed with transaction creator's private key
	Pkey []byte // previous transaction output owner connected with input. transaction creator's public key
}

func (in *TxIn) UnlockableWith(pkeyHash []byte) bool {
	return bytes.Compare(HashPkey(in.Pkey), pkeyHash) == 0
}
