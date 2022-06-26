package kecpvalidate

import "testing"

func TestIsValidUserName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"\\u202e Test",
			args{
				s: "\u202eaaa",
			},
			false,
		},
		{
			"\\u202f Test 1",
			args{
				s: "\u202faaa",
			},
			false,
		},
		{
			"\\u202f Test 2",
			args{
				s: "aaa\u202f",
			},
			false,
		},
		{
			"Normal Unicode Test",
			args{
				s: "中文",
			},
			true,
		},
		{
			"Length Test",
			args{
				s: "0123456789abcdefg",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAValidUserName(tt.args.s); got != tt.want {
				t.Errorf("IsValidUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
