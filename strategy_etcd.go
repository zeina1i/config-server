package config_server

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type StrategyEtcd struct {
	client           *clientv3.Client
	keyChannelMap    map[string]chan string
	keyCancelFuncMap map[string]context.CancelFunc
}

func NewStrategyEtcd(etcdConfig clientv3.Config) (*StrategyEtcd, error) {
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	return &StrategyEtcd{
		client:           client,
		keyCancelFuncMap: make(map[string]context.CancelFunc),
	}, nil
}

func (etcd *StrategyEtcd) Watch(key string) WatchChannel {
	ctx, cancel := context.WithCancel(context.Background())
	watchChan := etcd.client.Watch(ctx, key)
	rec := make(chan string)
	etcd.keyCancelFuncMap[key] = cancel
	etcd.keyChannelMap[key] = rec

	go func() {
		for resp := range watchChan {
			item := resp.Events[0].Kv.Value

			rec <- string(item)
		}
	}()

	return rec
}

func (etcd *StrategyEtcd) Get(key string) (string, error) {
	resp, err := etcd.client.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", err
	}

	return string(resp.Kvs[0].Value), nil
}

func (etcd *StrategyEtcd) Set(key string, val string) error {
	_, err := etcd.client.Put(context.Background(), key, val)
	if err != nil {
		return err
	}

	return nil
}

func (etcd *StrategyEtcd) Stop(key string) error {
	cancelFunc, ok := etcd.keyCancelFuncMap[key]
	if !ok {
		return fmt.Errorf("this key doesn't exist")
	}

	cancelFunc()
	delete(etcd.keyCancelFuncMap, key)

	channel, ok := etcd.keyChannelMap[key]
	if !ok {
		return fmt.Errorf("this channel doesn't exist")
	}

	close(channel)
	delete(etcd.keyChannelMap, key)

	return nil
}
