package kecpcrypto

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func hexDecodeString(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestHashSha256(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "One Block Message Sample",
			args: args{content: []byte("abc")},
			want: hexDecodeString("BA7816BF8F01CFEA414140DE5DAE2223B00361A396177A9CB410FF61F20015AD"),
		},
		{
			name: "Two Block Message Sample",
			args: args{content: []byte("abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq")},
			want: hexDecodeString("248D6A61D20638B8E5C026930C3E6039A33CE45964FF2167F6ECEDD419DB06C1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSha256(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HashSha256() = %v, want %v", got, tt.want)
			}
		})
	}
}
