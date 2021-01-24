# LightMQ

# Abstract
LightMQ is Client Server messaging protocol. It is and will be lightweight and easy in client implementation. Protocol is originally intented to work with Internet of Things, where easy in making own implementation combined with simplicity and lightweightness is important.


# Data representations
## Bits
Bits in a byte are labeled 7 through 0. Bit number 7 is the most significant bit, the least significant bit is assigned bit number 0.

## 16-bit unsigned integer
16-bit unsigned integers are in [big-endian](https://en.wikipedia.org/wiki/Endiannes) order, that means high order byte is [MSB(Most significant bit)](https://en.wikipedia.org/wiki/Bit_numbering#Most_significant_bit) and the low order byte is [LSB(Least significant bit)](https://en.wikipedia.org/wiki/Bit_numbering#Least_significant_bit). This data representation allows to hold values in following range [0-65535], that comes from:

# Packet structure
## Fixed header

First 8 bits of EVERY packet data. Used to identify packet type. Represents 8-bit unsigned value, it can be one of following:

| Name                | Value | Direction        | Description                                        |
| ------------------- | ----- | ---------------- | -------------------------------------------------- |
|                     | 0     |                  |                                                    |
| [CONNECT](#connect) | 1     | Client -> Server | Client request to connect to Server                |
| [CONNACK](#connack) | 2     | Server -> Client | Server acknowledges connection request from Client |
|                     | 3     |                  |                                                    |
|                     | 4     |                  |                                                    |
|                     | 5     |                  |                                                    |
|                     | 6     |                  |                                                    |
|                     | 7     |                  |                                                    |
|                     | 8     |                  |                                                    |
|                     | 9     |                  |                                                    |

## Variable header

First 8 bits AFTER [Fixed header](#fixed-header) of EVERY packet data. Used to describe the packet

## Signature
Digital signature created using [Ed25519 scheme](https://en.wikipedia.org/wiki/EdDSA). Signatures are 512 bits(64 bytes) in size <sup>[1](#references)</sup>. They're raw bytes without any encoding

