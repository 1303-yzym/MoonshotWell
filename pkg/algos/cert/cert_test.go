package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

const Country = "US"

func TestGenRoot(t *testing.T) {
	// 1. 生成 RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	// 2. 定义证书模板
	rootCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{Country},
			Organization:       []string{"ALO"},
			OrganizationalUnit: []string{"ALO ROOT CA"},
			SerialNumber:       "1024",
			CommonName:         "ALO",
			// 地域
			Locality:      []string{"ALO"},
			Province:      []string{"ALO"},
			StreetAddress: []string{"ALO"},
			PostalCode:    []string{"100101010"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 有效期 10 年
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  true, // 标记为 CA 证书
		MaxPathLen:            2,    // 子证书的最大层级
		SignatureAlgorithm:    x509.SHA384WithRSA,
	}

	// 3. 自签名生成证书
	rootCertBytes, err := x509.CreateCertificate(rand.Reader, rootCertTemplate, rootCertTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// 4. 将私钥写入文件
	privateKeyFile, err := os.Create("root.key")
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return
	}

	// 5. 将根证书写入文件
	rootCertFile, err := os.Create("root.crt")
	if err != nil {
		panic(err)
	}
	defer rootCertFile.Close()

	err = pem.Encode(rootCertFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCertBytes,
	})
	if err != nil {
		return
	}

	println("自签名根证书和私钥生成成功：root.crt 和 root.key")

}

func TestGenMiddle(t *testing.T) {
	fileName := "p1"
	// 1. 加载根证书和私钥
	rootCert, rootKey, err := loadRootCertificate("./root.crt", "./root.key")
	if err != nil {
		t.Fatalf("Failed to load root certificate or key: %v", err)
	}

	// 2. 生成中间证书的私钥
	intermediatePrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// 3. 定义中间证书模板
	intermediateCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Country:            []string{Country},
			Organization:       []string{"p1"},
			OrganizationalUnit: []string{"p1"},
			CommonName:         "p1",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.SHA384WithRSA,
	}

	// 4. 使用根证书和私钥签名中间证书
	intermediateCertBytes, err := x509.CreateCertificate(rand.Reader, intermediateCertTemplate, rootCert, &intermediatePrivateKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("Failed to create intermediate certificate: %v", err)
	}

	// 5. 将中间证书写入文件
	intermediateCertFile, err := os.Create(fileName + ".crt")
	if err != nil {
		t.Fatalf("Failed to create certificate file: %v", err)
	}
	defer func() {
		if cerr := intermediateCertFile.Close(); cerr != nil {
			t.Fatalf("Failed to close certificate file: %v", cerr)
		}
	}()

	err = pem.Encode(intermediateCertFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intermediateCertBytes,
	})
	if err != nil {
		t.Fatalf("Failed to encode intermediate certificate: %v", err)
	}

	// 6. 将中间证书的私钥写入文件
	intermediatePrivateKeyFile, err := os.Create(fileName + ".key")
	if err != nil {
		t.Fatalf("Failed to create private key file: %v", err)
	}
	defer func() {
		if cerr := intermediatePrivateKeyFile.Close(); cerr != nil {
			t.Fatalf("Failed to close private key file: %v", cerr)
		}
	}()

	err = pem.Encode(intermediatePrivateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(intermediatePrivateKey),
	})
	if err != nil {
		t.Fatalf("Failed to encode private key: %v", err)
	}

	t.Logf("中间证书和私钥生成成功：%s.crt 和 %s.key", fileName, fileName)
}

func TestGenEntity(t *testing.T) {
	fileName := "nginx"
	// 1. 加载根证书和私钥
	rootCert, rootKey, err := loadRootCertificate("./p1.crt", "./p1.key")
	if err != nil {
		panic(err)
	}
	// 加载中间证书（假设你有中间证书文件 intermediate.crt）
	intermediateCertPEM, err := os.ReadFile("./p1.crt")
	if err != nil {
		panic(err)
	}
	// 2. 生成实体证书的私钥
	entityPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 3. 定义实体证书模板
	entityCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Country:            []string{Country},
			Organization:       []string{"c"},
			OrganizationalUnit: []string{"c"},
			CommonName:         "c",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "nginx", "wlc_nginx"},
		// 配置 IP 地址
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("202.79.174.3"),
		},
	}

	// 4. 使用根证书和私钥签名实体证书
	entityCertBytes, err := x509.CreateCertificate(rand.Reader, entityCertTemplate, rootCert, &entityPrivateKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}

	// 5. 将实体证书写入文件
	entityCertFile, err := os.Create(fileName + ".crt")
	if err != nil {
		panic(err)
	}
	defer entityCertFile.Close()

	err = pem.Encode(entityCertFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: entityCertBytes,
	})
	if err != nil {
		panic(err)
	}

	// 写入中间证书内容
	_, err = entityCertFile.Write(intermediateCertPEM)
	if err != nil {
		panic(err)
	}

	// 6. 将实体证书的私钥写入文件
	entityPrivateKeyFile, err := os.Create(fileName + ".key")
	if err != nil {
		panic(err)
	}
	defer entityPrivateKeyFile.Close()

	err = pem.Encode(entityPrivateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(entityPrivateKey),
	})
	if err != nil {
		panic(err)
	}

	println("实体证书和私钥生成成功：entity.crt 和 entity.key")
}
