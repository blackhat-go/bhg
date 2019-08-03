package utils

func encodeDecode(input []byte, key string) []byte {
	var bArr = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		bArr[i] += input[i] ^ key[i%len(key)]
	}
	return bArr
}

//XorEncode returns encoded byte array
func XorEncode(decode []byte, key string) []byte {
	return encodeDecode(decode, key)
}

//XorDecode returns decoded byte array
func XorDecode(encode []byte, key string) []byte {
	return encodeDecode(encode, key)
}
