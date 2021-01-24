# LightMQ

# Abstract
LightMQ is Client Server messaging protocol. It is and will be lightweight and easy in client implementation. Protocol is originally intented to work with Internet of Things, where easy in making own implementation combined with simplicity and lightweightness is important.

- [LightMQ](#lightmq)
- [Abstract](#abstract)
- [Data representations](#data-representations)
  - [Bits](#bits)
  - [16-bit unsigned integer](#16-bit-unsigned-integer)
  - [UTF8 String](#utf8-string)
      - [Example](#example)
- [Packet structure](#packet-structure)
  - [Packet type](#packet-type)
  - [Variable header](#variable-header)
  - [Signature](#signature)
  - [Payload size](#payload-size)
  - [Payload](#payload)
  - [Packet structure table](#packet-structure-table)
- [Packet types](#packet-types)
  - [CONNECT](#connect)
    - [Payload structure](#payload-structure)
    - [ClientID](#clientid)
    - [Challenge](#challenge)
  - [CONNACK](#connack)
    - [Payload structure](#payload-structure-1)
    - [Return Code](#return-code)
  - [SEND](#send)
    - [Payload structure](#payload-structure-2)
    - [Message ID](#message-id)
    - [Message Flags](#message-flags)
      - [ACK](#ack)
    - [Data](#data)
  - [SENDRESP](#sendresp)
    - [Payload structure](#payload-structure-3)
    - [Message ID](#message-id-1)
- [References](#references)
- [TODO](#todo)


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
```go
import (
    "encoding/binary"
)

bytes := []byte{0x20, 0x10} // 16-bit integer in bytes
value := binary.BigEndian.Uint16(bytes) // 8208
```

## UTF8 String
UTF-8 Length prefixed strings means that length of a string is stored explicitly, before the actual text. Length **MUST** be single byte value. String can be up to 256 bytes long.

#### Example
| Bit             | Value |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
| --------------- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Byte 0 - Length |   5   |   0   |   0   |   0   |   0   |   0   |   1   |   0   |   1   |
| Byte 1 - Char   |   H   |   0   |   1   |   0   |   0   |   1   |   0   |   0   |   0   |
| Byte 2 - Char   |   E   |   0   |   1   |   0   |   0   |   0   |   1   |   0   |   1   |
| Byte 3 - Char   |   L   |   0   |   1   |   0   |   0   |   1   |   1   |   0   |   0   |
| Byte 4 - Char   |   L   |   0   |   1   |   0   |   0   |   1   |   1   |   0   |   0   |
| Byte 5 - Char   |   O   |   0   |   1   |   0   |   0   |   1   |   1   |   1   |   1   |


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

Server **SHOULD** verify if signature is valid.

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
| Byte 69...65535 - [Payload](#payload)              |   X   |   X   |   X   |   X   |   X   |   X   |   X   |   X   |

<br/>

# Packet types

## CONNECT
After a Network Connection is established by a Client to a Server, the first Packet sent from the Client to the Server **MUST** be a CONNECT Packet.

CONNECT packet may occur only once, second CONNECT packet **MUST** close connection.


### Payload structure

| Name                  | Size          |
| --------------------- | ------------- |
| ClientID size         | 1 byte        |
| [ClientID](#clientid) | Prefixed size |
| Challenge             | 8 bytes       |

<br/>

### ClientID
ClientID is [UTF-8 Length prefixed string](#utf8-string)

**MUST** be unique across different clients. 

If Client with same ClientID already exists, server **MUST** disconnect old client.

### Challenge
Challenge is random byte string of 8 bytes in size, it is used to prevent hijacking signature of CONNECT message.

## CONNACK
The CONNACK Packet is the packet sent by the Server in response to a CONNECT Packet received from a Client. The first packet sent from the Server to the Client **MUST** be a CONNACK Packet.

### Payload structure
| Name                        | Size   |
| --------------------------- | ------ |
| [Return Code](#return-code) | 1 byte |

### Return Code

| Value    | Description                  |
| -------- | ---------------------------- |
| 0x0      | Connection Accepted          |
| 0x1      | Unsupported Protocol Version |
| 0x2      | Server unavailable           |
| 0x3      | Malformed payload            |
| 0x4      | Unauthorized                 |
| 0x5-0xFF | Reserved for future use      |


## SEND
A SEND Packet is sent from a Client to a Server or from Server to a Client to transport an Application Message.


### Payload structure
| Name                    | Size              |
| ----------------------- | ----------------- |
| [ID](#message-id)       | 2 bytes           |
| [Flags](#message-flags) | 1 byte            |
| [Data](#data)           | Up to 65463 bytes |

### Message ID
Random bytes used as correlation data.

### Message Flags
Special flags for message, currently all are reserved for future use

| Bit           |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
| ------------- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Message Flags |   R   |   R   |   R   |   R   |   R   |   R   |   R   |   R   |
|               |   0   |   0   |   0   |   0   |   0   |   0   |   0   |   0   |

*R - Reserved for future use*

#### ACK
Flag which tells if acknowledgement is expected

### Data
Data of the message


## SENDRESP
A SENDRESP Packet is sent from a Client to a Server or from Server to a Client as a response to SEND.


### Payload structure
| Name              | Size              |
| ----------------- | ----------------- |
| [ID](#message-id) | 2 bytes           |
| [Data](#data)     | Up to 65464 bytes |

### Message ID
Random bytes used as correlation data.

# References

1. Daniel J. Bernstein, Niels Duif, Tanja Lange, Peter Schwabe, Bo-Yin Yang, [High-speed high-security signatures](https://ed25519.cr.yp.to/ed25519-20110926.pdf) 
    > Public keys are 32 bytes, and signatures are 64 bytes.

# TODO
- Figure out what exactly should be in [Variable header](#variable-header)
