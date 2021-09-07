package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	for i := 1; i <= 1024; i++ {
		time.Sleep(100 * time.Millisecond)
		go func(j int) {
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Println("[+]Port -> ", j)
		}(i)
	}
}
