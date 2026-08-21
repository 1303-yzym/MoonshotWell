package ucrypto

import (
	"testing"
)

func TestBcrypt(t *testing.T) {
	str := BcryptHash("123456")
	t.Log(str)
	// t.Log(BcryptCheck("12345678", "$2a$10$nDbqhqkuE.gzAQEA3SC/ge6n.29FoaMMbuvTcOZQn9T59Q.iML/TW"))
}
