package main

import (
	"encoding/hex"
	"testing"
)

func mustHexDecode(hexString string) []byte {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return b
}

func TestIsMagicHash(t *testing.T) {
	tests := []struct {
		hash     []byte
		expected bool
	}{
		{mustHexDecode("0000000000000000000000000000000000000000"), true},
		{mustHexDecode("0e00000000000000000000000000000000000000"), true},
		{mustHexDecode("00e0000000000000000000000000000000000000"), true},
		{mustHexDecode("0e44215523245623123156788068234657654781"), true},
		{mustHexDecode("00e4215523245623123156788068234657654781"), true},
		{mustHexDecode("1000000000000000000000000000000000000000"), false},
		{mustHexDecode("1e00000000000000000000000000000000000000"), false},
		{mustHexDecode("0e0000000000000000000000000000000000000a"), false},
	}

	for _, test := range tests {
		actual := isMagicHash(test.hash)
		if actual != test.expected {
			t.Errorf("expected '%v' for hash '%x', got: '%v'",
				test.expected, test.hash, actual)
		}
	}
}

func BenchmarkIsMagicHash(b *testing.B) {
	hashOk := mustHexDecode("0e44215523245623123156788068234657654781")
	hashNok := mustHexDecode("ae4421552324562312315678866823465765478a")
	for n := 0; n < b.N; n++ {
		isMagicHash(hashOk)
		isMagicHash(hashNok)
	}
}
