package ucrypto

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestAs(t *testing.T) {
	// 生成密钥对
	pivClient, pubClient := generateKeyPair()
	pivServer, pubServer := generateKeyPair()

	// 生成共享密钥
	sharedSecret1, _ := computeSharedSecret(pivClient, pubServer)
	sharedSecret2, _ := computeSharedSecret(pivServer, pubClient)

	fmt.Printf("Shared Secret (Party 1): %x\n", sharedSecret1)
	fmt.Printf("Shared Secret (Party 2): %x\n", sharedSecret2)

	log.Println(len(fmt.Sprintf("%x", sharedSecret2)))

	if hex.EncodeToString(sharedSecret1) == hex.EncodeToString(sharedSecret2) {
		fmt.Println("Shared secrets match!")
	} else {
		fmt.Println("Shared secrets do not match!")
	}

	data := []byte("hallo")

	nonce, ciphertext, err := encryptAESGCM(sharedSecret1, data)
	if err != nil {
		return
	}

	log.Println(nonce, ciphertext)

	decryptData, err := decryptAESGCM(sharedSecret2, nonce, ciphertext)
	if err != nil {
		return
	}

	if hex.EncodeToString(data) == hex.EncodeToString(decryptData) {
		fmt.Println("decrypt data match!")
	} else {
		fmt.Println("decrypt data do not match!")
	}

	log.Println(string(decryptData))
}
