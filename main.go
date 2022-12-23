package main

import (
	"context"
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.3.62:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("Cannot access etcd", err)
		os.Exit(1)
	}

	defer etcdClient.Close()

	res, err := etcdClient.Get(context.Background(), "/user/mdk")

	if err != nil {
		fmt.Println("Cannot get key", err)
	}

	for _, kv := range res.Kvs {
		fmt.Println(string(kv.Key), string(kv.Value))
	}
}
