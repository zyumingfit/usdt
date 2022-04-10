package utils

import (
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

var (
	ErrDecodeLength = errors.New("base58 decode length error")
	ErrDecodeCheck  = errors.New("base58 check failed")
	ErrEncodeLength = errors.New("base58 encode length error")
)

// Decode by base58 and check.
func Decode58Check(input string) ([]byte, error) {
	decodeCheck := base58.Decode(input)
	if len(decodeCheck) <= 4 {
		return nil, ErrDecodeLength
	}
	decodeData := decodeCheck[:len(decodeCheck)-4]
	hash0, err := Hash(decodeData)
	if err != nil {
		return nil, err
	}
	hash1, err := Hash(hash0)
	if hash1 == nil {
		return nil, err
	}
	if hash1[0] == decodeCheck[len(decodeData)] && hash1[1] == decodeCheck[len(decodeData)+1] &&
		hash1[2] == decodeCheck[len(decodeData)+2] && hash1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, ErrDecodeCheck
}

// Encode by base58 and check.
func Encode58Check(input []byte) (string, error) {
	h0, err := Hash(input)
	if err != nil {
		return "", err
	}
	h1, err := Hash(h0)
	if err != nil {
		return "", err
	}
	if len(h1) < 4 {
		return "", ErrEncodeLength
	}
	inputCheck := append(input, h1[:4]...)

	return base58.Encode(inputCheck), nil
}
