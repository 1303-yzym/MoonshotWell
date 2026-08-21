package ucrypto

import (
	"fmt"
	"testing"
)

const testAesKet = "63edb370ff55d3c3"

func TestNewAesKet(t *testing.T) {
	t.Log(NewAesKet("dc"))
}

func TestAESEncrypt(t *testing.T) {
	encrypt, err := AESEncrypt([]byte(testAesKet), []byte("dfacfcas"))
	if err != nil {
		return
	}
	t.Log(encrypt)
}

func TestAESDecrypt(t *testing.T) {
	key := []byte("63edb370ff55d3c3")
	iv := []byte("cc9cf5804edf1475")
	ciphertext := "PtllkPp05I0Lg1NOTudzLQ=="

	decrypt, err := AESCBCDecrypt(string(key), string(iv), ciphertext)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", decrypt)
}
