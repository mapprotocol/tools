package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func BlockNumber() {
	q := 0
	all := 10000
	startTime := time.Now().Unix()
	var waitMutx sync.WaitGroup
	waitMutx.Add(all)
	for i := 0; i < all; i++ {
		go func(params int) {
			conn, _ := dialConn()
			_, error := conn.BlockNumber(context.Background())
			if error != nil {
				q++
				//query eth_blockNumber err  Post "http://119.8.165.158:8082": read tcp 192.168.10.201:62584->119.8.165.158:8082: wsarecv: An existing connection was forcibly closed by the remote host.
				//wsarecv:远程主机强制关闭了现有连接。
				fmt.Println(params, "   query eth_blockNumber err ", error)
				waitMutx.Done()
				return
			}
			waitMutx.Done()
		}(i)
	}
	waitMutx.Wait()
	num0, err := strconv.ParseFloat(fmt.Sprintf("%.3f", (float64(all)-float64(q))/float64(all)), 64) // 保留2位小数
	if err != nil {
		fmt.Println(err)
		return
	}
	endTime := time.Now().Unix()
	fmt.Println("Success rate: ", num0*100, "%", "  err num:", q, "  timeduring:", endTime-startTime)
	// 10000次 21次失败  成功率99.8%
	// 5000次  8次失败   成功率99.8%
	// 4800次  12次失败  成功率99.8%
	// 4600次  3次失败   成功率99.9%
	// 4500次  0次失败   成功率100%
	// 4000次  0次失败   成功率100%
	// 3000次  0次失败   成功率100%
	// 1000次  0次失败   成功率100%

}

//  4k * 1,000,000 = 4,000,000k ≈ 4G内存
func SuggestGasPrice() {
	q := 0
	all := 5000
	startTime := time.Now().Unix()
	var waitMutx sync.WaitGroup
	waitMutx.Add(all)
	for i := 0; i < all; i++ {
		go func(params int) {
			conn, _ := dialConn()
			_, error := conn.SuggestGasPrice(context.Background())
			if error != nil {
				q++
				//eth_gasPrice
				//Post "http://119.8.165.158:8082": dial tcp 119.8.165.158:8082:
				//connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.
				//connectex：每个套接字地址（协议/网络地址/端口）通常只允许使用一次。
				//嵌套地址不够用
				fmt.Println(params, "   query SuggestGasPrice err ", error)
				waitMutx.Done()
				return
			}
			waitMutx.Done()
		}(i)
	}
	waitMutx.Wait()
	num0, err := strconv.ParseFloat(fmt.Sprintf("%.3f", (float64(all)-float64(q))/float64(all)), 64) // 保留2位小数
	if err != nil {
		fmt.Println(err)
		return
	}
	endTime := time.Now().Unix()
	fmt.Println("Success rate: ", num0*100, "%", "  err num:", q, "  timeduring:", endTime-startTime)
}
