package ucrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	piv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	return piv, &piv.PublicKey
}

func computeSharedSecret(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, error) {
	x, _ := publicKey.Curve.ScalarMult(publicKey.X, publicKey.Y, privateKey.D.Bytes())
	sharedSecret := sha256.Sum256(x.Bytes())

	return sharedSecret[:], nil
}

// 加密
func encryptAESGCM(key []byte, plaintext []byte) (nonce []byte, ciphertext []byte, err error) {
	// 创建一个 AES 密钥块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	// 使用 GCM 模式加密
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	// 生成随机的 nonce（GCM 要求 nonce 是唯一的，但不需要是保密的）
	nonce = make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	// 使用密钥和 nonce 加密明文，生成密文和认证标签
	ciphertext = aesGCM.Seal(nil, nonce, plaintext, nil)

	return nonce, ciphertext, nil
}

// AES-GCM 解密
func decryptAESGCM(key []byte, nonce []byte, ciphertext []byte) ([]byte, error) {
	// 创建一个 AES 密钥块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 使用密钥和 nonce 解密密文
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
