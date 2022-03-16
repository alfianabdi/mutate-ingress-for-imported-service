package mutate

import (
	"crypto/sha256"
	"encoding/base32"
	"strings"
)

// The name derived for the imported service by MCS Controller
func DerivedName(namespace string, name string) string {
	hash := sha256.New()
	hash.Write([]byte(namespace + name))
	return "imported-" + strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(hash.Sum(nil)))[:10]
}
