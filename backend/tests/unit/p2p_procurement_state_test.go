package unit_test

import (
	"errors"
	"testing"

	"cafe-erp-system/backend/pkg/p2pstate"
)

// TestP2P_PurchaseRequisition_StateMachine tests PR status transitions and illegal state jump prevention.
func TestP2P_PurchaseRequisition_StateMachine(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		targetStatus  string
		wantErr       error
	}{
		{
			name:          "Valid: Draft to Pending Approval",
			currentStatus: p2pstate.PRStatusDraft,
			targetStatus:  p2pstate.PRStatusPendingApproval,
			wantErr:       nil,
		},
		{
			name:          "Valid: Pending Approval to Approved",
			currentStatus: p2pstate.PRStatusPendingApproval,
			targetStatus:  p2pstate.PRStatusApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Pending Approval to Rejected",
			currentStatus: p2pstate.PRStatusPendingApproval,
			targetStatus:  p2pstate.PRStatusRejected,
			wantErr:       nil,
		},
		{
			name:          "Valid: Approved to PO Created",
			currentStatus: p2pstate.PRStatusApproved,
			targetStatus:  p2pstate.PRStatusPOCreated,
			wantErr:       nil,
		},
		{
			name:          "Valid: Rejected to Draft for Revisions",
			currentStatus: p2pstate.PRStatusRejected,
			targetStatus:  p2pstate.PRStatusDraft,
			wantErr:       nil,
		},
		{
			name:          "Invalid: Draft directly converted to PO (Skips Approval)",
			currentStatus: p2pstate.PRStatusDraft,
			targetStatus:  p2pstate.PRStatusPOCreated,
			wantErr:       p2pstate.ErrInvalidPRTransition,
		},
		{
			name:          "Invalid: Rejected PR converted to PO directly",
			currentStatus: p2pstate.PRStatusRejected,
			targetStatus:  p2pstate.PRStatusPOCreated,
			wantErr:       p2pstate.ErrInvalidPRTransition,
		},
		{
			name:          "Invalid: POCreated PR moved back to Draft (Immutable)",
			currentStatus: p2pstate.PRStatusPOCreated,
			targetStatus:  p2pstate.PRStatusDraft,
			wantErr:       p2pstate.ErrInvalidPRTransition,
		},
		{
			name:          "Idempotent: Transition to same status",
			currentStatus: p2pstate.PRStatusApproved,
			targetStatus:  p2pstate.PRStatusApproved,
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p2pstate.ValidatePRTransition(tt.currentStatus, tt.targetStatus)
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

// TestP2P_PurchaseOrder_StateMachine tests PO status flow and approval gates.
func TestP2P_PurchaseOrder_StateMachine(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		targetStatus  string
		wantErr       error
	}{
		{
			name:          "Valid: Draft to Submitted",
			currentStatus: p2pstate.POStatusDraft,
			targetStatus:  p2pstate.POStatusSubmitted,
			wantErr:       nil,
		},
		{
			name:          "Valid: Submitted to Manager Approved",
			currentStatus: p2pstate.POStatusSubmitted,
			targetStatus:  p2pstate.POStatusManagerApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Manager Approved to Owner Approved",
			currentStatus: p2pstate.POStatusManagerApproved,
			targetStatus:  p2pstate.POStatusOwnerApproved,
			wantErr:       nil,
		},
		{
			name:          "Valid: Sent to Received",
			currentStatus: p2pstate.POStatusSent,
			targetStatus:  p2pstate.POStatusReceived,
			wantErr:       nil,
		},
		{
			name:          "Valid: Sent to Partially Received",
			currentStatus: p2pstate.POStatusSent,
			targetStatus:  p2pstate.POStatusPartiallyReceived,
			wantErr:       nil,
		},
		{
			name:          "Valid: Received to Invoiced",
			currentStatus: p2pstate.POStatusReceived,
			targetStatus:  p2pstate.POStatusInvoiced,
			wantErr:       nil,
		},
		{
			name:          "Valid: Invoiced to Paid",
			currentStatus: p2pstate.POStatusInvoiced,
			targetStatus:  p2pstate.POStatusPaid,
			wantErr:       nil,
		},
		{
			name:          "Invalid: Draft directly marked as Received (Skips PO Send)",
			currentStatus: p2pstate.POStatusDraft,
			targetStatus:  p2pstate.POStatusReceived,
			wantErr:       p2pstate.ErrInvalidPOTransition,
		},
		{
			name:          "Invalid: Cancelled PO marked as Paid",
			currentStatus: p2pstate.POStatusCancelled,
			targetStatus:  p2pstate.POStatusPaid,
			wantErr:       p2pstate.ErrInvalidPOTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p2pstate.ValidatePOTransition(tt.currentStatus, tt.targetStatus)
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

// TestP2P_OwnerApprovalThreshold tests dynamic approval gating based on transaction value.
func TestP2P_OwnerApprovalThreshold(t *testing.T) {
	threshold := 10000000.0 // 10 Million IDR

	tests := []struct {
		name          string
		poAmount      float64
		requiresOwner bool
	}{
		{
			name:          "Routine Purchase Below Threshold (3M IDR)",
			poAmount:      3000000,
			requiresOwner: false,
		},
		{
			name:          "Purchase Exactly at Threshold (10M IDR)",
			poAmount:      10000000,
			requiresOwner: true,
		},
		{
			name:          "Major Equipment Purchase Above Threshold (45M IDR)",
			poAmount:      45000000,
			requiresOwner: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := p2pstate.RequiresOwnerApproval(tt.poAmount, threshold)
			if req != tt.requiresOwner {
				t.Errorf("RequiresOwnerApproval(%.2f, %.2f) = %v; want %v", tt.poAmount, threshold, req, tt.requiresOwner)
			}
		})
	}
}

// TestP2P_ThreeWayMatch tests financial reconciliation across PO, GRN, and Vendor Invoice.
func TestP2P_ThreeWayMatch(t *testing.T) {
	tests := []struct {
		name             string
		poQty            float64
		poUnitPrice      float64
		grnReceivedQty   float64
		invBilledQty     float64
		invUnitPrice     float64
		tolerancePercent float64
		wantStatus       p2pstate.ThreeWayMatchStatus
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
			wantStatus:       p2pstate.MatchSuccess,
			wantErr:          nil,
		},
		{
			name:             "Price Variance Within Allowed Tolerance (1.5% increase with 2% tolerance)",
			poQty:            50,
			poUnitPrice:      100000,
			grnReceivedQty:   50,
			invBilledQty:     50,
			invUnitPrice:     101500, // +1.5%
			tolerancePercent: 2.0,
			wantStatus:       p2pstate.MatchSuccess,
			wantErr:          nil,
		},
		{
			name:             "Price Variance Exceeding Allowed Tolerance (5% increase with 2% tolerance)",
			poQty:            50,
			poUnitPrice:      100000,
			grnReceivedQty:   50,
			invBilledQty:     50,
			invUnitPrice:     105000, // +5%
			tolerancePercent: 2.0,
			wantStatus:       p2pstate.MatchPriceError,
			wantErr:          p2pstate.ErrPriceVarianceExceeded,
		},
		{
			name:             "Over-Billed Quantity: Billed More Than Received (GRN)",
			poQty:            100,
			poUnitPrice:      50000,
			grnReceivedQty:   80,  // Only 80 arrived
			invBilledQty:     100, // Vendor tried to bill for 100
			invUnitPrice:     50000,
			tolerancePercent: 0.0,
			wantStatus:       p2pstate.MatchQuantityError,
			wantErr:          p2pstate.ErrOverBilledQuantity,
		},
		{
			name:             "Billed Quantity Exceeds PO Quantity",
			poQty:            100,
			poUnitPrice:      50000,
			grnReceivedQty:   120, // Warehouse received extra unauthorized items
			invBilledQty:     120, // Vendor billed 120
			invUnitPrice:     50000,
			tolerancePercent: 0.0,
			wantStatus:       p2pstate.MatchQuantityError,
			wantErr:          p2pstate.ErrPOQuantityExceeded,
		},
		{
			name:             "Valid Partial Delivery Matching",
			poQty:            100,
			poUnitPrice:      50000,
			grnReceivedQty:   40, // 40 received in first batch
			invBilledQty:     40, // Vendor bills for 40
			invUnitPrice:     50000,
			tolerancePercent: 0.0,
			wantStatus:       p2pstate.MatchSuccess,
			wantErr:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := p2pstate.ValidateThreeWayMatch(
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
