package p2pstate

import (
	"errors"
	"testing"
)

func TestPurchaseRequisition_StateMachine(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		targetStatus  string
		wantErr       error
	}{
		{
			name:          "Valid: Draft to Pending Approval",
			currentStatus: PRStatusDraft,
			targetStatus:  PRStatusPendingApproval,
			wantErr:       nil,
		},
		{
			name:          "Valid: Pending Approval to Approved",
			currentStatus: PRStatusPendingApproval,
			targetStatus:  PRStatusApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Pending Approval to Rejected",
			currentStatus: PRStatusPendingApproval,
			targetStatus:  PRStatusRejected,
			wantErr:       nil,
		},
		{
			name:          "Valid: Approved to PO Created",
			currentStatus: PRStatusApproved,
			targetStatus:  PRStatusPOCreated,
			wantErr:       nil,
		},
		{
			name:          "Valid: Rejected to Draft for Revisions",
			currentStatus: PRStatusRejected,
			targetStatus:  PRStatusDraft,
			wantErr:       nil,
		},
		{
			name:          "Invalid: Draft directly converted to PO",
			currentStatus: PRStatusDraft,
			targetStatus:  PRStatusPOCreated,
			wantErr:       ErrInvalidPRTransition,
		},
		{
			name:          "Invalid: Rejected PR converted to PO directly",
			currentStatus: PRStatusRejected,
			targetStatus:  PRStatusPOCreated,
			wantErr:       ErrInvalidPRTransition,
		},
		{
			name:          "Invalid: POCreated PR moved back to Draft",
			currentStatus: PRStatusPOCreated,
			targetStatus:  PRStatusDraft,
			wantErr:       ErrInvalidPRTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePRTransition(tt.currentStatus, tt.targetStatus)
			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error wrapping %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPurchaseOrder_StateMachine(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		targetStatus  string
		wantErr       error
	}{
		{
			name:          "Valid: Draft to Submitted",
			currentStatus: POStatusDraft,
			targetStatus:  POStatusSubmitted,
			wantErr:       nil,
		},
		{
			name:          "Valid: Submitted to Manager Approved",
			currentStatus: POStatusSubmitted,
			targetStatus:  POStatusManagerApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Manager Approved to Owner Approved",
			currentStatus: POStatusManagerApproved,
			targetStatus:  POStatusOwnerApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Sent to Received",
			currentStatus: POStatusSent,
			targetStatus:  POStatusReceived,
			wantErr:       nil,
		},
		{
			name:          "Valid: Received to Invoiced",
			currentStatus: POStatusReceived,
			targetStatus:  POStatusInvoiced,
			wantErr:       nil,
		},
		{
			name:          "Valid: Invoiced to Paid",
			currentStatus: POStatusInvoiced,
			targetStatus:  POStatusPaid,
			wantErr:       nil,
		},
		{
			name:          "Invalid: Draft directly marked as Received",
			currentStatus: POStatusDraft,
			targetStatus:  POStatusReceived,
			wantErr:       ErrInvalidPOTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePOTransition(tt.currentStatus, tt.targetStatus)
			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error wrapping %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestThreeWayMatch(t *testing.T) {
	tests := []struct {
		name             string
		poQty            float64
		poUnitPrice      float64
		grnReceivedQty   float64
		invBilledQty     float64
		invUnitPrice     float64
		tolerancePercent float64
		wantStatus       ThreeWayMatchStatus
		wantErr          error
	}{
		{
			name:             "Perfect 3-Way Match",
			poQty:            100,
			poUnitPrice:      50000,
			grnReceivedQty:   100,
			invBilledQty:     100,
			invUnitPrice:     50000,
			tolerancePercent: 0.0,
			wantStatus:       MatchSuccess,
			wantErr:          nil,
		},
		{
			name:             "Price Variance Within Allowed Tolerance",
			poQty:            50,
			poUnitPrice:      100000,
			grnReceivedQty:   50,
			invBilledQty:     50,
			invUnitPrice:     101500, // +1.5%
			tolerancePercent: 2.0,
			wantStatus:       MatchSuccess,
			wantErr:          nil,
		},
		{
			name:             "Price Variance Exceeding Allowed Tolerance",
			poQty:            50,
			poUnitPrice:      100000,
			grnReceivedQty:   50,
			invBilledQty:     50,
			invUnitPrice:     105000, // +5%
			tolerancePercent: 2.0,
			wantStatus:       MatchPriceError,
			wantErr:          ErrPriceVarianceExceeded,
		},
		{
			name:             "Over-Billed Quantity: Billed More Than Received (GRN)",
			poQty:            100,
			poUnitPrice:      50000,
			grnReceivedQty:   80,
			invBilledQty:     100,
			invUnitPrice:     50000,
			tolerancePercent: 0.0,
			wantStatus:       MatchQuantityError,
			wantErr:          ErrOverBilledQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ValidateThreeWayMatch(
				tt.poQty, tt.poUnitPrice,
				tt.grnReceivedQty,
				tt.invBilledQty, tt.invUnitPrice,
				tt.tolerancePercent,
			)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error wrapping %v, got %v", tt.wantErr, err)
				}
				if res.Status != tt.wantStatus {
					t.Errorf("res.Status = %v; want %v", res.Status, tt.wantStatus)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Status != tt.wantStatus {
				t.Errorf("res.Status = %v; want %v", res.Status, tt.wantStatus)
			}
		})
	}
}
