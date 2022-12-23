package main

import (
	"fmt"
	"time"

	etcdclient "github.com/madokast/etcd_access/etcdClient"
)

func main() {
	client := etcdclient.New("192.168.3.62:2379")
	defer client.Close()

	watchKey := "/user/zrx"

	go func() {
		client.Watch(watchKey, func(value string) {
			fmt.Println("新值", value)
		})
	}()

	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		client.Put(watchKey, fmt.Sprintf("w%d", i))
		time.Sleep(time.Second)
	}
}
