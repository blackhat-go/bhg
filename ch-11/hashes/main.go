package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

var md5hash = "77f62e3524cd583d698d51fa24fdff4f"
var sha256hash = "95a5e1547df73abdd4781b6c9e55f3377c15d08884b11738c2727dbd887d4ced"

func main() {
	f, err := os.Open("wordlist.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		password := scanner.Text()
		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		if hash == md5hash {
			fmt.Printf("[+] Password found (MD5): %s\n", password)
		}

		hash = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		if hash == sha256hash {
			fmt.Printf("[+] Password found (SHA-256): %s\n", password)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
