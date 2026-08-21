package ucrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func NewAesKet(str string) (key string) {
	h := md5.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))[8:24]
}

// Pad 填充数据
func Pad(src []byte, blockSize int) []byte {
	padSize := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)

	return append(src, pad...)
}

// Pad 去掉填充数据
func unPad(src []byte) ([]byte, error) {
	n := len(src)

	unPadNum := int(src[n-1])
	if n <= unPadNum {
		return nil, errors.New("pad err")
	}

	return src[:n-unPadNum], nil
}

func AESEncrypt(key []byte, data []byte) ([]byte, error) {
	data = Pad(data, aes.BlockSize)
	c, _ := aes.NewCipher(key)
	out := make([]byte, len(data))
	c.Encrypt(out, data)

	return out, nil
}

func AESDecrypt(key []byte, data []byte) ([]byte, error) {
	c, _ := aes.NewCipher(key)
	out := make([]byte, len(data))
	c.Decrypt(out, data)
	out, _ = unPad(out)

	return out, nil
}

// AESCBCDecrypt 解析aes-cbc 有base64
func AESCBCDecrypt(key string, iv string, dts string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(dts)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(data, data)

	data, err = unPad(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
