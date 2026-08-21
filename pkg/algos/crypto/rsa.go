package ucrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

// GenKeys 生成秘钥对
func GenKeys(path string) {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // 使用2048位的密钥长度
	if err != nil {
		panic(err)
	}

	privateKeyFile, err := os.Create(path + "/app_private.pem")
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()

	publicKeyFile, err := os.Create(path + "/app_public.pem")
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()
	// 存储私钥
	err = pem.Encode(privateKeyFile, &pem.Block{
		Type:  "DC_SP RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		panic(err)
	}
	// 存储公钥
	publicKeyPEM := &pem.Block{
		Type:  "DC_SP PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		panic(err)
	}
}

type DRsa struct {
	P *rsa.PrivateKey
}

func NewRsa(path string) (DRsa, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return DRsa{}, errors.New("无法打开私钥文件" + err.Error())
	}
	// 解码私钥数据
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return DRsa{}, errors.New("无法解析私钥文件")
	}
	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return DRsa{}, errors.New("无法解析私钥" + err.Error())
	}

	return DRsa{P: privateKey}, nil
}

// DecryptBase64 解析base64的加密数据
func (u DRsa) DecryptBase64(data string) (decrypt []byte, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return decrypt, errors.New("base64错误 RFC 4648")
	}

	return u.P.Decrypt(nil, ciphertext, &rsa.PKCS1v15DecryptOptions{})
}
