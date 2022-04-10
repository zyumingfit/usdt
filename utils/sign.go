package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"

	eth "github.com/ethereum/go-ethereum/crypto"
)

// Sign data by ecdsa private key.
func Sign(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	hash, err := Hash(data)
	if err != nil {
		return nil, err
	}

	signature, err := eth.Sign(hash, key)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

//Package goLang sha256 hash algorithm.
func Hash(s []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}
