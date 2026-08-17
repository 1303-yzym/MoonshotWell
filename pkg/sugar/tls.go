package sugar

import (
	"crypto/x509"
	"fmt"
	"os"
)

// LoadCertPool 从多个证书文件路径加载证书池
func LoadCertPool(certPaths ...string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, path := range certPaths {
		// 读取证书文件
		certBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate file '%s': %w", path, err)
		}
		// 添加证书到 CertPool
		if !certPool.AppendCertsFromPEM(certBytes) {
			return nil, fmt.Errorf("failed to append certificate from file '%s'", path)
		}
	}

	return certPool, nil
}
