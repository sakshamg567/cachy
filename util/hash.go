package util

import (
	"crypto/sha256"
	"math/big"
)

const total_slots = uint64(1) << 32

func Hash(id string) uint32 {
	h := sha256.Sum256([]byte(id))
	num := new(big.Int).SetBytes(h[:])
	return uint32(new(big.Int).Mod(num, big.NewInt(int64(total_slots))).Int64())
}
