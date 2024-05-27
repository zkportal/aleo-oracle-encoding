package aleoOracleEncoding

import (
	"aleo-oracle-encoding/positionRecorder"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_pow(t *testing.T) {
	type args struct {
		base     uint
		exponent uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "basic",
			args: args{
				base:     10,
				exponent: 2,
			},
			want: 100,
		},
		{
			name: "basic 2",
			args: args{
				base:     2,
				exponent: 10,
			},
			want: 1024,
		},
		{
			name: "exponent 1",
			args: args{
				base:     2,
				exponent: 1,
			},
			want: 2,
		},
		{
			name: "exponent 1",
			args: args{
				base:     2,
				exponent: 0,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pow(uint64(tt.args.base), uint64(tt.args.exponent)); got != uint64(tt.want) {
				t.Errorf("pow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPadding(t *testing.T) {
	type args struct {
		arr       []byte
		alignment int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "pad empty",
			args: args{
				arr:       []byte{},
				alignment: 16,
			},
			want: []byte{},
		},
		{
			name: "pad short",
			args: args{
				arr:       []byte{1, 1, 1, 1},
				alignment: 16,
			},
			want: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "pad long",
			args: args{
				arr:       []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				alignment: 16,
			},
			want: []byte{0, 0},
		}, {
			name: "no pad",
			args: args{
				arr:       []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				alignment: 16,
			},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPadding(tt.args.arr, tt.args.alignment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NumberToBytes(t *testing.T) {
	tests := []struct {
		name   string
		number uint64
		want   []byte
	}{
		{
			name:   "number 1",
			number: 200,
			want:   []byte{200, 0, 0, 0, 0, 0, 0, 0},
		}, {
			name:   "number 2",
			number: 64250,
			want:   []byte{0xfa, 0xfa, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "number 3",
			number: 0xdeadbeefdeadbeef,
			want:   []byte{0xef, 0xbe, 0xad, 0xde, 0xef, 0xbe, 0xad, 0xde},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberToBytes(tt.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BytesToNumber(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		want uint64
	}{
		{
			name: "number 1",
			buf:  []byte{200, 0, 0, 0, 0, 0, 0, 0},
			want: 200,
		}, {
			name: "number 2",
			buf:  []byte{0xfa, 0xfa, 0, 0, 0, 0, 0, 0},
			want: 64250,
		},
		{
			name: "number 3",
			buf:  []byte{0xef, 0xbe, 0xad, 0xde, 0xef, 0xbe, 0xad, 0xde},
			want: 0xdeadbeefdeadbeef,
		},
		{
			name: "short buffer",
			buf:  []byte{0xde, 0xad},
			want: 44510,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToNumber(tt.buf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesToNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockToNumbers(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{
			name: "valid input",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: []uint64{0, 0},
		},
		{
			name: "short input",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "nil input",
			args: args{
				buf: nil,
			},
			want: nil,
		},
		{
			name: "long input",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "valid input 2",
			args: args{
				buf: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			},
			want: []uint64{506097522914230528, 1084818905618843912},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BlockToNumbers(tt.args.buf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlockToNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMetaHeader(t *testing.T) {
	header := make([]byte, TARGET_ALIGNMENT*2)

	type args struct {
		header             []byte
		attestationDataLen uint16
		methodLen          uint16
		urlLen             uint16
		selectorLen        uint16
		headersLen         uint16
		optionalFieldsLen  uint16
	}
	tests := []struct {
		name       string
		args       args
		wantHeader []byte
		wantErr    bool
	}{
		{
			name: "nil buffer",
			args: args{
				header:             nil,
				attestationDataLen: 1,
				methodLen:          1,
				urlLen:             1,
				selectorLen:        1,
				headersLen:         1,
				optionalFieldsLen:  1,
			},
			wantHeader: nil,
			wantErr:    true,
		},
		{
			name: "short buffer",
			args: args{
				header:             make([]byte, TARGET_ALIGNMENT*2-1),
				attestationDataLen: 1,
				methodLen:          1,
				urlLen:             1,
				selectorLen:        1,
				headersLen:         1,
				optionalFieldsLen:  1,
			},
			wantHeader: nil,
			wantErr:    true,
		},
		{
			name: "long buffer",
			args: args{
				header:             make([]byte, TARGET_ALIGNMENT*2+1),
				attestationDataLen: 1,
				methodLen:          1,
				urlLen:             1,
				selectorLen:        1,
				headersLen:         1,
				optionalFieldsLen:  1,
			},
			wantHeader: nil,
			wantErr:    true,
		},
		{
			name: "valid",
			args: args{
				header:             header,
				attestationDataLen: 10,
				methodLen:          5,
				urlLen:             40,
				selectorLen:        30,
				headersLen:         256,
				optionalFieldsLen:  64,
			},
			wantHeader: []byte{10, 0, 8, 0, 8, 0, 5, 0, 1, 0, 40, 0, 30, 0, 16, 0, 0, 1, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = CreateMetaHeader(tt.args.header, tt.args.attestationDataLen, tt.args.methodLen, tt.args.urlLen, tt.args.selectorLen, tt.args.headersLen, tt.args.optionalFieldsLen); (err != nil) != tt.wantErr {
				t.Errorf("CreateMetaHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(tt.wantHeader, tt.args.header) {
				t.Errorf("CreateMetaHeader() = %v, want %v", tt.args.header, tt.wantHeader)
			}
		})
	}
}

func TestDecodeMetaHeader(t *testing.T) {
	type args struct {
		header []byte
	}
	tests := []struct {
		name             string
		args             args
		wantParsedHeader *MetaHeader
		wantErr          bool
	}{
		{
			name: "nil buffer",
			args: args{
				header: nil,
			},
			wantParsedHeader: nil,
			wantErr:          true,
		},
		{
			name: "short buffer",
			args: args{
				header: make([]byte, TARGET_ALIGNMENT*2-1),
			},
			wantParsedHeader: nil,
			wantErr:          true,
		},
		{
			name: "long buffer",
			args: args{
				header: make([]byte, TARGET_ALIGNMENT*2+1),
			},
			wantParsedHeader: nil,
			wantErr:          true,
		},
		{
			name: "valid buffer",
			args: args{
				header: make([]byte, TARGET_ALIGNMENT*2),
			},
			wantParsedHeader: &MetaHeader{
				AttestationDataLen: 0,
				TimestampLen:       0,
				StatusCodeLen:      0,
				MethodLen:          0,
				ResponseFormatLen:  0,
				UrlLen:             0,
				SelectorLen:        0,
				EncodingOptionsLen: 0,
				HeadersLen:         0,
				OptionalFieldsLen:  0,
			},
			wantErr: false,
		},
		{
			name: "valid buffer 2",
			args: args{
				header: []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantParsedHeader: &MetaHeader{
				AttestationDataLen: 1,
				TimestampLen:       2,
				StatusCodeLen:      3,
				MethodLen:          4,
				ResponseFormatLen:  5,
				UrlLen:             6,
				SelectorLen:        7,
				EncodingOptionsLen: 8,
				HeadersLen:         9,
				OptionalFieldsLen:  10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParsedHeader, err := DecodeMetaHeader(tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMetaHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParsedHeader, tt.wantParsedHeader) {
				t.Errorf("DecodeMetaHeader() = %v, want %v", gotParsedHeader, tt.wantParsedHeader)
			}
		})
	}
}

func Test_prepareDataAsInteger(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid short",
			data:    "200",
			want:    []byte{200, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		}, {
			name:    "valid short 2",
			data:    "64250",
			want:    []byte{0xfa, 0xfa, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "max decimal",
			data:    "18446744073709551615",
			want:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			wantErr: false,
		},
		{
			name:    "invalid",
			data:    "xyz",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			data:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no hex",
			data:    "FFFF",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no hex 2",
			data:    "0xffff",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "too big decimal",
			data:    "18446744073709551616",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "zero",
			data:    "0",
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareDataAsInteger(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareDataAsInteger() error = %v, wantErr %v, value %s", err, tt.wantErr, tt.data)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareDataAsInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareDataAsFloat(t *testing.T) {
	type args struct {
		data      string
		precision uint
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid float",
			args: args{
				data:      "3.01",
				precision: 2,
			},
			want:    []byte{45, 1, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, just enough precision",
			args: args{
				data:      "3.14159",
				precision: 5,
			},
			want:    []byte{47, 203, 4, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, bigger precision",
			args: args{
				data:      "3.1415",
				precision: 5,
			},
			want:    []byte{38, 203, 4, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float without fractions",
			args: args{
				data:      "3",
				precision: 1,
			},
			want:    []byte{30, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, without fractions, zero precision",
			args: args{
				data:      "3",
				precision: 0,
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "float with dot, zero precision",
			args: args{
				data:      "3.",
				precision: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "float with dot, with precision",
			args: args{
				data:      "3.",
				precision: 2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid float, with redundant .0",
			args: args{
				data:      "3.0",
				precision: 0,
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, with redundant .00",
			args: args{
				data:      "3.00",
				precision: 0,
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, with redundant .00 and precision",
			args: args{
				data:      "3.00",
				precision: 1,
			},
			want:    []byte{30, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, not enough precision",
			args: args{
				data:      "3.1415",
				precision: 2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "big precision",
			args: args{
				data:      "3.14",
				precision: 10,
			},
			want:    []byte{0, 250, 149, 79, 7, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "too big precision",
			args: args{
				data:      "3.14",
				precision: 20,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - string",
			args: args{
				data:      "float",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - hex",
			args: args{
				data:      "0xabcd",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - decimalless notation - not allowed",
			args: args{
				data:      "1234p-9",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - scientific notation - not allowed",
			args: args{
				data:      "0.1234e+9",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - scientific notation 2 - not allowed",
			args: args{
				data:      "0.1234e-09",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float - hex scientific notation - not allowed",
			args: args{
				data:      "0x0.1234p+09",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bigger than precision float",
			args: args{
				data:      "1234.123456789123456789",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bigger than precision float 2",
			args: args{
				data:      "999999999.1234567891",
				precision: 6,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "info loss, redundant zeroes",
			args: args{
				data:      "123456789123456789123456789.00",
				precision: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				data:      "",
				precision: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareDataAsFloat(tt.args.data, tt.args.precision)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareDataAsFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareDataAsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_WriteWithPadding(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		want    []byte
	}{
		{
			name:    "short",
			data:    []byte{1, 1, 1, 1},
			want:    []byte{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "long",
			data:    []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want:    []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
			wantErr: false,
		}, {
			name:    "no pad",
			data:    []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want:    []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var b bytes.Buffer
		recorder := positionRecorder.NewPositionRecorder(&b, 16)

		t.Run(tt.name, func(t *testing.T) {
			if _, err := WriteWithPadding(recorder, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteWithPadding() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := b.Bytes()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareDataAsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EncodeAttestationData(t *testing.T) {
	type args struct {
		data    string
		options *EncodingOptions
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "empty string, valid string encoding",
			args: args{
				data: "",
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid string, valid encoding",
			args: args{
				data: "a",
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:    []byte{0x61, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid string, invalid encoding",
			args: args{
				data: "a",
				options: &EncodingOptions{
					Value: "hex",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid string, empty encoding",
			args: args{
				data: "a",
				options: &EncodingOptions{
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid short integer",
			args: args{
				data: "200",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    []byte{200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid short integer 2",
			args: args{
				data: "64250",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    []byte{0xfa, 0xfa, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid short integer 2, precision is ignored",
			args: args{
				data: "64250",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 10,
				},
			},
			want:    []byte{0xfa, 0xfa, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "too big decimal integer",
			args: args{
				data: "9999999999999999999999999999999999999999",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid integer",
			args: args{
				data: "xyz",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid integer 2",
			args: args{
				data: "abc",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty integer string",
			args: args{
				data: "",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "max decimal integer",
			args: args{
				data: "18446744073709551615",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "negative integer",
			args: args{
				data: "-42",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "zero",
			args: args{
				data: "0",
				options: &EncodingOptions{
					Value:     "int",
					Precision: 0,
				},
			},
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, just enough precision",
			args: args{
				data: "3.14159",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 5,
				},
			},
			want:    []byte{47, 203, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, bigger precision",
			args: args{
				data: "3.1415",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 5,
				},
			},
			want:    []byte{38, 203, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float without fractions",
			args: args{
				data: "3",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 1,
				},
			},
			want:    []byte{30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, without fractions, zero precision",
			args: args{
				data: "3",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "float with dot, zero precision",
			args: args{
				data: "3.",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "float with dot, with precision",
			args: args{
				data: "3.",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 3,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid float, with redundant zero",
			args: args{
				data: "3.0",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, with redundant zeroes, with precision",
			args: args{
				data: "3.00",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 2,
				},
			},
			want:    []byte{44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "valid float, not enough precision",
			args: args{
				data: "3.1415",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "big precision",
			args: args{
				data: "3.14",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 10,
				},
			},
			want:    []byte{0, 250, 149, 79, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "too big precision",
			args: args{
				data: "3.14",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 20,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid float",
			args: args{
				data: "float",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 6,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				data: "",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "big valid float",
			args: args{
				data: "123456789.123456",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 6,
				},
			},
			want:    []byte{128, 145, 15, 134, 72, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "too big float",
			args: args{
				data: "1234567890000000000000.123456",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 6,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative float",
			args: args{
				data: "-3.14",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 6,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "too precise float",
			args: args{
				data: "0.1234567899999999999999999999",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 12,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "too big integer, parsed as float",
			args: args{
				data: "99999999999999999999999999999",
				options: &EncodingOptions{
					Value:     "float",
					Precision: 12,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Error(t.Name(), "wasn't expected to panic")
				}
			}()

			got, err := EncodeAttestationData(tt.args.data, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeAttestationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeAttestationData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeAttestationData(t *testing.T) {
	type args struct {
		buf       []byte
		stringLen int
		options   *EncodingOptions
	}
	tests := []struct {
		name           string
		args           args
		want           string
		wantErr        bool
		checkRoundTrip bool
	}{
		{
			name: "nil",
			args: args{
				buf:     nil,
				options: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "nil options",
			args: args{
				buf:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				options: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid encoding options",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				options: &EncodingOptions{
					Value: "integer",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid string len",
			args: args{
				buf:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 20,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           "",
			wantErr:        true,
			checkRoundTrip: false,
		},
		{
			name: "zero-byte string",
			args: args{
				buf:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 10,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           string([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "zero-byte string, cut all",
			args: args{
				buf:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 0,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           "",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "zero-byte string, two blocks",
			args: args{
				buf:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 10,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           string([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			wantErr:        false,
			checkRoundTrip: false, // encoding this will create only 1 block for an empty string
		},
		{
			name: "normal string, no padding",
			args: args{
				buf:       []byte{0x61, 0x61, 0x61, 0x61, 0x62, 0x62, 0x62, 0x62, 0x63, 0x63, 0x63, 0x63, 0x7a, 0x7a, 0x7a, 0x7a},
				stringLen: 16,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           "aaaabbbbcccczzzz",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "normal string, no padding, cut string",
			args: args{
				buf:       []byte{0x61, 0x61, 0x61, 0x61, 0x62, 0x62, 0x62, 0x62, 0x63, 0x63, 0x63, 0x63, 0x7a, 0x7a, 0x7a, 0x7a},
				stringLen: 10,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           "aaaabbbbcc",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "normal string, with padding",
			args: args{
				buf:       []byte{0x61, 0x61, 0x61, 0x61, 0x62, 0x62, 0x62, 0x62, 0x63, 0x63, 0x63, 0x63, 0x7a, 0x7a, 0x7a, 0x7a, 0x30, 0x30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 18,
				options: &EncodingOptions{
					Value: "string",
				},
			},
			want:           "aaaabbbbcccczzzz00",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "zero integer",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				options: &EncodingOptions{
					Value: "int",
				},
			},
			want:           "0",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "normal integer",
			args: args{
				buf: []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0},
				options: &EncodingOptions{
					Value: "int",
				},
			},
			want:           "506097522914230528",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "normal integer, only 8 bytes are taken into account",
			args: args{
				buf: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				options: &EncodingOptions{
					Value: "int",
				},
			},
			want:           "506097522914230528",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "small float, big precision",
			args: args{
				buf:       []byte{0, 250, 149, 79, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 4,
				options: &EncodingOptions{
					Value:     "float",
					Precision: 10,
				},
			},
			want:           "3.14",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "small float, big precision, redundant zeroes",
			args: args{
				buf:       []byte{0, 250, 149, 79, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 6,
				options: &EncodingOptions{
					Value:     "float",
					Precision: 10,
				},
			},
			want:           "3.1400",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float, zero precision, redundant bytes",
			args: args{
				buf: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:           "506097522914230528",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float, zero precision",
			args: args{
				buf: []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0},
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:           "506097522914230528",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float 2, zero precision",
			args: args{
				buf: []byte{128, 145, 15, 134, 72, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				options: &EncodingOptions{
					Value:     "float",
					Precision: 0,
				},
			},
			want:           "123456789123456",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float 2, 1 precision",
			args: args{
				buf:       []byte{128, 145, 15, 134, 72, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 16,
				options: &EncodingOptions{
					Value:     "float",
					Precision: 1,
				},
			},
			want:           "12345678912345.6",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float 2, 6 precision",
			args: args{
				buf:       []byte{128, 145, 15, 134, 72, 112, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 16,
				options: &EncodingOptions{
					Value:     "float",
					Precision: 6,
				},
			},
			want:           "123456789.123456",
			wantErr:        false,
			checkRoundTrip: true,
		},
		{
			name: "float 3, 2 precision, redundant zeroes",
			args: args{
				buf:       []byte{44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				stringLen: 4,
				options: &EncodingOptions{
					Value:     "float",
					Precision: 2,
				},
			},
			want:           "3.00",
			wantErr:        false,
			checkRoundTrip: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeAttestationData(tt.args.buf, tt.args.stringLen, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeAttestationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeAttestationData() = %v, want %v", got, tt.want)
				return
			}
			if err != nil || !tt.checkRoundTrip {
				return
			}

			encoded, err := EncodeAttestationData(got, tt.args.options)
			if err != nil {
				t.Errorf("Encoding successfully decoded data shouldn't fail, got err = %v", err)
				return
			}

			decoded, err := DecodeAttestationData(encoded, tt.args.stringLen, tt.args.options)
			if err != nil {
				t.Errorf("Shouldn't fail to decode -> encode -> decode, got err = %v", err)
				return
			}

			if got != decoded {
				t.Errorf("Input should be equal to output after a roundtrip decode -> encode -> decode, expected \"%v\", got \"%v\"", got, decoded)
			}
		})
	}
}

func Test_EncodeResponseFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		want    []byte
		wantErr bool
	}{
		{
			name:    "json",
			format:  "json",
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "html",
			format:  "html",
			want:    []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "xml",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Error(t.Name(), "wasn't expected to panic")
				}
			}()

			got, err := EncodeResponseFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeResponseFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeResponseFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeResponseFormat(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				buf: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "short",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "long",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "json",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    RESPONSE_FORMAT_JSON,
			wantErr: false,
		},
		{
			name: "html",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    RESPONSE_FORMAT_HTML,
			wantErr: false,
		},
		{
			name: "invalid format",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid format 2",
			args: args{
				buf: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeResponseFormat(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeResponseFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeResponseFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EncodeEncodingOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *EncodingOptions
		want    []byte
		wantErr bool
	}{
		{
			name: "string",
			options: &EncodingOptions{
				Value: "string",
			},
			want: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "string, precision is ignored",
			options: &EncodingOptions{
				Value:     "string",
				Precision: 5,
			},
			want: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "int",
			options: &EncodingOptions{
				Value: "int",
			},
			want: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "int, precision is ignored",
			options: &EncodingOptions{
				Value:     "int",
				Precision: 5,
			},
			want: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "float",
			options: &EncodingOptions{
				Value: "float",
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "float, with precision 5",
			options: &EncodingOptions{
				Value:     "float",
				Precision: 5,
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "float, with precision 10",
			options: &EncodingOptions{
				Value:     "float",
				Precision: 10,
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0x0a, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "invalid value",
			options: &EncodingOptions{
				Value:     "bytes",
				Precision: 10,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid precision",
			options: &EncodingOptions{
				Value:     "float",
				Precision: 100,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Error(t.Name(), "wasn't expected to panic")
				}
			}()

			got, err := EncodeEncodingOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeEncodingOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeEncodingOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeEncodingOptions(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *EncodingOptions
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				buf: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "short buffer",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "long buffer",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid format",
			args: args{
				buf: []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: &EncodingOptions{
				Value: "string",
			},
			wantErr: false,
		},
		{
			name: "int",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: &EncodingOptions{
				Value: "int",
			},
			wantErr: false,
		},
		{
			name: "float",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: &EncodingOptions{
				Value:     "float",
				Precision: 0,
			},
			wantErr: false,
		},
		{
			name: "float, with precision 5",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0},
			},
			want: &EncodingOptions{
				Value:     "float",
				Precision: 5,
			},
			wantErr: false,
		},
		{
			name: "float, with precision 10",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 0x0a, 0, 0, 0, 0, 0, 0, 0},
			},
			want: &EncodingOptions{
				Value:     "float",
				Precision: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeEncodingOptions(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeEncodingOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeEncodingOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EncodeHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    []byte
	}{
		{
			name:    "no headers",
			headers: make(map[string]string),
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "one short header",
			headers: map[string]string{
				"a": "b",
			},
			want: []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0x61, 0x3a, 0x62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "two short headers",
			headers: map[string]string{
				"a": "b",
				"c": "?",
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0x61, 0x3a, 0x62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0x63, 0x3a, 0x3f, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "one short, one long headers",
			headers: map[string]string{
				"a": "b",
				"c": "abcdefghijklmnopqrstuvwxyz",
			},
			want: []byte{2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0x61, 0x3a, 0x62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 28, 0, 0x63, 0x3a, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeHeaders(tt.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeHeaders(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name           string
		args           args
		want           map[string]string
		wantErr        bool
		checkRoundTrip bool
	}{
		{
			name: "nil buffer",
			args: args{
				buf: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "buffer too short",
			args: args{
				buf: []byte{0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty map",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "encodes 1 block, doesn't have blocks",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encodes 1 header, doesn't have blocks",
			args: args{
				buf: []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded length and content mismatch",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded length and content mismatch",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded header entry mismatch - encodes entry of 1 char, no entries",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded header entry - encodes entry with no separator",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded header entry - encodes entry with no key",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 14, 0, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded header entry is valid, mismatch with number of headers meta info",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 14, 0, 97, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "encoded header entry is valid with invalid padding",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 15, 0, 97, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid encoding with 1 header, no padding",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 14, 0, 97, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97},
			},
			want: map[string]string{
				"a": "aaaaaaaaaaaa",
			},
			wantErr: false,
		},
		{
			name: "valid encoding with 1 header, padding",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 15, 0, 97, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: map[string]string{
				"a": "aaaaaaaaaaaaa",
			},
			wantErr: false,
		},
		{
			name: "valid encoding with 2 headers",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 5, 0, 98, 98, 0x3a, 99, 99, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 97, 0x3a, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: map[string]string{
				"a":  "aaaaaaaaaaaaa",
				"bb": "cc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeHeaders(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeHeaders() error = %v, wantErr %v, got %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeHeaders() = %v, want %v", got, tt.want)
			}
			if err != nil || !tt.checkRoundTrip {
				return
			}

			encoded := EncodeHeaders(got)
			decoded, err := DecodeHeaders(encoded)
			if err != nil {
				t.Errorf("Shouldn't fail to decode -> encode -> decode, got err = %v", err)
				return
			}

			if reflect.DeepEqual(got, decoded) {
				t.Errorf("Input should be equal to output after a roundtrip decode -> encode -> decode, expected \"%v\", got \"%v\"", got, decoded)
			}
		})
	}
}

func Test_EncodeOptionalFields(t *testing.T) {
	htmlResultElement := "element"
	htmlResultValue := "value"
	htmlResultInvalid := "string"

	exampleContentType := "text/plain"
	exampleLongerContentType := "multipart/form-data"

	exampleBody := "short body text"
	exampleBodyLong := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam laoreet mattis diam eget finibus"

	type args struct {
		htmlResultType     *string
		requestContentType *string
		requestBody        *string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			want:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "invalid html result",
			args: args{
				htmlResultType: &htmlResultInvalid,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "html element",
			args: args{
				htmlResultType: &htmlResultElement,
			},
			want:    []byte{1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "html value",
			args: args{
				htmlResultType: &htmlResultValue,
			},
			want:    []byte{1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "content type",
			args: args{
				requestContentType: &exampleContentType,
			},
			want:    []byte{2, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "longer content type",
			args: args{
				requestContentType: &exampleLongerContentType,
			},
			want:    []byte{2, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 19, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x2f, 0x66, 0x6f, 0x72, 0x6d, 0x2d, 0x64, 0x61, 0x74, 0x61, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "body",
			args: args{
				requestBody: &exampleBody,
			},
			want:    []byte{4, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			wantErr: false,
		},
		{
			name: "body long",
			args: args{
				requestBody: &exampleBodyLong,
			},
			want:    []byte{4, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 76, 111, 114, 101, 109, 32, 105, 112, 115, 117, 109, 32, 100, 111, 108, 111, 114, 32, 115, 105, 116, 32, 97, 109, 101, 116, 44, 32, 99, 111, 110, 115, 101, 99, 116, 101, 116, 117, 114, 32, 97, 100, 105, 112, 105, 115, 99, 105, 110, 103, 32, 101, 108, 105, 116, 46, 32, 78, 117, 108, 108, 97, 109, 32, 108, 97, 111, 114, 101, 101, 116, 32, 109, 97, 116, 116, 105, 115, 32, 100, 105, 97, 109, 32, 101, 103, 101, 116, 32, 102, 105, 110, 105, 98, 117, 115},
			wantErr: false,
		},
		{
			name: "html element, content type",
			args: args{
				htmlResultType:     &htmlResultElement,
				requestContentType: &exampleContentType,
			},
			want:    []byte{3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "html element, body",
			args: args{
				htmlResultType: &htmlResultElement,
				requestBody:    &exampleBody,
			},
			want:    []byte{5, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			wantErr: false,
		},
		{
			name: "content type, body",
			args: args{
				requestContentType: &exampleContentType,
				requestBody:        &exampleBody,
			},
			want:    []byte{6, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			wantErr: false,
		},
		{
			name: "html result, content type, body",
			args: args{
				htmlResultType:     &htmlResultValue,
				requestContentType: &exampleContentType,
				requestBody:        &exampleBody,
			},
			want:    []byte{7, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Error(t.Name(), "wasn't expected to panic")
				}
			}()

			got, err := EncodeOptionalFields(tt.args.htmlResultType, tt.args.requestContentType, tt.args.requestBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeOptionalFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeOptionalFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareStringPtr(ptr1, ptr2 *string) bool {
	if ptr1 == nil || ptr2 == nil {
		return ptr1 == ptr2
	}

	return *ptr1 == *ptr2
}

func TestDecodeOptionalFields(t *testing.T) {
	htmlResultElement := "element"
	htmlResultValue := "value"

	exampleContentType := "text/plain"
	exampleLongerContentType := "multipart/form-data"

	exampleBody := "short body text"
	exampleBodyLong := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam laoreet mattis diam eget finibus"

	type args struct {
		buf []byte
	}
	tests := []struct {
		name                   string
		args                   args
		wantHtmlResultType     *string
		wantRequestContentType *string
		wantRequestBody        *string
		wantErr                bool
		checkRoundTrip         bool
	}{
		{
			name: "nil buffer",
			args: args{
				buf: nil,
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "short buffer",
			args: args{
				buf: make([]byte, 4*TARGET_ALIGNMENT-1),
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "zero buffer",
			args: args{
				buf: make([]byte, 4*TARGET_ALIGNMENT),
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "no optional fields",
			args: args{
				buf: append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0}, make([]byte, 3*TARGET_ALIGNMENT)...),
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "no optional fields, incorrect block count in header",
			args: args{
				buf: append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0}, make([]byte, 3*TARGET_ALIGNMENT)...),
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "unknown html result type",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "html element",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     &htmlResultElement,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "html value",
			args: args{
				buf: []byte{1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     &htmlResultValue,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "invalid content type meta header - too long",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "content type",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x74, 0x65, 0x78, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: &exampleContentType,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "longer content type",
			args: args{
				buf: []byte{2, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 19, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x2f, 0x66, 0x6f, 0x72, 0x6d, 0x2d, 0x64, 0x61, 0x74, 0x61, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: &exampleLongerContentType,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "invalid body meta header - too long",
			args: args{
				buf: []byte{4, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        nil,
			wantErr:                true,
			checkRoundTrip:         false,
		},
		{
			name: "body",
			args: args{
				buf: []byte{4, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        &exampleBody,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "body long",
			args: args{
				buf: []byte{4, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 76, 111, 114, 101, 109, 32, 105, 112, 115, 117, 109, 32, 100, 111, 108, 111, 114, 32, 115, 105, 116, 32, 97, 109, 101, 116, 44, 32, 99, 111, 110, 115, 101, 99, 116, 101, 116, 117, 114, 32, 97, 100, 105, 112, 105, 115, 99, 105, 110, 103, 32, 101, 108, 105, 116, 46, 32, 78, 117, 108, 108, 97, 109, 32, 108, 97, 111, 114, 101, 101, 116, 32, 109, 97, 116, 116, 105, 115, 32, 100, 105, 97, 109, 32, 101, 103, 101, 116, 32, 102, 105, 110, 105, 98, 117, 115},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: nil,
			wantRequestBody:        &exampleBodyLong,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "html element, content type",
			args: args{
				buf: []byte{3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantHtmlResultType:     &htmlResultElement,
			wantRequestContentType: &exampleContentType,
			wantRequestBody:        nil,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "html element, body",
			args: args{
				buf: []byte{5, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			},
			wantHtmlResultType:     &htmlResultElement,
			wantRequestContentType: nil,
			wantRequestBody:        &exampleBody,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "content type, body",
			args: args{
				buf: []byte{6, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			},
			wantHtmlResultType:     nil,
			wantRequestContentType: &exampleContentType,
			wantRequestBody:        &exampleBody,
			wantErr:                false,
			checkRoundTrip:         true,
		},
		{
			name: "html result, content type, body",
			args: args{
				buf: []byte{7, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 115, 104, 111, 114, 116, 32, 98, 111, 100, 121, 32, 116, 101, 120, 116, 0},
			},
			wantHtmlResultType:     &htmlResultValue,
			wantRequestContentType: &exampleContentType,
			wantRequestBody:        &exampleBody,
			wantErr:                false,
			checkRoundTrip:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHtmlResultType, gotRequestContentType, gotRequestBody, err := DecodeOptionalFields(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeOptionalFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareStringPtr(gotHtmlResultType, tt.wantHtmlResultType) {
				t.Errorf("DecodeOptionalFields() gotHtmlResultType = %v, want %v", gotHtmlResultType, tt.wantHtmlResultType)
			}
			if !compareStringPtr(gotRequestContentType, tt.wantRequestContentType) {
				t.Errorf("DecodeOptionalFields() gotRequestContentType = %v, want %v", gotRequestContentType, tt.wantRequestContentType)
			}
			if !compareStringPtr(gotRequestBody, tt.wantRequestBody) {
				t.Errorf("DecodeOptionalFields() gotRequestBody = %v, want %v", gotRequestBody, tt.wantRequestBody)
			}
			if err != nil || !tt.checkRoundTrip {
				return
			}

			encoded, err := EncodeOptionalFields(gotHtmlResultType, gotRequestContentType, gotRequestBody)
			if err != nil {
				t.Errorf("DecodeOptionalFields() encode error = %v", err)
				return
			}

			decodedHtmlResult, decodedContentType, decodedBody, err := DecodeOptionalFields(encoded)
			if err != nil {
				t.Errorf("Shouldn't fail to decode -> encode -> decode, got err = %v", err)
				return
			}

			htmlResultEqual := compareStringPtr(gotHtmlResultType, decodedHtmlResult)
			contentTypeEqual := compareStringPtr(gotRequestContentType, decodedContentType)
			bodyEqual := compareStringPtr(gotRequestBody, decodedBody)

			var errorStringBuilder strings.Builder

			if !htmlResultEqual {
				errorStringBuilder.WriteString(fmt.Sprintf("HTML result type should be equal to output after decode -> encode -> decode,\n\texpected = \"%v\", got = \"%v\"\n", *gotHtmlResultType, *decodedHtmlResult))
			}

			if !contentTypeEqual {
				errorStringBuilder.WriteString(fmt.Sprintf("Request content type should be equal to output after decode -> encode -> decode,\n\texpected = \"%v\", got = \"%v\"\n", *gotRequestContentType, *decodedContentType))
			}

			if !bodyEqual {
				errorStringBuilder.WriteString(fmt.Sprintf("Request body should be equal to output after decode -> encode -> decode,\n\texpected:\n\t\t\"%v\"\ngot:\n\t\t\"%v\"\n", *gotRequestContentType, *decodedContentType))
			}

			errorString := errorStringBuilder.String()

			if errorString != "" {
				t.Error(errorString)
			}
		})
	}
}
