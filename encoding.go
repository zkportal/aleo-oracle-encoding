package aleoOracleEncoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/zkportal/aleo-oracle-encoding/positionRecorder"
)

var (
	ErrEncodingMetaHeaderInvalidSize              = errors.New("encoding general meta header requires a 2-block buffer")
	ErrIntValueParseFailure                       = errors.New("extracted value expected to be int but failed to parse as int")
	ErrFloatValueEncodingPrecisionTooBig          = errors.New("encoding precision is too big")
	ErrFloatNegativeUnsupported                   = errors.New("negative numbers are not supported for floats")
	ErrFloatValueDecimallessScientificUnsupported = errors.New("decimalless scientific notation is not supported for floats")
	ErrFloatValueScientificUnsupported            = errors.New("scientific notation is not supported for floats")
	ErrFloatValueParseFailure                     = errors.New("extracted value expected to be float but failed to parse as float")
	ErrFloatValueInfoLoss                         = errors.New("cannot parse float without losing information")
	ErrFloatValueNotEnoughPrecision               = errors.New("extracted value is more precise than given precision")
	ErrWritePaddingFailure                        = errors.New("failed to write to buffer")
	ErrValueEncodingUnknown                       = errors.New("unknown value type")
	ErrResponseFormatUnknown                      = errors.New("unknown response type")
	ErrHtmlResultTypeUnknown                      = errors.New("HTML result type is unknown")

	ErrDecodingInvalidMetaHeader = errors.New("invalid general meta header")

	ErrDecodingBufferTooShort        = errors.New("cannot decode buffer of unexpected size")
	ErrDecodingUnexpectedPadding     = errors.New("buffer contains unexpected padding")
	ErrDecodingAttestationImpossible = errors.New("cannot decode attestation data without encoding options information")

	ErrDecodingHeadersInvalidBlockHeader     = errors.New("buffer doesn't have meta header with encoding information")
	ErrDecodingHeadersCountLengthMismatch    = errors.New("buffer length doesn't match encoded headers length in meta header")
	ErrDecodingHeadersInvalidHeaderLength    = errors.New("encoded header length is bigger than buffer")
	ErrDecodingHeadersInvalidHeaderEncoding  = errors.New("encoded headers entry is not a valid encoded header - no separator found")
	ErrDecodingHeadersEmptyHeader            = errors.New("encoded headers entry is an empty header")
	ErrDecodingHeadersCountProcessedMismatch = errors.New("number of processed headers doesn't match number of headers in meta header")

	ErrDecodingOptionalsCountLengthMismatch      = errors.New("buffer length doesn't match encoded length in meta header")
	ErrDecodingOptionalsInvalidContentTypeLength = errors.New("encoded request content type length is bigger than buffer")
	ErrDecodingOptionalsInvalidBodyLength        = errors.New("encoded request body length is bigger than buffer")
	ErrDecodingOptionalsInvalidEncoding          = errors.New("could not parse the whole buffer")
)

const (
	TARGET_ALIGNMENT = 16

	RESPONSE_FORMAT_JSON_VALUE = 0 // value used for encoding response format for Aleo
	RESPONSE_FORMAT_HTML_VALUE = 1 // value used for encoding response format for Aleo

	ENCODING_OPTION_STRING_VALUE = 0 // value used for encoding encoding value format for Aleo
	ENCODING_OPTION_INT_VALUE    = 1 // value used for encoding encoding value format for Aleo
	ENCODING_OPTION_FLOAT_VALUE  = 2 // value used for encoding encoding value format for Aleo

	OPTIONAL_FIELDS_HEADER_HAS_HTML_RESULT_TYPE = 1 // bit flag used for encoding presence of HTML result type for Aleo
	OPTIONAL_FIELDS_HEADER_HAS_CONTENT_TYPE     = 2 // bit flag used for encoding presence of request content type for Aleo
	OPTIONAL_FIELDS_HEADER_HAS_REQUEST_BODY     = 4 // bit flag used for encoding presence of request body for Aleo

	HTML_RESULT_TYPE_ELEMENT_VALUE = 1 // value used for encoding HTML result type for Aleo
	HTML_RESULT_TYPE_VALUE_VALUE   = 2 // value used for encoding HTML result type for Aleo

	// Value encoding format options

	// Extracted value is a string
	ENCODING_OPTION_STRING = "string"
	// Extracted value is an unsigned decimal or hexadecimal integer up to 64 bits in size
	ENCODING_OPTION_INT = "int"
	// Extracted value is an unsigned floating point number up to 64 bits in size
	ENCODING_OPTION_FLOAT = "float"

	RESPONSE_FORMAT_HTML = "html"
	RESPONSE_FORMAT_JSON = "json"

	HTML_RESULT_TYPE_ELEMENT = "element"
	HTML_RESULT_TYPE_VALUE   = "value"

	ENCODING_OPTION_FLOAT_MAX_PRECISION = 12
)

type EncodingOptions struct {
	Value     string `json:"value"`
	Precision uint   `json:"precision"`
}

type ProofPositionalInfo struct {
	Data            positionRecorder.PositionInfo `json:"data"`
	Timestamp       positionRecorder.PositionInfo `json:"timestamp"`
	StatusCode      positionRecorder.PositionInfo `json:"statusCode"`
	Method          positionRecorder.PositionInfo `json:"method"`
	ResponseFormat  positionRecorder.PositionInfo `json:"responseFormat"`
	Url             positionRecorder.PositionInfo `json:"url"`
	Selector        positionRecorder.PositionInfo `json:"selector"`
	EncodingOptions positionRecorder.PositionInfo `json:"encodingOptions"`
	RequestHeaders  positionRecorder.PositionInfo `json:"requestHeaders"`
	OptionalFields  positionRecorder.PositionInfo `json:"optionalFields"`
}

// math.Pow works only with floats, which is inconvenient
func pow(base uint64, exponent uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < exponent; i++ {
		result *= base
	}
	return result
}

// returns an array of zeros so that when concatenated with arr the resulting array is padded to TARGET_ALIGNMENT.
// returns an empty array if arr doesn't need padding
func getPadding(arr []byte, alignment int) []byte {
	var paddingSize int
	overflow := len(arr) % alignment
	if overflow != 0 {
		paddingSize = alignment - overflow
	} else {
		paddingSize = 0
	}
	padding := make([]byte, paddingSize)
	return padding
}

// converts a number to 8 bytes in little-endian order
func NumberToBytes(number uint64) []byte {
	bytes := make([]byte, TARGET_ALIGNMENT/2)

	binary.LittleEndian.PutUint64(bytes, number)

	return bytes
}

// Converts 8 bytes in little-endian order to a number
func BytesToNumber(buf []byte) uint64 {
	b := make([]byte, len(buf))
	copy(b, buf)

	if len(b) < TARGET_ALIGNMENT/2 {
		padding := getPadding(b, TARGET_ALIGNMENT/2)
		b = append(b, padding...)
	}
	number := binary.LittleEndian.Uint64(b[:TARGET_ALIGNMENT/2])
	return number
}

// Converts 16-byte block to 2 uint64 numbers
func BlockToNumbers(buf []byte) []uint64 {
	if len(buf) != TARGET_ALIGNMENT {
		return nil
	}

	return []uint64{
		BytesToNumber(buf[0 : TARGET_ALIGNMENT/2]),
		BytesToNumber(buf[TARGET_ALIGNMENT/2 : TARGET_ALIGNMENT]),
	}
}

// creates a 2-block header, which encodes the byte length of all encoded elements
func CreateMetaHeader(header []byte, attestationDataLen, methodLen, urlLen, selectorLen, headersLen, optionalFieldsLen uint16) error {
	if len(header) != TARGET_ALIGNMENT*2 {
		return ErrEncodingMetaHeaderInvalidSize
	}

	// write attestation data length
	binary.LittleEndian.PutUint16(header[0:2], attestationDataLen)

	// write timestamp length
	binary.LittleEndian.PutUint16(header[2:4], 8) // timestamp is encoded as uint64 so it's always 8 bytes

	// write status code length
	binary.LittleEndian.PutUint16(header[4:6], 8) // status code is encoded as uint64 so it's always 8 bytes

	// write method length
	binary.LittleEndian.PutUint16(header[6:8], methodLen)

	// write response format length
	binary.LittleEndian.PutUint16(header[8:10], 1) // response format is always encoded as 1 byte

	// write URL length
	binary.LittleEndian.PutUint16(header[10:12], urlLen)

	// write selector length
	binary.LittleEndian.PutUint16(header[12:14], selectorLen)

	// write encoding options length
	binary.LittleEndian.PutUint16(header[14:16], 16) // encoding options are encoded as 2 uint64 numbers so it's always 16 bytes

	// write headers length
	binary.LittleEndian.PutUint16(header[16:18], headersLen)

	// write optional fields length
	binary.LittleEndian.PutUint16(header[18:20], optionalFieldsLen)

	return nil
}

type MetaHeader struct {
	AttestationDataLen int
	TimestampLen       int
	StatusCodeLen      int
	MethodLen          int
	ResponseFormatLen  int
	UrlLen             int
	SelectorLen        int
	EncodingOptionsLen int
	HeadersLen         int
	OptionalFieldsLen  int
}

func DecodeMetaHeader(header []byte) (parsedHeader *MetaHeader, err error) {
	if len(header) != TARGET_ALIGNMENT*2 {
		err = ErrDecodingInvalidMetaHeader
		return
	}

	return &MetaHeader{
		AttestationDataLen: int(binary.LittleEndian.Uint16(header[0:2])),
		TimestampLen:       int(binary.LittleEndian.Uint16(header[2:4])),
		StatusCodeLen:      int(binary.LittleEndian.Uint16(header[4:6])),
		MethodLen:          int(binary.LittleEndian.Uint16(header[6:8])),
		ResponseFormatLen:  int(binary.LittleEndian.Uint16(header[8:10])),
		UrlLen:             int(binary.LittleEndian.Uint16(header[10:12])),
		SelectorLen:        int(binary.LittleEndian.Uint16(header[12:14])),
		EncodingOptionsLen: int(binary.LittleEndian.Uint16(header[14:16])),
		HeadersLen:         int(binary.LittleEndian.Uint16(header[16:18])),
		OptionalFieldsLen:  int(binary.LittleEndian.Uint16(header[18:20])),
	}, nil
}

// parses the data string as a decimal 64-bit number and converts it to 8 bytes in little-endian order
func prepareDataAsInteger(data string) ([]byte, error) {
	var attestedNumber uint64
	var err error
	attestedNumber, err = strconv.ParseUint(data, 10, 64)
	if err != nil {
		log.Println("PrepareProofData: prepareDataAsInteger:", err)
		return nil, ErrIntValueParseFailure
	}

	return NumberToBytes(attestedNumber), nil
}

// parses the data string as a 64-bit float, multiplies it by 10^precision, and returns it as 8 bytes in little-endian order.
// If there are still fractions after multiplying by 10^precision, then returns an error
func prepareDataAsFloat(str string, precision uint) ([]byte, error) {
	if precision > ENCODING_OPTION_FLOAT_MAX_PRECISION {
		return nil, ErrFloatValueEncodingPrecisionTooBig
	}
	data := strings.ToLower(str)

	dotPos := strings.Index(data, ".")
	if dotPos == len(data)-1 {
		return nil, ErrFloatValueParseFailure
	}
	// trim all redundant zeroes if dot is present in number. this will not hurt us in decoding because we also
	// have the length of the original string encoded in the meta header
	if dotPos != -1 {
		data = strings.TrimRight(data, "0")
	}

	// if the dot is the last character, then all characters after the dot
	// were redundant zeroes. we can remove the dot
	if dotPos == len(data)-1 {
		data = strings.TrimRight(data, ".")
	} else if dotPos != -1 {
		fraction := data[dotPos+1:]
		// slice exponent to set precision
		if len(fraction) > int(precision) {
			return nil, ErrFloatValueNotEnoughPrecision
		}
	}

	// check for unsupported formats
	decimallessNotation := strings.Index(data, "p-")
	if decimallessNotation != -1 {
		return nil, ErrFloatValueDecimallessScientificUnsupported
	}

	exponentNotation := strings.Index(data, "e+")
	if exponentNotation != -1 {
		return nil, ErrFloatValueScientificUnsupported
	}
	exponentNotation = strings.Index(data, "e-")
	if exponentNotation != -1 {
		return nil, ErrFloatValueScientificUnsupported
	}

	hexExponentNotation := strings.Index(data, "p+")
	if hexExponentNotation != -1 {
		return nil, ErrFloatValueScientificUnsupported
	}

	negative := strings.Index(data, "-")
	if negative != -1 {
		return nil, ErrFloatNegativeUnsupported
	}

	bigNumber, _, err := big.ParseFloat(data, 10, 64, big.AwayFromZero)
	if err != nil {
		log.Println("PrepareProofData: prepareDataAsFloat:", err)
		return nil, ErrFloatValueParseFailure
	}
	bigNumber.SetMode(bigNumber.Mode())

	// the numbers are encoded as bytes of unsigned 64-bit integers so negative numbers
	// are not supported.
	// We do an explicit check for floats since there's no "unsigned float" but when parsing an integer
	// in prepareDataAsInteger we can parse a string as uint, which will return an error for negative numbers.
	if bigNumber.Sign() == -1 {
		return nil, ErrFloatNegativeUnsupported
	}

	bigMagnitude := new(big.Float).SetUint64(pow(10, uint64(precision)))
	bigAdjustedNumber := bigNumber.Mul(bigNumber, bigMagnitude)

	adjustedNumber, _ := bigAdjustedNumber.Float64()
	anotherAdjustedNumber, _ := bigAdjustedNumber.Uint64()

	// test recovery of all meaningful digits that will happen during decoding
	adjustedPrecision := int(precision)
	testStr := bigAdjustedNumber.Quo(bigAdjustedNumber, bigMagnitude).Text('f', int(precision))

	lenDiff := len(testStr) - len(data)

	if adjustedPrecision != 0 && len(data) != 0 && lenDiff != 0 {
		adjustedPrecision -= lenDiff
	}

	if data != bigAdjustedNumber.Text('f', adjustedPrecision) {
		return nil, ErrFloatValueInfoLoss
	}

	// we need to make sure that the provided precision is big enough
	// to cover all the digits after the comma in the extracted value
	if adjustedNumber != math.Floor(adjustedNumber) {
		return nil, ErrFloatValueNotEnoughPrecision
	}

	return NumberToBytes(anotherAdjustedNumber), nil
}

// writes data to the buffer, padding it to TARGET_ALIGNMENT bytes if needed.
// Returns the position info for the written aligned blocks
func WriteWithPadding(rec positionRecorder.PositionRecorder, data []byte) (*positionRecorder.PositionInfo, error) {
	padding := getPadding(data, TARGET_ALIGNMENT)
	buffer := make([]byte, 0, len(data)+len(padding))
	buffer = append(buffer, data...)
	buffer = append(buffer, padding...)

	n, err := rec.Write(buffer)
	if n != len(data)+len(padding) || err != nil {
		log.Printf("writeWithPadding: n=%d err=%s\n", n, err)
		return nil, ErrWritePaddingFailure
	}

	return rec.GetLastWrite(), nil
}

// Encodes data according to encoding options.
//
// If options.Value is "string", then data is encoded as character codes.
//
// If options.Value is "int", then data is parsed as a decimal or hexadecimal 64-bit number and encoded to 8 little endian bytes.
//
// If options.Value is "float", then data is parsed as a 64-bit float, then multiplied by 10^options.Precision, and encoded as 8 bytes in little-endian order.
// If there are still fractions after multiplying by 10^precision, then returns an error. If the data string is not equal to parsed string converted back to string,
// then you will be losing information, therefore an error is returned.
func EncodeAttestationData(data string, options *EncodingOptions) ([]byte, error) {
	var attestationDataBuffer []byte
	var err error

	switch options.Value {
	case ENCODING_OPTION_STRING:
		if data == "" {
			attestationDataBuffer = make([]byte, TARGET_ALIGNMENT)
		} else {
			attestationDataBuffer = []byte(data)
		}
	case ENCODING_OPTION_INT:
		attestationDataBuffer, err = prepareDataAsInteger(data)
	case ENCODING_OPTION_FLOAT:
		attestationDataBuffer, err = prepareDataAsFloat(data, options.Precision)
	default:
		err = ErrValueEncodingUnknown
	}
	if err != nil {
		return nil, err
	}

	padding := getPadding(attestationDataBuffer, TARGET_ALIGNMENT)
	attestationDataBuffer = append(attestationDataBuffer, padding...)
	return attestationDataBuffer, nil
}

func DecodeAttestationData(buf []byte, stringLen int, options *EncodingOptions) (string, error) {
	if len(buf) < TARGET_ALIGNMENT {
		return "", ErrDecodingBufferTooShort
	}

	if options == nil {
		return "", ErrDecodingAttestationImpossible
	}

	switch options.Value {
	case ENCODING_OPTION_STRING:
		if stringLen > len(buf) {
			return "", ErrDecodingBufferTooShort
		}
		return string(buf[:stringLen]), nil

	case ENCODING_OPTION_INT:
		number := BytesToNumber(buf[:TARGET_ALIGNMENT/2])
		return strconv.FormatUint(number, 10), nil

	case ENCODING_OPTION_FLOAT:
		number := BytesToNumber(buf[:TARGET_ALIGNMENT/2])
		float := new(big.Float).SetUint64(number).SetPrec(64).SetMode(big.ToNearestAway)
		magnitude := new(big.Float).SetUint64(pow(10, uint64(options.Precision)))

		float = float.Quo(float, magnitude)

		// since we know the length of the original string, we can figure out how many
		// redundant zeroes we are supposed to have
		adjustedPrecision := int(options.Precision)
		testStr := float.Text('f', int(options.Precision))

		lenDiff := len(testStr) - stringLen

		if adjustedPrecision != 0 && stringLen != 0 && lenDiff != 0 {
			adjustedPrecision -= lenDiff
		}

		return float.Text('f', adjustedPrecision), nil

	default:
		return "", ErrValueEncodingUnknown
	}
}

// Encodes response format type as 1 block. The first little-endian byte encodes the format type - 1 for HTML, 0 for JSON
func EncodeResponseFormat(format string) ([]byte, error) {
	buf := make([]byte, TARGET_ALIGNMENT)
	switch format {
	case RESPONSE_FORMAT_HTML:
		buf[0] = RESPONSE_FORMAT_HTML_VALUE
	case RESPONSE_FORMAT_JSON:
		buf[0] = RESPONSE_FORMAT_JSON_VALUE
	default:
		return nil, ErrResponseFormatUnknown
	}

	return buf, nil
}

// Decodes byte slice to the original format string
func DecodeResponseFormat(buf []byte) (string, error) {
	if len(buf) != TARGET_ALIGNMENT {
		return "", ErrDecodingBufferTooShort
	}

	switch buf[0] {
	case RESPONSE_FORMAT_HTML_VALUE:
		return RESPONSE_FORMAT_HTML, nil
	case RESPONSE_FORMAT_JSON_VALUE:
		return RESPONSE_FORMAT_JSON, nil
	default:
		return "", ErrResponseFormatUnknown
	}
}

// Encodes encoding options as 1 block. The first 8 little-endian bytes contain value type, where the first byte is the encoded value -
// 0 for string, 1 for int, 2 for float. If the encoded value type is float, then the second 8 bytes encode the floating point precision as little-endian bytes
// representing the number.
func EncodeEncodingOptions(options *EncodingOptions) ([]byte, error) {
	var valueTypeByte byte
	var precisionByte byte

	switch options.Value {
	case ENCODING_OPTION_STRING:
		valueTypeByte = ENCODING_OPTION_STRING_VALUE
		precisionByte = 0
	case ENCODING_OPTION_INT:
		valueTypeByte = ENCODING_OPTION_INT_VALUE
		precisionByte = 0
	case ENCODING_OPTION_FLOAT:
		if options.Precision > ENCODING_OPTION_FLOAT_MAX_PRECISION {
			return nil, ErrFloatValueEncodingPrecisionTooBig
		}
		valueTypeByte = ENCODING_OPTION_FLOAT_VALUE
		precisionByte = byte(options.Precision)
	default:
		return nil, ErrValueEncodingUnknown
	}

	valueBytes := NumberToBytes(uint64(valueTypeByte))
	precisionBytes := NumberToBytes(uint64(precisionByte))

	return append(valueBytes, precisionBytes...), nil
}

// Decodes byte slice to the original encoding options
func DecodeEncodingOptions(buf []byte) (*EncodingOptions, error) {
	if len(buf) != TARGET_ALIGNMENT {
		return nil, ErrDecodingBufferTooShort
	}

	valueTypeByte := buf[0]
	var precisionByte byte
	if valueTypeByte == ENCODING_OPTION_FLOAT_VALUE {
		precisionByte = buf[8]
	}

	switch valueTypeByte {
	case ENCODING_OPTION_STRING_VALUE:
		return &EncodingOptions{Value: ENCODING_OPTION_STRING, Precision: 0}, nil
	case ENCODING_OPTION_INT_VALUE:
		return &EncodingOptions{Value: ENCODING_OPTION_INT, Precision: 0}, nil
	case ENCODING_OPTION_FLOAT_VALUE:
		return &EncodingOptions{Value: ENCODING_OPTION_FLOAT, Precision: uint(precisionByte)}, nil
	default:
		return nil, ErrValueEncodingUnknown
	}
}

// encodes headers in the following format:
// 1 block - ((number of headers << 64) | number of blocks of headers)
// 2+ blocks - 2 bytes of "header:value" length + "header:value" + pad to TARGET_ALIGNMENT, repeat for all headers
// the headers are sorted alphabetically
func EncodeHeaders(headers map[string]string) []byte {
	buf := make([]byte, TARGET_ALIGNMENT, len(headers)*TARGET_ALIGNMENT+TARGET_ALIGNMENT)

	// collect keys first and sort them
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := headers[key]
		entry := []byte(fmt.Sprintf("%s:%s", key, val))
		lenBuf := make([]byte, 2)
		binary.LittleEndian.PutUint16(lenBuf, uint16(len(entry)))
		entry = append(lenBuf, entry...)

		padding := getPadding(entry, TARGET_ALIGNMENT)

		buf = append(buf, entry...)
		buf = append(buf, padding...)
	}

	numHeaders := uint64(len(headers))
	copy(buf[:TARGET_ALIGNMENT/2], NumberToBytes(numHeaders))

	// the first block is the one where we're writing length so we're not counting it
	numBlocks := uint64(len(buf)/TARGET_ALIGNMENT - 1)
	copy(buf[TARGET_ALIGNMENT/2:TARGET_ALIGNMENT], NumberToBytes(numBlocks))

	return buf
}

func DecodeHeaders(buf []byte) (map[string]string, error) {
	if len(buf) < TARGET_ALIGNMENT {
		return nil, ErrDecodingBufferTooShort
	}

	headers := make(map[string]string)

	// if there's only one block, then it's a block header with no content
	if len(buf) == TARGET_ALIGNMENT {
		return headers, nil
	}

	parsedBlockHeader := BlockToNumbers(buf[:TARGET_ALIGNMENT])
	if len(parsedBlockHeader) != 2 {
		return nil, ErrDecodingHeadersInvalidBlockHeader
	}

	headerCount := parsedBlockHeader[0]
	blockCount := parsedBlockHeader[1]

	// verify that the encoded block length + block header matches the buffer length
	if len(buf) != int(blockCount+1)*TARGET_ALIGNMENT {
		return nil, ErrDecodingHeadersCountLengthMismatch
	}

	headersProcessed := 0
	byteOffset := TARGET_ALIGNMENT
	for byteOffset < len(buf) {
		// read 2 bytes of header length and convert it to a number
		entryLenBuf := buf[byteOffset : byteOffset+2]
		byteOffset += 2

		// decode the length of an entry
		entryLen := int(binary.LittleEndian.Uint16(entryLenBuf))
		if byteOffset+entryLen > len(buf) {
			return nil, ErrDecodingHeadersInvalidHeaderLength
		}

		// get the entry
		entry := buf[byteOffset : byteOffset+entryLen]
		byteOffset += entryLen

		// an entry is formatted as "header:value", split around the first colon
		header, value, found := strings.Cut(string(entry), ":")
		if !found {
			return nil, ErrDecodingHeadersInvalidHeaderEncoding
		}
		if header == "" {
			return nil, ErrDecodingHeadersEmptyHeader
		}

		headers[header] = value
		headersProcessed += 1

		// check if this header entry was padded, skip the padding if it was
		currentAlignment := byteOffset % TARGET_ALIGNMENT
		if currentAlignment != 0 {
			paddingBytes := TARGET_ALIGNMENT - currentAlignment
			padding := buf[byteOffset : byteOffset+paddingBytes]
			expectedPadding := make([]byte, paddingBytes)
			// may be an overkill to verify that the padding used is made of zeroes?
			if !bytes.Equal(padding, expectedPadding) {
				return nil, ErrDecodingUnexpectedPadding
			}

			byteOffset += paddingBytes
		}
	}

	// verify that we've got the same number of headers as the encoding say
	if headersProcessed != int(headerCount) {
		return nil, ErrDecodingHeadersCountProcessedMismatch
	}

	return headers, nil
}

// Encodes fields that are optional for the oracle into blocks.
//
// The first block is a header - the first little endian byte is a bitmask, which encodes existence of HTML result type (1st bit),
// request content type (2nd bit), request body (3rd bit). The last 8 little-endian bytes are a little-endian byte representation of
// the number of blocks following the header. There will always be at least 3 blocks after the header - zero-byte blocks encoding the lengths of the (non-existent) components.
//
// The header is followed by the following content:
//
// 1. 1 block encoding HTML result type. The first little endian byte encodes the value - 1 for "element", 2 for "value".
// If there's no HTML result type, then the whole block is 0.
//
// 2. At least 1 block encoding request content type. The first 8 little endian bytes encode the number of the following blocks encoding the actual content type as character
// codes. If there is no content type, there's 1 block of 0, followed by 0 blocks of content.
//
// 3. At least 1 block encoding request body. The first 8 little endian bytes encode the number of the following blocks encoding the actual request body as character codes.
// If there is no request body, there's 1 block of 0, followed by 0 blocks of content.
func EncodeOptionalFields(htmlResultType, requestContentType, requestBody *string) ([]byte, error) {
	header := make([]byte, TARGET_ALIGNMENT)
	var htmlResultTypeBuf, contentTypeBuf, requestBodyBuf []byte

	// if there's HTML result type, set the byte in the header,
	// encode the type in 1 block.
	// if there's no HTML type, write one block of zeros.
	htmlResultTypeBuf = make([]byte, TARGET_ALIGNMENT)
	if htmlResultType != nil {
		header[0] |= OPTIONAL_FIELDS_HEADER_HAS_HTML_RESULT_TYPE

		switch *htmlResultType {
		case HTML_RESULT_TYPE_ELEMENT:
			htmlResultTypeBuf[0] = HTML_RESULT_TYPE_ELEMENT_VALUE
		case HTML_RESULT_TYPE_VALUE:
			htmlResultTypeBuf[0] = HTML_RESULT_TYPE_VALUE_VALUE
		default:
			return nil, ErrHtmlResultTypeUnknown
		}
	}

	// if there's a request content type, set the bit in the header,
	// encode the content type as a string with padding at the end, prepend with 1 block of
	// length.
	// if there's no content type, write one block of zeroes.
	contentTypeBuf = make([]byte, TARGET_ALIGNMENT)
	if requestContentType != nil {
		header[0] |= OPTIONAL_FIELDS_HEADER_HAS_CONTENT_TYPE

		data := []byte(*requestContentType)
		padding := getPadding(data, TARGET_ALIGNMENT)

		// write the string length to the request content type meta header
		copy(contentTypeBuf[:TARGET_ALIGNMENT/2], NumberToBytes(uint64(len(data))))
		contentTypeBuf = append(contentTypeBuf, data...)
		contentTypeBuf = append(contentTypeBuf, padding...)
	}

	// if there's a request body, set the bit in the header,
	// encode the request body as a string with padding at the end, prepend with 1 block of
	// length.
	// if there's no request body, write one block of zeroes.
	requestBodyBuf = make([]byte, TARGET_ALIGNMENT)
	if requestBody != nil {
		header[0] |= OPTIONAL_FIELDS_HEADER_HAS_REQUEST_BODY

		data := []byte(*requestBody)
		padding := getPadding(data, TARGET_ALIGNMENT)

		// write the string length to the request body meta header
		copy(requestBodyBuf[:TARGET_ALIGNMENT/2], NumberToBytes(uint64(len(data))))
		requestBodyBuf = append(requestBodyBuf, data...)
		requestBodyBuf = append(requestBodyBuf, padding...)
	}

	dataSize := len(htmlResultTypeBuf) + len(contentTypeBuf) + len(requestBodyBuf)

	// now that we know how much data we actually have, we can write the block length to the header
	copy(header[TARGET_ALIGNMENT/2:TARGET_ALIGNMENT], NumberToBytes(uint64(dataSize/TARGET_ALIGNMENT)))

	result := make([]byte, 0, TARGET_ALIGNMENT+dataSize)

	result = append(result, header...)
	result = append(result, htmlResultTypeBuf...)
	result = append(result, contentTypeBuf...)
	result = append(result, requestBodyBuf...)

	return result, nil
}

func DecodeOptionalFields(buf []byte) (htmlResultType, requestContentType, requestBody *string, err error) {
	if len(buf) < 4*TARGET_ALIGNMENT {
		err = ErrDecodingBufferTooShort
		return
	}

	header := buf[0:TARGET_ALIGNMENT]
	blockCount := BytesToNumber(header[TARGET_ALIGNMENT/2 : TARGET_ALIGNMENT])

	// check that the header encodes the correct number of of the following blocks
	if (blockCount+1)*TARGET_ALIGNMENT != uint64(len(buf)) {
		err = ErrDecodingOptionalsCountLengthMismatch
		return
	}

	// no optional fields are present
	if header[0] == 0 {
		return
	}

	hasHtmlResultType := (header[0] & OPTIONAL_FIELDS_HEADER_HAS_HTML_RESULT_TYPE) == OPTIONAL_FIELDS_HEADER_HAS_HTML_RESULT_TYPE
	hasRequestContentType := (header[0] & OPTIONAL_FIELDS_HEADER_HAS_CONTENT_TYPE) == OPTIONAL_FIELDS_HEADER_HAS_CONTENT_TYPE
	hasRequestBody := (header[0] & OPTIONAL_FIELDS_HEADER_HAS_REQUEST_BODY) == OPTIONAL_FIELDS_HEADER_HAS_REQUEST_BODY

	blockOffset := 1 * TARGET_ALIGNMENT

	if hasHtmlResultType {
		// html result type is encoded in the first byte of the first block after the meta header
		switch buf[blockOffset] {
		case HTML_RESULT_TYPE_ELEMENT_VALUE:
			htmlResultType = new(string)
			*htmlResultType = HTML_RESULT_TYPE_ELEMENT
		case HTML_RESULT_TYPE_VALUE_VALUE:
			htmlResultType = new(string)
			*htmlResultType = HTML_RESULT_TYPE_VALUE
		default:
			err = ErrHtmlResultTypeUnknown
			return
		}
	}
	blockOffset += 1 * TARGET_ALIGNMENT

	if hasRequestContentType {
		contentTypeHeader := buf[blockOffset : blockOffset+TARGET_ALIGNMENT]
		// parse the length of the content
		contentTypeLen := BytesToNumber(contentTypeHeader[:TARGET_ALIGNMENT/2])
		// check if encoded length makes sense. We add one more block to the length here to account for the request body header block, which comes after content type
		if blockOffset+int(contentTypeLen)+TARGET_ALIGNMENT > len(buf) {
			err = ErrDecodingOptionalsInvalidContentTypeLength
			return
		}

		// skip the content type header and read the content type, then figure out if there was padding and skip it
		contentType := buf[blockOffset+TARGET_ALIGNMENT : blockOffset+TARGET_ALIGNMENT+int(contentTypeLen)]
		padding := getPadding(contentType, TARGET_ALIGNMENT)

		requestContentType = new(string)
		*requestContentType = string(contentType)
		blockOffset += int(contentTypeLen) + len(padding)
	}
	blockOffset += 1 * TARGET_ALIGNMENT

	if hasRequestBody {
		requestBodyHeader := buf[blockOffset : blockOffset+TARGET_ALIGNMENT]
		// parse the length of the content
		requestBodyLen := BytesToNumber(requestBodyHeader[:TARGET_ALIGNMENT/2])
		// check if encoded length makes sense. This is the last header, so no need to add anything like we did for content type.
		if blockOffset+int(requestBodyLen) > len(buf) {
			err = ErrDecodingOptionalsInvalidContentTypeLength
			return
		}

		// skip the content type header and read the content type, then figure out if there was padding and skip it
		body := buf[blockOffset+TARGET_ALIGNMENT : blockOffset+TARGET_ALIGNMENT+int(requestBodyLen)]
		padding := getPadding(body, TARGET_ALIGNMENT)

		requestBody = new(string)
		*requestBody = string(body)
		blockOffset += int(requestBodyLen) + len(padding)
	}
	blockOffset += 1 * TARGET_ALIGNMENT

	if blockOffset != len(buf) {
		err = ErrDecodingOptionalsInvalidEncoding
		return
	}

	return
}
