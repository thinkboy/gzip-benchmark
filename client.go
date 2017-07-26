package main

// Start Commond eg: ./push_room 3000 60 http://127.0.0.1:8088/oldGzip
// 参数1: QPS
// 参数2: goroutine并发数
// 参数3: 被压测url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	count     = int64(0)
	beginTime int64
)

func main() {
	qps, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	rountineNum, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	qpsPerRoutine := qps / rountineNum                           //每个goroutine负担的qps数量
	sleePerRoutine := time.Second / time.Duration(qpsPerRoutine) //每个goroutine需要sleep的时间

	urlStr := os.Args[3]

	beginTime = time.Now().Unix()

	go coun()

	for i := 0; i < rountineNum; i++ {
		go run(urlStr, sleePerRoutine)
	}
	time.Sleep(9999 * time.Hour)
}

func coun() {
	for {
		time.Sleep(10 * time.Second)
		fmt.Printf("QPS: %d\n", count/(time.Now().Unix()-beginTime))
	}

}

func run(urlStr string, inter time.Duration) {
	for {
		t := get(urlStr) // 返回值为http耗时补齐
		time.Sleep(inter - t)
	}
}

func get(urlStr string) time.Duration {
	bTime := time.Now()
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if len(b) != 32763 {
		panic("压缩内容是不是变了？变了的话改下这里的返回长度判断")
	}

	atomic.AddInt64(&count, 1)
	return time.Now().Sub(bTime)
}
