package util

import (
	"hash/fnv"
)

func Hash(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
