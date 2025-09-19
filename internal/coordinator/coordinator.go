package coordinator

import (
	"context"
	"sort"

	"github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
	"github.com/sakshamg567/cachy/util"
)

type Coordinator struct {
	ring *HashRing
}

func NewCoordinator(addresses []string) *Coordinator {

	ring := NewHashRing(addresses)

	return &Coordinator{
		ring: ring,
	}
}

func (c *Coordinator) Get(ctx context.Context, key string) (string, error) {

	n := c.ring.getNode(key)

	val, err := n.client.Get(ctx, &cacheNodepb.GetRequest{Key: key})
	if err != nil {
		return "", err
	}

	return val.Value, nil
}

func (c *Coordinator) Set(ctx context.Context, key, value string) bool {
	n := c.ring.getNode(key)

	_, err := n.client.Set(ctx, &cacheNodepb.SetRequest{Key: key, Value: value})
	return err == nil
}

func (c *Coordinator) AddNode(addr string) {
	c.ring.addNode(addr)

	c.ring.mu.RLock()
	h := util.Hash(addr)
	idx := sort.Search(len(c.ring.keys), func(i int) bool { return c.ring.keys[i] >= h })
	if idx == len(c.ring.keys) {
		idx = 0
	}
	successorAddr := c.ring.nodes[c.ring.keys[idx]].addr
	predIdx := idx - 1
	if predIdx < 0 {
		predIdx = len(c.ring.keys) - 1
	}
	predecessorHash := c.ring.keys[predIdx]
	c.ring.mu.RUnlock()

	go c.ring.migrateData(successorAddr, addr, predecessorHash)
}
