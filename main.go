package main

import (
	"encoding/json"
	"fmt"
	"net"
)

var (
	channel chan map[string]interface{}
)

func main() {
	var (
		err  error
		addr = "127.0.0.1:12315"
	)
	if err = initConfig(); err != nil {
		panic(err)
	}
	if err = startES(); err != nil {
		panic(err)
	}
	channel = make(chan map[string]interface{}, 1024)
	go listenChan()
	udps, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		panic(err)
	}

	conn, err := net.ListenUDP("udp", udps)
	if err != nil {
		fmt.Println("Error listening:", err)
		panic(err)
	}
	defer conn.Close()
	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	var (
		n   int
		err error
	)
	data := make([]byte, 1024)
	udpData := make(map[string]interface{})
	n, _, err = conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}
	err = json.Unmarshal(data[:n], &udpData)
	if err != nil {
		fmt.Println("json.Unmarshal() err !", err)
		return
	}
	channel <- udpData
	fmt.Println("read from server:", udpData)
}
