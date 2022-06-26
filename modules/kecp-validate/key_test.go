package kecpvalidate

import "testing"

func TestIsAValidCryptoKey(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Shannon Entropy Test 1",
			args{s: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			false,
		},
		{
			"Shannon Entropy Test 2",
			args{s: "aaagaaotanaaqaabaaaaaapaaafaaamaaahaaaraasacaadaaaajaeaalaaaaaaa"},
			false,
		},
		{
			"Shannon Entropy Test 3",
			args{s: "K7EIM7IXDtxljByBCPQmo-f8v3vOHgAkNLQCiEInpIHtqHa0CXflYZopjG13cYev"},
			true,
		},
		{
			"Length Test",
			args{s: "K7EIM7IXDtxljByBCPQmo-f8v3vOHgAkNLQCiEInpIHtqHa0CXflYZopjG13cYe"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAValidCryptoKey(tt.args.s); got != tt.want {
				t.Errorf("IsAValidCryptoKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
