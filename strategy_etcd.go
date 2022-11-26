package config_server

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type StrategyEtcd struct {
	client *clientv3.Client
}

func NewStrategyEtcd(etcdConfig clientv3.Config) (*StrategyEtcd, error) {
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	return &StrategyEtcd{client: client}, nil
}

func (etcd *StrategyEtcd) Watch(key string) <-chan string {
	watchChan := etcd.client.Watch(context.Background(), key)
	rec := make(chan string)

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
