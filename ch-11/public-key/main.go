package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {
	var (
		err                                              error
		privateKey                                       *rsa.PrivateKey
		publicKey                                        *rsa.PublicKey
		message, plaintext, ciphertext, signature, label []byte
	)

	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		log.Fatalln(err)
	}
	publicKey = &privateKey.PublicKey

	label = []byte("")
	message = []byte("Some super secret message, maybe a session key even")
	ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, label)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, label)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Plaintext: %s\n", plaintext)

	h := sha256.New()
	h.Write(message)
	signature, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil), nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Signature: %x\n", signature)

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, h.Sum(nil), signature, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Signature verified")
}
