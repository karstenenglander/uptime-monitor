package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

func GetFNVHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func GetSHA256Hash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
