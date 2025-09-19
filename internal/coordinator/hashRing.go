package coordinator

import (
	"context"
	"log"
	"sort"
	"sync"

	"github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
	"github.com/sakshamg567/cachy/util"
	"google.golang.org/grpc"
)

type node struct {
	addr   string
	client cacheNodepb.CacheClient
}

type HashRing struct {
	nodes map[uint32]node
	keys  []uint32
	mu    sync.RWMutex
}

func NewHashRing(addresses []string) *HashRing {
	n := make(map[uint32]node)
	var keys []uint32

	for _, addr := range addresses {
		h := util.Hash(addr)
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		client := cacheNodepb.NewCacheClient(conn)
		n[h] = node{
			addr:   addr,
			client: client,
		}
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

func (r *HashRing) addNode(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h := util.Hash(addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := cacheNodepb.NewCacheClient(conn)
	r.nodes[h] = node{
		addr:   addr,
		client: client,
	}
	r.keys = append(r.keys, h)
	sort.Slice(r.keys, func(i, j int) bool {
		return r.keys[i] < r.keys[j]
	})
}

func (r *HashRing) removeNode(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h := util.Hash(addr)
	delete(r.nodes, h)
	for i, key := range r.keys {
		if key == h {
			r.keys = append(r.keys[:i], r.keys[i+1:]...)
			break
		}
	}
}

func (r *HashRing) migrateData(addrFrom, addrTo string, predHash uint32) error {
	r.mu.RLock()
	hFrom := util.Hash(addrFrom)
	hTo := util.Hash(addrTo)

	nodeFrom := r.nodes[hFrom]
	nodeTo := r.nodes[hTo]
	r.mu.RUnlock()

	keys, err := nodeFrom.client.GetAllKeys(context.Background(), &cacheNodepb.GetAllKeysRequest{})
	if err != nil {
		return err
	}

	for _, key := range keys.Keys {
		kh := util.Hash(key)

		shouldMove := false
		if predHash > hTo {
			if kh > predHash || kh <= hTo {
				shouldMove = true
			}
		} else {
			if kh > predHash && kh <= hTo {
				shouldMove = true
			}
		}
		if shouldMove {
			getRes, err := nodeFrom.client.Get(context.Background(), &cacheNodepb.GetRequest{Key: key})
			if err != nil || !getRes.Found {
				continue
			}

			_, err = nodeTo.client.Set(context.Background(), &cacheNodepb.SetRequest{
				Key:   key,
				Value: getRes.Value,
			})
			if err != nil {
				continue
			}

			_, err = nodeFrom.client.Delete(context.Background(), &cacheNodepb.DeleteRequest{Key: key})
			if err != nil {
				continue
			}
		}
	}
	return nil
}

func (r *HashRing) getNode(key string) node {
	h := util.Hash(key)
	log.Printf("key hash : %v", h)
	log.Printf("node hashes : %v", r.keys)

	r.mu.RLock()
	defer r.mu.RUnlock()

	idx := sort.Search(len(r.keys), func(i int) bool {
		return r.keys[i] >= h
	})
	if idx == len(r.keys) {
		idx = 0
	}
	n := r.nodes[r.keys[idx]]
	log.Printf("[ring] key=%s hash=%d -> node=%s nodeHash=%d (idx=%d)", key, h, n.addr, r.keys[idx], idx)

	log.Printf("ring.keys=%v (sorted=%t)", r.keys, sort.SliceIsSorted(r.keys, func(i, j int) bool { return r.keys[i] < r.keys[j] }))

	return n
}
