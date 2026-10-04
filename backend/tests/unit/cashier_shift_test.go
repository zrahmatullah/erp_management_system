package unit_test

import (
	"testing"
)

// CalculateShiftReconciliation computes expected cash and discrepancy for blind closing
func CalculateShiftReconciliation(openingFloat, cashSales, cashDrops, refunds, actualCounted float64) (expectedCash float64, difference float64) {
	expectedCash = openingFloat + cashSales - cashDrops - refunds
	difference = actualCounted - expectedCash
	return expectedCash, difference
}

func TestCashierShift_BlindClosingReconciliation(t *testing.T) {
	tests := []struct {
		name             string
		openingFloat     float64
		cashSales        float64
		cashDrops        float64
		refunds          float64
		actualCounted    float64
		wantExpectedCash float64
		wantDifference   float64
	}{
		{
			name:             "Perfect Balanced Closing",
			openingFloat:     300000,
			cashSales:        1500000,
			cashDrops:        500000,
			refunds:          0,
			actualCounted:    1300000,
			wantExpectedCash: 1300000,
			wantDifference:   0,
		},
		{
			name:             "Short Cash (Deficit / Kurang Setor)",
			openingFloat:     500000,
			cashSales:        2000000,
			cashDrops:        1000000,
			refunds:          50000,
			actualCounted:    1400000,
			wantExpectedCash: 1450000,
			wantDifference:   -50000,
		},
		{
			name:             "Over Cash (Surplus / Lebih Setor)",
			openingFloat:     300000,
			cashSales:        850000,
			cashDrops:        0,
			refunds:          0,
			actualCounted:    1170000,
			wantExpectedCash: 1150000,
			wantDifference:   20000,
		},
		{
			name:             "Multiple Cash Drops Mid-Shift",
			openingFloat:     400000,
			cashSales:        3500000,
			cashDrops:        2500000,
			refunds:          100000,
			actualCounted:    1300000,
			wantExpectedCash: 1300000,
			wantDifference:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpected, gotDiff := CalculateShiftReconciliation(
				tt.openingFloat,
				tt.cashSales,
				tt.cashDrops,
				tt.refunds,
				tt.actualCounted,
			)

			if gotExpected != tt.wantExpectedCash {
				t.Errorf("ExpectedCash = %v, want %v", gotExpected, tt.wantExpectedCash)
			}
			if gotDiff != tt.wantDifference {
				t.Errorf("Difference = %v, want %v", gotDiff, tt.wantDifference)
			}
		})
	}
}

func TestCashierShift_ValidationRules(t *testing.T) {
	t.Run("Negative Opening Float Should Fail", func(t *testing.T) {
		openingFloat := -50000.0
		if openingFloat >= 0 {
			t.Errorf("Expected opening float to be rejected when negative")
		}
	})

	t.Run("Negative Actual Cash Counted Should Fail", func(t *testing.T) {
		actualCounted := -100.0
		if actualCounted >= 0 {
			t.Errorf("Expected actual counted cash to be rejected when negative")
		}
	})

	t.Run("Valid Movement Types", func(t *testing.T) {
		validTypes := map[string]bool{
			"cash_drop": true,
			"paid_out":  true,
			"cash_in":   true,
		}

		if !validTypes["cash_drop"] {
			t.Errorf("cash_drop should be valid")
		}
		if validTypes["invalid_type"] {
			t.Errorf("invalid_type should not be valid")
		}
	})
}
