package checkpointstore

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

/*
创建一个类，实现接口

	type CheckPointStore interface {
		Get(ctx context.Context, checkPointID string) ([]byte, bool, error)
		Set(ctx context.Context, checkPointID string, checkPoint []byte) error
	}
*/

type InMemoryStore struct {
	mem map[string][]byte //  基于本地内存 存储KV，在生产环境中最好基于Redis来存储KV，因为一旦go进程挂了，数据就丢失了
}

func (i *InMemoryStore) Set(ctx context.Context, key string, value []byte) error {
	// log.Printf("store %s", key)
	i.mem[key] = value
	return nil
}

func (i *InMemoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	// log.Printf("read %s", key)
	v, ok := i.mem[key]
	return v, ok, nil
}

func NewInMemoryStore() compose.CheckPointStore {
	return &InMemoryStore{
		mem: map[string][]byte{},
	}
}
