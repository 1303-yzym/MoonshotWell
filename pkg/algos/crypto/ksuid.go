package ucrypto

import (
	"github.com/segmentio/ksuid"
)

func NewKsUid() string {
	id := ksuid.New()

	return id.String()
}
