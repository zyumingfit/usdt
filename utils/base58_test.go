package utils

import (
	"testing"
)

var (
	testBase58DecodeString = "test"
	testBase58EncodeString = "LUC1eAJa5jW"
)

// Test decode58 function.
func TestDecode58Check(t *testing.T) {
	decode, err := Decode58Check(testBase58EncodeString)
	if err != nil {
		t.Errorf("Decode58 address error, reasons: [%v]", err)
		return
	}

	if testBase58DecodeString != string(decode) {
		t.Error("Test Decode58Check function failed!")
	} else {
		t.Log("Test Decode58Check function success!")
	}
}

// Test encode58 function.
func TestEncode58Check(t *testing.T) {
	encode, err := Encode58Check([]byte(testBase58DecodeString))
	if err != nil {
		t.Errorf("Encode58 address error, reasons: [%v]", err)
		return
	}

	if encode != testBase58EncodeString {
		t.Error("Test Encode58Check function failed!")
	} else {
		t.Log("Test Encode58Check function success!")
	}
}
