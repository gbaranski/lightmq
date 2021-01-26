# LightMQ

# Abstract
LightMQ is Client Server messaging protocol. It is and will be lightweight and easy in client implementation. Protocol is originally intented to work with Internet of Things, where easy in making own implementation combined with simplicity and lightweightness is important.

- [LightMQ](#lightmq)
- [Abstract](#abstract)
- [Data representations](#data-representations)
  - [Bits](#bits)
  - [16-bit unsigned integer](#16-bit-unsigned-integer)
  - [Length prefixed string](#length-prefixed-string)
      - [Example](#example)
- [Opcodes](#opcodes)
- [Frame format](#frame-format)
  - [Example SEND frame](#example-send-frame)
- [Control frames](#control-frames)
  - [CONNECT](#connect)
    - [Payload structure](#payload-structure)
    - [ClientID](#clientid)
  - [CONNACK](#connack)
    - [Payload structure](#payload-structure-1)
    - [Return Code](#return-code)
  - [PING](#ping)
    - [Payload structure](#payload-structure-2)
  - [PONG](#pong)
    - [Payload structure](#payload-structure-3)
- [Data frames](#data-frames)
  - [SEND](#send)
    - [Payload structure](#payload-structure-4)
    - [Message ID](#message-id)
    - [Message Flags](#message-flags)
    - [Data](#data)
  - [SENDRESP](#sendresp)
    - [Payload structure](#payload-structure-5)
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

## Length prefixed string
Length prefixed strings means that length of a string is stored explicitly, before the actual text as a single byte. Length **MUST** be single byte value. String can be up to 255 bytes long.

#### Example
| Bit             | Value |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
| --------------- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Byte 0 - Length |   5   |   0   |   0   |   0   |   0   |   0   |   1   |   0   |   1   |
| Byte 1 - Char   |   H   |   0   |   1   |   0   |   0   |   1   |   0   |   0   |   0   |
| Byte 2 - Char   |   E   |   0   |   1   |   0   |   0   |   0   |   1   |   0   |   1   |
| Byte 3 - Char   |   L   |   0   |   1   |   0   |   0   |   1   |   1   |   0   |   0   |
| Byte 4 - Char   |   L   |   0   |   1   |   0   |   0   |   1   |   1   |   0   |   0   |
| Byte 5 - Char   |   O   |   0   |   1   |   0   |   0   |   1   |   1   |   1   |   1   |

# Opcodes
|         Name          |   Hex    |     Direction     | Description                                        |
| :-------------------: | :------: | :---------------: | -------------------------------------------------- |
|       Reserved        |   0x0    |       None        | For use in future                                  |
|  [CONNECT](#connect)  |   0x1    | Client -> Server  | Client request to connect to Server                |
|  [CONNACK](#connack)  |   0x2    | Server -> Client  | Server acknowledges connection request from Client |
|     [PING](#ping)     |   0x3    | Server <-> Client | Check if network connection is active              |
|     [PONG](#pong)     |   0x4    | Server <-> Client | Response to PING                                   |
|     [SEND](#send)     |   0x5    | Server <-> Client | Send messages with data                            |
| [SENDRESP](#sendresp) |   0x6    | Server <-> Client | Send response to message                           |
|       Reserved        | 0x7-0xFF |       None        | For use in future                                  |

# Frame format

|      Name      |  Length  |       Presence        | Description                                         |
| :------------: | :------: | :-------------------: | --------------------------------------------------- |
|     opcode     |  1 byte  |      Every frame      | Must be one of [registered opcodes](#opcodes)       |
| Payload length | 2 bytes  |      Every frame      | [16-bit unsigned integer](#16-bit-unsigned-integer) |
|    Payload     |    ^>    | If payload length > 0 |                                                     |
|    Reserved    | 0x7-0xFF |                       | None                                                |

## Example SEND frame

| Bit                                             |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
| ----------------------------------------------- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Byte 0 - [Packet type](#packet-type)(0x5)       |   0   |   0   |   0   |   0   |   0   |   1   |   0   |   1   |
| Byte 1 - [Payload size MSB](#payload-size)(0x2) |   0   |   0   |   0   |   0   |   0   |   0   |   0   |   0   |
| Byte 2 - [Payload size LSB](#payload-size)(0x0) |   0   |   0   |   0   |   0   |   0   |   0   |   0   |   0   |
| Byte 3 - [Payload byte 0](#payload)('H')        |   0   |   1   |   0   |   0   |   1   |   0   |   0   |   0   |
| Byte 4 - [Payload byte 1](#payload)('i')        |   0   |   1   |   1   |   0   |   1   |   0   |   0   |   1   |

# Control frames

## CONNECT
After a Network Connection is established by a Client to a Server, the first Packet sent from the Client to the Server **MUST** be a CONNECT Packet.

CONNECT packet may occur only once, second CONNECT packet **MUST** close connection.

### Payload structure

| Name                  | Size          |
| --------------------- | ------------- |
| ClientID size         | 1 byte        |
| [ClientID](#clientid) | Prefixed size |

<br/>

### ClientID
ClientID is [UTF-8 Length prefixed string](#utf8-string)

**MUST** be unique across different clients. 

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

## PING
The PING Packet is sent from a Client to the Server or from Server to Client. It can be used to test if the network connection is active.

### Payload structure

| Name                                  | Size    |
| ------------------------------------- | ------- |
| [Random ID](#16-bit-unsigned-integer) | 2 bytes |

## PONG
The PONG Packet is sent from a Client to the Server or from Server to Client. It acknowledges that it received PING. This packet type doesn't have payload.

### Payload structure

| Name                                                   | Size    |
| ------------------------------------------------------ | ------- |
| [Random ID, same as PING ID](#16-bit-unsigned-integer) | 2 bytes |

# Data frames

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

# TODO
- Figure out what exactly should be in [Variable header](#variable-header)
