package main

import (
	"fmt"
	"net"
)

func main() {
	for i := 1; i <= 1024; i++ {
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// port is closed or filtered.
			continue
		}
		err = conn.Close()
		if err != nil {
			fmt.Printf("%d open\n", i)
			return
		}
		fmt.Printf("%d open\n", i)
	}
}
