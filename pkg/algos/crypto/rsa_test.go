package ucrypto

import (
	"testing"
)

func TestGenKeys(t *testing.T) {
	// GenKeys("./certs")
}

//func TestDecryption(t *testing.T) {
//	rs, _ := NewRsa("./certs/app_private.pem")
//	// str := "Hello world!!!"
//
//	c := "DZxvltk7cp0mz6gyYu/Icdbd5evF4ATtx5C5B9Nhh9fae1ittEr8r71nwkmCQmwAfGj68fiHuaa3Otg6c4eDgDpw/hZgaktQIpayoQoRcT+2LQ9OziDWKev49hJjjXr4ulIkaJgd/3DFv5x+TTavuKrI7wi/HxYdHILM74/HiNYpuRdviTPM3FT6w+VpnuApDBl28SxlkeyhnJg0OJWoRfUUr7Pe4FKLMK6LbWq0GpuuHALir2Z+RY/NzfR3VHZAbijeZDuDh2OWFuoX0Z8iMsdoz6qF+B5L9miGMXo13SX/rOrNkMdf5W+vewTVFPh7EXDBLc7su+RUcAs9Rf6BpA=="
//	decrypt, err := rs.DecryptBase64(c)
//	if err != nil {
//		panic(err)
//		return
//	}
//	t.Log(string(decrypt))
//}

func TestEncryption(t *testing.T) {
	//// // 读取公钥文件
	//publicKeyData, err := os.ReadFile("./certs/app_public.pem")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// // 解码公钥数据
	//block, _ := pem.Decode(publicKeyData)
	//if block == nil {
	//	log.Fatal("无法解码公钥文件")
	//}
	//// // 解析公钥
	//publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// // 要加密的数据
	//plaintext := []byte("17502362036#1723535538")
	//// // 使用公钥加密数据
	//ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//encryptedBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	//fmt.Println(encryptedBase64)
}
