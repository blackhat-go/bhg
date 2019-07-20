package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var storedHash = "$2a$10$Zs3ZwsjV/nF.KuvSUE.5WuwtDrK6UVXcBpQrH84V8q3Opg1yNdWLu"

func main() {
	var password string
	if len(os.Args) != 2 {
		log.Fatalln("Usage: bcrypt password")
	}
	password = os.Args[1]

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("hash = %s\n", hash)

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		log.Println("[!] Authentication failed")
		return
	}
	log.Println("[+] Authentication successful")
}
