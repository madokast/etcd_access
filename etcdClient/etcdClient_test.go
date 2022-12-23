package etcdclient

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

const endpoint = "192.168.3.62:2379"

func TestSync(t *testing.T) {
	client := New(endpoint)
	defer client.Close()

	cntKey := "/user/count"

	done := make(chan bool)

	code := func(worker string) {
		for i := 0; i < 5; i++ {
			if client.Exists(cntKey) {
				cnt, _ := strconv.Atoi(client.Get(cntKey))
				t.Log(worker, "读取到", cnt, "写入", cnt+1)
				client.Put(cntKey, strconv.Itoa(cnt+1))
			} else {
				t.Log(worker, "不存在", cntKey, "写入", 1)
				client.Put(cntKey, strconv.Itoa(1))
			}
			time.Sleep(time.Second)
		}
		done <- true
	}

	go client.Sync("/sync/count", func() { code("w1") })
	go client.Sync("/sync/count", func() { code("w2") })
	go client.Sync("/sync/count", func() { code("w3") })

	<-done
	<-done
	<-done
}

func TestWatch(t *testing.T) {
	client := New(endpoint)
	defer client.Close()

	watchKey := "/user/zrx"

	go func() {
		client.Watch(watchKey, func(value string) {
			t.Log("新值", value)
		})
	}()

	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		client.Put(watchKey, fmt.Sprintf("w%d", i))
		time.Sleep(time.Second)
	}
}
