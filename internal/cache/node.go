package cache

import (
	"context"

	cachepb "github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
)

type CacheNode struct {
	cachepb.UnimplementedCacheServer
	lru *LruCache
}

func NewCacheNode(cap int) cachepb.CacheServer {
	return &CacheNode{
		lru: &LruCache{
			capacity: cap,
		},
	}
}

func (cn *CacheNode) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	val, err := cn.lru.get(KEY(req.Key))
	if err != nil {
		return &cachepb.GetResponse{Found: false}, nil
	}

	return &cachepb.GetResponse{Value: val, Found: true}, nil
}

func (cn *CacheNode) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.SetResponse, error) {
	_ = cn.lru.set(KEY(req.Key), req.Value)

	return &cachepb.SetResponse{Success: true}, nil
}
