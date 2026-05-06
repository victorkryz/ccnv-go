package main

import "testing"

type currencyDescriptor struct {
	minorUnits int64
	majorUnits float64
	precision  int
	code       string
}

func TestCurrencyInitialization(t *testing.T) {
	tests := []currencyDescriptor{
		{minorUnits: 2500, majorUnits: 25, precision: 2, code: "usd"},
		{minorUnits: 4500, majorUnits: 45, precision: 2, code: "uah"},
		{minorUnits: 37, majorUnits: 0.37, precision: 2, code: "eur"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			sut := NewCurrencyFromMinorUnits(tt.minorUnits, tt.code, tt.precision)

			if got := sut.RawMinorUnits(); got != tt.minorUnits {
				t.Fatalf("RawMinorUnits() = %d, want %d", got, tt.minorUnits)
			}
			if got := sut.ToDouble(); got != tt.majorUnits {
				t.Fatalf("ToDouble() = %v, want %v", got, tt.majorUnits)
			}
			if got := sut.CurrencyCode(); got != tt.code {
				t.Fatalf("CurrencyCode() = %q, want %q", got, tt.code)
			}
			if got := sut.Precision(); got != tt.precision {
				t.Fatalf("Precision() = %d, want %d", got, tt.precision)
			}
		})
	}
}

func TestConvertCurrency(t *testing.T) {
	tests := []struct {
		name        string
		source      currencyDescriptor
		rate        float64
		codeTo      string
		precisionTo int
		want        currencyDescriptor
	}{
		{
			name:        "usd to uah",
			source:      currencyDescriptor{minorUnits: 2500, majorUnits: 25, precision: 2, code: "usd"},
			rate:        41.71458864,
			codeTo:      "uah",
			precisionTo: 2,
			want:        currencyDescriptor{minorUnits: 104286, majorUnits: 1042.86, precision: 2, code: "uah"},
		},
		{
			name:        "usd to eur",
			source:      currencyDescriptor{minorUnits: 1000, majorUnits: 10, precision: 2, code: "usd"},
			rate:        0.86389992,
			codeTo:      "eur",
			precisionTo: 2,
			want:        currencyDescriptor{minorUnits: 864, majorUnits: 8.64, precision: 2, code: "eur"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := NewCurrencyFromMinorUnits(tt.source.minorUnits, tt.source.code, tt.source.precision)

			target := ConvertCurrency(from, tt.codeTo, tt.rate, tt.precisionTo)

			if got := target.RawMinorUnits(); got != tt.want.minorUnits {
				t.Fatalf("RawMinorUnits() = %d, want %d", got, tt.want.minorUnits)
			}
			if got := target.ToDouble(); got != tt.want.majorUnits {
				t.Fatalf("ToDouble() = %v, want %v", got, tt.want.majorUnits)
			}
			if got := target.CurrencyCode(); got != tt.want.code {
				t.Fatalf("CurrencyCode() = %q, want %q", got, tt.want.code)
			}
			if got := target.Precision(); got != tt.want.precision {
				t.Fatalf("Precision() = %d, want %d", got, tt.want.precision)
			}
		})
	}
}
