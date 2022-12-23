package etcdclient

import (
	"context"
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type client struct {
	proxy *clientv3.Client
}

func New(endpoints ...string) client {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot assess estd server %v", err)
		os.Exit(1)
	}
	return client{proxy: c}
}

func (c client) Close() {
	c.proxy.Close()
}

func (c client) Get(key string) string {
	res, err := c.proxy.Get(context.Background(), key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot get key %v. %v", key, err)
		panic(err)
	}
	for _, kv := range res.Kvs {
		return string(kv.Value)
	}

	return ""
}

func (c client) Put(key, value string) (oldValue string) {
	pr, err := c.proxy.Put(context.Background(), key, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot put  %v %v. %v", key, value, err)
		panic(err)
	}
	if pr.PrevKv != nil {
		return string(pr.PrevKv.Value)
	}
	return ""
}

func (c client) GetAll(prefix string) map[string]string {
	res, err := c.proxy.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot get prefix %v. %v", prefix, err)
		panic(err)
	}
	kvMap := make(map[string]string)
	for _, kv := range res.Kvs {
		kvMap[string(kv.Key)] = string(kv.Value)
	}

	return kvMap
}

func (c client) Exists(key string) bool {
	gr, err := c.proxy.Get(context.Background(), key, clientv3.WithCountOnly())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot get key %v. %v", key, err)
		panic(err)
	}
	return gr.Count > 0
}

func (c client) PutTimed(key, value string, alive time.Duration) {
	lease := clientv3.NewLease(c.proxy)
	lgr, err := lease.Grant(context.Background(), int64(alive.Seconds()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot grant lease on %v %v. %v", key, value, err)
		panic(err)
	}
	c.proxy.Put(context.Background(), key, value, clientv3.WithLease(lgr.ID))
}

func (c client) Sync(key string, code func()) {
	s, err := concurrency.NewSession(c.proxy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ETCD] Cannot create seccion. %v", err)
		panic(err)
	}
	m := concurrency.NewMutex(s, key)
	m.Lock(context.Background())
	defer m.Unlock(context.Background())
	code()
}

func (c client) Watch(key string, react func(value string)) {
	wc := c.proxy.Watch(context.Background(), key)
	for wres := range wc {
		for _, e := range wres.Events {
			react(string(e.Kv.Value))
		}
	}
}
