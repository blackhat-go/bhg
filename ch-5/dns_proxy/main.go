package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
)

func parse(filename string) (map[string]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	records := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("%s is not a valid line", line)
		}
		records[parts[0]] = parts[1]
	}
	log.Println("records set to:")
	for k, v := range records {
		fmt.Printf("%s -> %s\n", k, v)
	}
	return records, scanner.Err()
}

func main() {
	records, err := parse("proxy.config")
	if err != nil {
		log.Fatalf("Error processing configuration file: %s\n", err.Error())
	}

	var recordLock sync.RWMutex
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		if len(req.Question) < 1 {
			dns.HandleFailed(w, req)
			return
		}
		name := req.Question[0].Name
		parts := strings.Split(name, ".")
		if len(parts) > 1 {
			name = strings.Join(parts[len(parts)-2:], ".")
		}
		recordLock.RLock()
		match, ok := records[name]
		recordLock.RUnlock()
		if !ok {
			dns.HandleFailed(w, req)
			return
		}
		resp, err := dns.Exchange(req, match)
		if err != nil {
			dns.HandleFailed(w, req)
			return
		}
		if err := w.WriteMsg(resp); err != nil {
			dns.HandleFailed(w, req)
			return
		}
	})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		for sig := range sigs {
			switch sig {
			case syscall.SIGUSR1:
				log.Println("SIGUSR1: reloading records")
				recordsUpdate, err := parse("proxy.config")
				if err != nil {
					log.Printf("Error processing configuration file: %s\n", err.Error())
				} else {
					recordLock.Lock()
					records = recordsUpdate
					recordLock.Unlock()
				}
			}
		}
	}()

	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
