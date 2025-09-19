package coordinator

import (
	"context"

	"github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
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

	val, err := n.Get(ctx, &cacheNodepb.GetRequest{Key: key})
	if err != nil {
		return "", err
	}

	return val.Value, nil
}

func (c *Coordinator) Set(ctx context.Context, key, value string) bool {
	n := c.ring.getNode(key)

	_, err := n.Set(ctx, &cacheNodepb.SetRequest{Key: key, Value: value})
	if err != nil {
		return false
	}

	return true
}
