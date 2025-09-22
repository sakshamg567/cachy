package cache

import (
	"context"
	"log"

	cachepb "github.com/sakshamg567/cachy/shared/proto/cacheNodepb"
)

type CacheNode struct {
	cachepb.UnimplementedCacheServer
	lru *LruCache
}

func NewCacheNode(cap int) cachepb.CacheServer {
	lru := NewLruCache(cap)

	return &CacheNode{
		lru: lru,
	}
}

func (cn *CacheNode) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	log.Printf("RPC Get key=%q", req.Key)
	val, err := cn.lru.get(req.Key)
	if err != nil {
		return &cachepb.GetResponse{Found: false}, nil
	}
	return &cachepb.GetResponse{Value: val, Found: true}, nil
}

func (cn *CacheNode) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.SetResponse, error) {
	log.Printf("RPC Set key=%q value=%q", req.Key, req.Value)
	_ = cn.lru.set(req.Key, req.Value)
	return &cachepb.SetResponse{Success: true}, nil
}

func (cn *CacheNode) GetAllKeys(ctx context.Context, req *cachepb.GetAllKeysRequest) (*cachepb.GetAllKeysResponse, error) {
	keys := cn.lru.GetAllKeys()
	log.Printf("RPC GetAllKeys count=%d", len(keys))
	return &cachepb.GetAllKeysResponse{Keys: keys}, nil
}

func (cn *CacheNode) Delete(ctx context.Context, req *cachepb.DeleteRequest) (*cachepb.DeleteResponse, error) {
	log.Printf("RPC Delete key=%q", req.Key)
	success := cn.lru.delete(req.Key)
	return &cachepb.DeleteResponse{Success: success}, nil
}
