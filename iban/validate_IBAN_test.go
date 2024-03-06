package iban_test

import (
	"fmt"
	"testing"
	"web-iban/iban"
)

func TestValidateIBAN(t *testing.T) {
	type args struct {
		iban string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid IBAN for Germany",
			args:    args{"DE89370400440532013000"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid IBAN for United Kingdom",
			args:    args{"GB29NWBK60161331926819"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid IBAN (checksum failure)",
			args:    args{"GB29NWBK60161331926818"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid IBAN (length mismatch)",
			args:    args{"GB29NWBK60161331926819X"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid IBAN (invalid character)",
			args:    args{"GB29NWBK6016131.926X19"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid IBAN (empty string)",
			args:    args{""},
			want:    false,
			wantErr: true,
		},
	}
	err := iban.InitIbanData("data/")
	if err != nil {
		fmt.Println(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := iban.ValidateIBAN(tt.args.iban)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIBAN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateIBAN() = %v, want %v", got, tt.want)
			}
		})
	}
}
