imgInject
==========

## Objective
```
A venture into the world of steganography through means of
image formats while adhering to the actual specifications and
byte ordering. This utility currently supports payload injection
into arbitrary byte chunks. PNG format supported w/ more to come.
```


## Usage
```
$ ./imgInject -h
Example Usage: ./imgInject -i in.png -o out.png --inject --offset 0x85258 --payload 1234
Example Encode Usage: ./imgInject -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --encode --key secret
Example Decode Usage: ./imgInject -i encode.png -o decode.png --offset 0x85258 --decode --key secret
Flags: ./imginject {OPTION]...
  -i, --input string           Path to the original image file
  -o, --output string          Path to output the new image file
  -m, --meta                   Display the actual image meta details
  -s, --suppress               Suppress the chunk hex data (can be large)
      --offset string          The offset location to initiate data injection
      --inject                 Enable this to inject data at the offset location specified
      --payload string         Payload is data that will be read as a byte stream
      --type string[="rNDm"]   Type is the name of the Chunk header to inject (default "rNDm")
      --key string             The XOR key for payload
      --encode                 XOR encode the payload
      --decode                 XOR decode the payload
```

## Installation
```
# Installation
# -----------------------------------------------
# imgInject only uses Go core packages
# Go version 1.11.4
```

## Sample Run
```
$  ./imgInject -i images/battlecat.png -m -s
Valid PNG so let us continue!
---- Chunk # 1 ----
Chunk Offset: 0x08
Chunk Length: 13 bytes
Chunk Type: IHDR
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 99e4ec6
---- Chunk # 2 ----
Chunk Offset: 0x21
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: df8ffef3
---- Chunk # 3 ----
Chunk Offset: 0x802d
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 2142ea4a
---- Chunk # 4 ----
Chunk Offset: 0x10039
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 2c01ffc0
---- Chunk # 5 ----
Chunk Offset: 0x18045
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: e209a19
---- Chunk # 6 ----
Chunk Offset: 0x20051
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 529473dc
---- Chunk # 7 ----
Chunk Offset: 0x2805d
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 3c7ae73f
---- Chunk # 8 ----
Chunk Offset: 0x30069
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 5ec64e92
---- Chunk # 9 ----
Chunk Offset: 0x38075
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 6d91cf72
---- Chunk # 10 ----
Chunk Offset: 0x40081
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 540ee217
---- Chunk # 11 ----
Chunk Offset: 0x4808d
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 9951d57f
---- Chunk # 12 ----
Chunk Offset: 0x50099
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: d7fd0589
---- Chunk # 13 ----
Chunk Offset: 0x580a5
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 87acf32e
---- Chunk # 14 ----
Chunk Offset: 0x600b1
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: f93b6a8e
---- Chunk # 15 ----
Chunk Offset: 0x680bd
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: cd0c89ba
---- Chunk # 16 ----
Chunk Offset: 0x700c9
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 12f271e5
---- Chunk # 17 ----
Chunk Offset: 0x780d5
Chunk Length: 32768 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 4deb638d
---- Chunk # 18 ----
Chunk Offset: 0x800e1
Chunk Length: 20843 bytes
Chunk Type: IDAT
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: 9247d44d
---- Chunk # 19 ----
Chunk Offset: 0x85258
Chunk Length: 0 bytes
Chunk Type: IEND
Chunk Importance: Critical
Chunk Data: Suppressed
Chunk CRC: ae426082
```

## Developing
```
Alpha code under active development
```

## Contact
```
# Author: Chris Patten
# Contact (Email): chris[t.a.]stacktitan[t.o.d]com
# Contact (Twitter): @packetassailant
```


