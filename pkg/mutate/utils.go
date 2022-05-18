package mutate

import (
	"crypto/sha256"
	"encoding/base32"
	"regexp"
	"strings"
)

// The name derived for the imported service by MCS Controller
func DerivedName(namespace string, name string) string {
	hash := sha256.New()
	hash.Write([]byte(namespace + name))
	return "imported-" + strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(hash.Sum(nil)))[:10]
}

func MatchDerivedName(name string) bool {
	result := false
	regex, err := regexp.Compile(`imported-[a-z0-9]{10}`)
	if err != nil {
		return result
	}
	result = regex.MatchString(name)
	return result
}
