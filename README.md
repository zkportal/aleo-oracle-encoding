# aleo-oracle-encoding

Golang package with functions for encoding oracle data for Aleo.

Since Leo language is missing features such as:
- strings,
- variable-length arrays,
- arrays bigger than 32 elements,
- structs with more than 32 fields,
- number types bigger than 16 bytes,
the oracle data is encoded as a struct of 32 structs of 32 `u128`s to allow the most flexibility of interpreting the data and making asserts on it. Before the data can be formatted as a struct of structs,
it needs to be encoded and aligned in a way that makes all the data points distinguishable and human readable (if possible) when formatted. For example, if a number 200 is encoded as byte 200, and is followed by other data without any padding,
then, after formatting, you will not get a struct field `200u128`. Instead, this number 200 must be encoded and aligned in a way that it becomes `200u128` as it's own separate field of the struct. The same should happen, when possible,
for all of the data points the oracle needs.

This package provides functions for encoding and aligning separate elements of different types without providing an implementation of encoding all of the oracle data as one blob.

All of encoding functions use alignment to 16 bytes, padded with zeroes if needed, unless specified otherwise. A buffer of 16 bytes is called a block.

This package intentionally doesn't specify how to encode everything needed for oracle contract and provides small functions for separate components instead, leaving some implementation details flexible.

For more information on meaning of all different components and values, see [Aleo Oracle SDK documentation](https://github.com/summitto/oracle-sdk/blob/main/documentation/doc.md).

## Encoding API

### `CreateMetaHeader` - encoding

Given the lengths of different data points it creates a 2-block meta header, which encodes the lengths. Every length integer is encoded using 2 little endian bytes.

| Byte positions | Data | Encoded value |
| --------- | ---- | ----- |
| 0-1 | attestations data length | variable |
| 2-3 | timestamp length | 8 |
| 4-5 | status code length | 8 |
| 6-7 | request method length | variable |
| 8-9 | response format length | 1 |
| 10-11 | URL length | variable |
| 12-13 | selector length | variable |
| 14-15 | encoding options length | 16 |
| 16-17 | length of request headers encoded with [`EncodeHeaders`](./README.md#encodeheaders---encoding) | variable |
| 18-19 | length of optional fields encoded with [`EncodeOptionalFields`](./README.md#encodeoptionalfields---encoding) | variable |
| 20-31 | reserved | 0 |

A meta header must be included in the full encoded blob in the known positions, otherwise decoding is very hard or impossible without knowing the original data.

### `EncodeAttestationData` - encoding

Encodes a given data string according to the format provided by the options. Can encode:
- strings
- positive floating-point numbers that fit into 64 bits
- unsigned integers up to 64 bits

#### Encoding a string

If a string is empty, will encode 1 block of zeroes, otherwise will encode the string as character codes and apply padding to 16 bytes.

| Byte positions | Data |
| --- | --- |
| 0-N | bytes of the string |
| N-M | padding to 16 with zeroes |

#### Encoding an integer

Parses a string as an unsigned 64-bit decimal integer and encodes it as 8 little endian bytes.

| Byte positions | Data |
| --- | --- |
| 0-7 | integer as 8 little endian bytes |
| 8-15 | reserved, 0 |

#### Encoding a float

Parses a string as a 64-bit decimal floating point number. Only positive numbers are supported at the moment. Exponent notation is not supported.

The encoder will return an `cannot parse float without losing information` error if it cannot parse the float, encode and then decode it back to the original string. It can happen because of possible bugs in the encoder or if the floating number is too big to be accurately represented in 8 bytes.

Encoding floats uses `EncodingOptions.Precision`. It must be a positive number less than or equal 12. Before encoding the parsed number is multiplied as `number * 10^precision`. If the resulting number is not an integer, then encoder returns an error `extracted value is more precise than given precision`. Essentially, the precision must be more or equal to the number of digits after the comma.

The resulting integer is encoded as 8 little-endian bytes.

| Byte positions | Data |
| --- | --- |
| 0-7 | float * (10^precision) as 8 little endian bytes |
| 8-15 | reserved, 0 |

### `EncodeResponseFormat` - encoding

Encodes response format type into 1 block where the first byte is 0 for JSON and 1 for HTML.

| Byte positions | Data |
| --- | --- |
| 0 | response format flag |
| 1-15 | reserved, 0 |

### `EncodeEncodingOptions` - encoding

Encodes encoding options object into one block of the following structure:

| Byte positions | Data | Comment |
| --- | --- | --- |
| 0 | value type | string=`0`, int=`1`, float=`2` |
| 1-7 | 0 | |
| 8-15 | encoding options precision | Little endian byte representation of the number. Due to the limit on Encoding options precision, this will always be only one byte with the actual value. If the value type is not float, this will be 0. |

### `EncodeHeaders` - encoding

Encodes a map of headers using the following components:

- 1 block of encoding meta header
- N blocks of encoded headers

Encoded structure:
| Byte positions | Data |
| --- | --- |
| 0-7 | number of headers, represented as 8 little endian bytes |
| 8-15 | number of blocks encoding the headers, represented as 8 little endian bytes |
| 16-N | encoded headers |

Header encoding uses the following pseudo code for every key-value pair:

```
result := []
for key, value in headers:
  entry := bytes("$key:$value")
  entryLen := len(entry) as uint16
  encodedHeader := entryLen + entry
  encodedHeader = addPadding(encodedHeader)
  result = result + encodedHeader
```

The encoded structure for `encodedHeader` therefore is:

| Byte positions | Data |
| --- | --- |
| 0-1 | `entryLen` |
| 2-N | `entry` |
| N-M | padding to 16 with zeroes |

The headers are sorted alphabetically by key. An empty map of headers is encoded into 1 block of zeroes.

### `EncodeOptionalFields` - encoding

Encodes optional notarization fields such as HTML result type (used only when response format is HTML), request content type (can only be used with POST request method) and request body (can only be used with POST request method).

This function always produces at least 4 blocks:
- meta header block, which encodes present optional fields and number of blocks they take. If no optional fields are present, the number of blocks will be set to 3
- 1 block of encoded HTML result type or 1 block of zeroes if HTML result type is not present
- N blocks of encoded request content type or 1 block of zeroes if request content type is not present
- N blocks of encoded request body or 1 block of zeroes if request body is not present

Meta header block structure:

| Byte positions | Data |
| --- | --- |
| 0 | field presence bitmask |
| 1-7 | reserved, 0 |
| 8-15 | number of blocks encoding the fields, will be at least 3 |

The field presence bitmask structure:

| Bit | Value |	Description |
| --- | --- | --- |
| 7 |	0 |	Reserved |
| 6 |	0 | Reserved |
| 5 |	0 | Reserved |
| 4 |	0 | Reserved |
| 3 |	0 | Reserved |
| 2 |	1 |	Bit is set when request body is present |
| 1 |	1 |	Bit is set when request content type is present |
| 0 |	1 |	Bit is set when HTML result type is present |

HTML result type encoding:

| Byte positions | Data |
| --- | --- |
| 0 | element=`1`, value=`2` |
| 1-15 | reserved, 0 |

Request content type and request body are encoded similarly using the following structure:
| Byte positions | Data |
| --- | --- |
| 0-7 | length of the string, represented as 8 little endian bytes |
| 8-15 | reserved, 0 |
| 16-N | bytes of the string |
| N-M | padding to 16 with zeroes |

## Decoding API

### `DecodeMetaHeader` - decoding

Decodes a meta header created with [`CreateMetaHeader`](./README.md#createmetaheader---encoding). The input buffer must be 2 blocks.

### `DecodeAttestationData` - decoding

Decodes attestation data created with [`EncodeAttestationData`](./README.md#encodeattestationdata---encoding). The input must be at least one block. This function requires encoding options and the length of the original string to decode the buffer. Encoding options must be parsed with [`DecodeEncodingOptions`](./README.md#decodeencodingoptions---decoding) or known before using this function. A meta header must be parsed with [`DecodeMetaHeader`](./README.md#decodemetaheader---decoding) to get the length of the original string or it must be known before calling this function.

If the provided length of the original string is not correct, the decoded data string may get trimmed.

### `DecodeResponseFormat` - decoding

Decodes response format created with [`EncodeResponseFormat`](./README.md#encoderesponseformat---encoding). The buffer must be 1 block.

### `DecodeEncodingOptions` - decoding

Decodes encoding options created with [`EncodeEncodingOptions`](./README.md#encodeencodingoptions---encoding). The buffer must be 1 block.

### `DecodeHeaders` - decoding

Decodes headers created with [`EncodeHeaders`](./README.md#encodeheaders---encoding). The buffer must be at least 1 block.

### `DecodeOptionalFields` - decoding

Decodes optional fields created with [`EncodeOptionalFields`](./README.md#encodeoptionalfields---encoding). The buffer must be at least 4 blocks.

## Utility API

### `NumberToBytes` - utility, no padding

Converts an unsigned 64-bit number to 8 little endian bytes

### `BytesToNumber` - utility, padding N/A

Converts a byte buffer to an unsigned 64-bit number by interpreting the first 8 bytes as little endian bytes of the aforementioned integer type. If the buffer is not
big enough, it will be padded internally on the right with zeroes.

### `BlockToNumbers` - utility

Converts 1 block to 2 unsigned 64-bit numbers using [`BytesToNumber`](./README.md#bytestonumber---utility-padding-na).

### `WriteWithPadding` - utility

Writes whatever buffer is provided to a position recorder (with an underlying `Writer`) and applies padding to 16 bytes to the buffer if needed.

Position recorder documentation can be found in the [positionRecorder package](./positionRecorder/README.md).
