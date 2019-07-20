package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
)

func pad(buf []byte) []byte {
	// Assumes valid lengths. Should add additional checks.
	length := len(buf)
	padding := aes.BlockSize - (length % aes.BlockSize)
	if padding == 0 {
		padding = aes.BlockSize
	}
	padded := make([]byte, length+padding)
	copy(padded, buf)
	copy(padded[length:], bytes.Repeat([]byte{byte(padding)}, padding))
	return padded
}

func unpad(buf []byte) []byte {
	// Assume valid length and padding. Should add checks
	padding := int(buf[len(buf)-1])
	return buf[:len(buf)-padding]
}

func encrypt(plaintext, key []byte) ([]byte, error) {
	var (
		ciphertext []byte
		iv         []byte
		block      cipher.Block
		mode       cipher.BlockMode
		err        error
	)

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	iv = make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalln(err)
	}

	mode = cipher.NewCBCEncrypter(block, iv)

	plaintext = pad(plaintext)
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
	var (
		plaintext []byte
		iv        []byte
		block     cipher.Block
		mode      cipher.BlockMode
		err       error
	)

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Invalid ciphertext length: too short")
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid ciphertext length: not a multiple of blocksize")
	}

	iv = ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	mode = cipher.NewCBCDecrypter(block, iv)
	plaintext = make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext = unpad(plaintext)

	return plaintext, nil
}

func main() {
	var (
		err        error
		plaintext  []byte
		ciphertext []byte
		key        []byte
	)

	key = make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		log.Fatalln(err)
	}

	plaintext = []byte("4321123456789090")
	if ciphertext, err = encrypt(plaintext, key); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("key        = %x\n", key)
	fmt.Printf("ciphertext = %x\n", ciphertext)

	if plaintext, err = decrypt(ciphertext, key); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("plaintext  = %s\n", plaintext)
}
