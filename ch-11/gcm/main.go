package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
)

var key = make([]byte, 32)

func encrypt(plaintext []byte) ([]byte, error) {
	var (
		ciphertext []byte
		nonce      []byte
		block      cipher.Block
		aead       cipher.AEAD
		err        error
	)

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	if aead, err = cipher.NewGCM(block); err != nil {
		return nil, err
	}

	nonce = make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln(err)
	}

	ciphertext = aead.Seal(nil, nonce, plaintext, nil)

	ciphertext = append(nonce, ciphertext...)
	return ciphertext, nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	var (
		plaintext []byte
		nonce     []byte
		block     cipher.Block
		aead      cipher.AEAD
		err       error
	)

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	if aead, err = cipher.NewGCM(block); err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("Invalid ciphertext length")
	}
	nonce = ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	if plaintext, err = aead.Open(nil, nonce, ciphertext, nil); err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	var (
		err        error
		plaintext  []byte
		ciphertext []byte
	)

	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		log.Fatalln(err)
	}

	plaintext = []byte("for queen and country")
	if ciphertext, err = encrypt(plaintext); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("ciphertext = %x\n", ciphertext)

	if plaintext, err = decrypt(ciphertext); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("plaintext = %s\n", plaintext)
}
