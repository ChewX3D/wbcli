package ordercmd

import "testing"

func TestValidatePositiveFloatFlag(t *testing.T) {
	tests := []struct {
		name      string
		flag      string
		value     float64
		wantError string
	}{
		{name: "valid", flag: "--amount", value: 1},
		{name: "zero", flag: "--amount", value: 0, wantError: "--amount must be greater than 0"},
		{name: "negative", flag: "--price", value: -1, wantError: "--price must be greater than 0"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validatePositiveFloatFlag(testCase.flag, testCase.value)
			if testCase.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", testCase.wantError)
			}
			if err.Error() != testCase.wantError {
				t.Fatalf("unexpected error. got %q want %q", err.Error(), testCase.wantError)
			}
		})
	}
}

func TestValidateNonNegativeInt64Flag(t *testing.T) {
	tests := []struct {
		name      string
		flag      string
		value     int64
		wantError string
	}{
		{name: "valid", flag: "--expiration", value: 0},
		{name: "positive", flag: "--expiration", value: 10},
		{name: "negative", flag: "--expiration", value: -1, wantError: "--expiration must be greater than or equal to 0"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateNonNegativeInt64Flag(testCase.flag, testCase.value)
			if testCase.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", testCase.wantError)
			}
			if err.Error() != testCase.wantError {
				t.Fatalf("unexpected error. got %q want %q", err.Error(), testCase.wantError)
			}
		})
	}
}

func TestValidateBase(t *testing.T) {
	tests := []struct {
		name      string
		options   baseOptions
		wantError string
	}{
		{
			name: "valid",
			options: baseOptions{
				Market: "BTC_PERP",
				Side:   "buy",
			},
		},
		{
			name: "valid case insensitive side",
			options: baseOptions{
				Market: "BTC_PERP",
				Side:   "SELL",
			},
		},
		{
			name: "missing market",
			options: baseOptions{
				Side: "buy",
			},
			wantError: "--market is required",
		},
		{
			name: "missing side",
			options: baseOptions{
				Market: "BTC_PERP",
			},
			wantError: "--side is required",
		},
		{
			name: "unsupported side",
			options: baseOptions{
				Market: "BTC_PERP",
				Side:   "hold",
			},
			wantError: "--side must be one of: buy, sell",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateBase(testCase.options)
			if testCase.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", testCase.wantError)
			}
			if err.Error() != testCase.wantError {
				t.Fatalf("unexpected error. got %q want %q", err.Error(), testCase.wantError)
			}
		})
	}
}

func TestValidateAmountMode(t *testing.T) {
	tests := []struct {
		name      string
		mode      string
		wantError string
	}{
		{name: "constant", mode: "constant"},
		{name: "arithmetic", mode: "arithmetic"},
		{name: "geometric", mode: "geometric"},
		{name: "capped geometric", mode: "capped-geometric"},
		{name: "fibonacci", mode: "fibonacci"},
		{name: "custom list", mode: "custom-list"},
		{name: "unsupported", mode: "random", wantError: amountModeErrorMessage},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateAmountMode(testCase.mode)
			if testCase.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", testCase.wantError)
			}
			if err.Error() != testCase.wantError {
				t.Fatalf("unexpected error. got %q want %q", err.Error(), testCase.wantError)
			}
		})
	}
}
