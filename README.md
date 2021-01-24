# LightMQ

# Abstract
LightMQ is Client Server messaging protocol. It is and will be lightweight and easy in client implementation. Protocol is originally intented to work with Internet of Things, where easy in making own implementation combined with simplicity and lightweightness is important.


# Data representations
## Bits
Bits in a byte are labeled 7 through 0. Bit number 7 is the most significant bit, the least significant bit is assigned bit number 0.

## 16-bit unsigned integer
16-bit unsigned integers are in [big-endian](https://en.wikipedia.org/wiki/Endiannes) order, that means high order byte is [MSB(Most significant bit)](https://en.wikipedia.org/wiki/Bit_numbering#Most_significant_bit) and the low order byte is [LSB(Least significant bit)](https://en.wikipedia.org/wiki/Bit_numbering#Least_significant_bit). This data representation allows to hold values in following range [0-65535].

Calculating 16-bit unsigned integer from two bytes in C
```c
uint8_t bytes[2] = {0x20, 0x10}; // 16-bit integer in bytes
uint16_t value = (bytes[0] << 8) | bytes[1]; // 8208
```

Calculating 16-bit unsigned integer from two bytes in Go
```c
import (
    "encoding/binary"
)

bytes := []byte{0x20, 0x10} // 16-bit integer in bytes
value := binary.BigEndian.Uint16(bytes) // 8208
```

# Packet structure
## Packet type

**Position**: Starts at byte 0

**Size**: 1 byte(8 bits)

**MUST** exist in every packet data

Represented in 1 byte(8 bits). **MUST** be one of following:

| Name                | Dec | Bin        | Direction        | Description                                        |
| ------------------- | --- | ---------- | ---------------- | -------------------------------------------------- |
| [CONNECT](#connect) | 1   | `00000001` | Client -> Server | Client request to connect to Server                |
| [CONNACK](#connack) | 2   | `00000010` | Server -> Client | Server acknowledges connection request from Client |

If Client send invalid Packet type, Server **MAY** close the connection.

## Variable header

**Position**: Starts at byte 1

**Size**: 1 byte(8 bits)

**MUST** exist in every packet data

Used to describe the packet.

## Signature

**Position**: Starts at byte 2

**Size**: 64 bytes(512 bits)<sup>[1](#references)</sup>

**MUST** exist in every packet data

Digital signature created using [Ed25519 scheme](https://en.wikipedia.org/wiki/EdDSA) by signing the [payload](#payload) with private key, so server can verify [payload](#signature) using Client Public Key.


## Payload size
**Position**: Starts at byte 66

**Size**: 2 bytes(16 bits)

**MUST** exist in every packet data

Represented as [16-bit unsigned integer](#16-bit-unsigned-integer). Used to define size for [payload](#payload). Can be equal 0 meaning payload does not exist.

## Payload
**Position**: Starts at byte 68

**Size**: Defined by [Payload length](#payload-length)

## Packet structure table

Each packet have following structure:

| Bit                                                |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
| -------------------------------------------------- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Byte 0 - [Packet type](#packet-type)               |   X   |   X   |   X   |   X   |   X   |   X   |   X   |   X   |
| Byte 1 - [Variable header](#variable-header)       |   X   |   X   |   X   |   X   |   X   |   X   |   X   |   X   |
| Byte 2...66 - [Paylod signature](#variable-header) |   -   |   -   |   -   |   -   |   -   |   -   |   -   |   -   |
| Byte 67 - [Payload size MSB](#payload-size)        |   X   |   X   |   X   |   X   |   X   |   X   |   X   |   X   |
| Byte 68 - [Payload size LSB](#payload-size)        |   X   |   X   |   X   |   X   |   X   |   X   |   X   |   X   |



# Packet types

## CONNECT



## CONNACK


# References

1. Daniel J. Bernstein, Niels Duif, Tanja Lange, Peter Schwabe, Bo-Yin Yang, [High-speed high-security signatures](https://ed25519.cr.yp.to/ed25519-20110926.pdf) 
    > Public keys are 32 bytes, and signatures are 64 bytes.

# TODO
- Figure out what exactly should be in [Variable header](#variable-header)