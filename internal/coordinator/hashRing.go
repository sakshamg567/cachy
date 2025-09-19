package coordinator

import (
	"sort"

	"github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
	"github.com/sakshamg567/cachy/util"
	"google.golang.org/grpc"
)

type HashRing struct {
	nodes map[uint32]cacheNodepb.CacheClient
	keys  []uint32
}

func NewHashRing(addresses []string) *HashRing {
	n := make(map[uint32]cacheNodepb.CacheClient)
	var keys []uint32

	for _, addr := range addresses {
		h := util.Hash(addr)
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		client := cacheNodepb.NewCacheClient(conn)
		n[h] = client
		keys = append(keys, h)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return &HashRing{
		nodes: n,
		keys:  keys,
	}
}

// func (r *HashRing) AddNode(node *CacheNode) {
// 	h := hash(node.ip)
// 	r.nodes[h] = node
// 	r.keys = append(r.keys, h)
// 	sort.Slice(r.keys, func(i, j int) bool {
// 		return r.keys[i] < r.keys[j]
// 	})
// }

func (r *HashRing) getNode(key string) cacheNodepb.CacheClient {
	h := util.Hash(key)

	idx := sort.Search(len(r.keys), func(i int) bool {
		return r.keys[i] >= h
	})

	if idx == len(r.keys) {
		idx = 0
	}

	return r.nodes[r.keys[idx]]
}
