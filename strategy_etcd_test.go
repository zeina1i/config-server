package config_server

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

func Test_StrategyEtcd_changeKeyScenario(t *testing.T) {
	DefaultListenPeerURLs := "http://localhost:2382"
	DefaultListenClientURLs := "http://localhost:2381"
	lpurl, _ := url.Parse(DefaultListenPeerURLs)
	lcurl, _ := url.Parse(DefaultListenClientURLs)
	etcdCfg := embed.NewConfig()
	etcdCfg.LPUrls = []url.URL{*lpurl}
	etcdCfg.LCUrls = []url.URL{*lcurl}
	etcdCfg.Dir = "./data"
	e, err := embed.StartEtcd(etcdCfg)
	if err != nil {
		t.Fatal(err)
	}

	<-e.Server.ReadyNotify()
	t.Log("Server is ready!")

	etcdClientCfg := &clientv3.Config{
		Endpoints: []string{DefaultListenClientURLs},
	}
	etcd, err := NewStrategyEtcd(*etcdClientCfg)
	assert.NoError(t, err)

	err = etcd.Set("abcd", "efg")
	assert.NoError(t, err)

	val, err := etcd.Get("abcd")
	assert.NoError(t, err)
	assert.Equal(t, "efg", val)

	watchChan := etcd.Watch("abcd")
	assert.NoError(t, err)

	err = etcd.Set("abcd", "hij")
	assert.NoError(t, err)

	ret := <-watchChan
	assert.Equal(t, "hij", ret)

	err = etcd.Set("abcd", "klm")
	assert.NoError(t, err)

	ret = <-watchChan
	assert.Equal(t, "klm", ret)
}
