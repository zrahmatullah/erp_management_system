package p2pstate

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidPRTransition   = errors.New("invalid purchase requisition status transition")
	ErrInvalidPOTransition   = errors.New("invalid purchase order status transition")
	ErrOverBilledQuantity    = errors.New("billed quantity exceeds goods received note (GRN) quantity")
	ErrPOQuantityExceeded    = errors.New("billed quantity exceeds purchase order quantity")
	ErrPriceVarianceExceeded = errors.New("invoice unit price exceeds PO price beyond allowed tolerance")
)

// PR Statuses
const (
	PRStatusDraft           = "draft"
	PRStatusPendingApproval = "pending_approval"
	PRStatusApproved        = "approved"
	PRStatusRejected        = "rejected"
	PRStatusPOCreated       = "po_created"
	PRStatusCancelled       = "cancelled"
)

// Allowed PR state transitions
var validPRTransitions = map[string][]string{
	PRStatusDraft:           {PRStatusPendingApproval, PRStatusCancelled},
	PRStatusPendingApproval: {PRStatusApproved, PRStatusRejected, PRStatusCancelled},
	PRStatusApproved:        {PRStatusPOCreated, PRStatusCancelled},
	PRStatusRejected:        {PRStatusDraft}, // Allows edit & resubmit
	PRStatusPOCreated:       {},              // Terminal state
	PRStatusCancelled:       {},              // Terminal state
}

// ValidatePRTransition validates if a PR status transition is legal.
func ValidatePRTransition(currentStatus, targetStatus string) error {
	if currentStatus == targetStatus {
		return nil
	}

	allowed, exists := validPRTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("%w: unknown current status '%s'", ErrInvalidPRTransition, currentStatus)
	}

	for _, next := range allowed {
		if next == targetStatus {
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from '%s' to '%s'", ErrInvalidPRTransition, currentStatus, targetStatus)
}

// PO Statuses
const (
	POStatusDraft             = "draft"
	POStatusSubmitted         = "submitted"
	POStatusManagerApproved   = "manager_approved"
	POStatusOwnerApproved     = "owner_approved"
	POStatusRejected          = "rejected"
	POStatusSent              = "sent"
	POStatusPartiallyReceived = "partially_received"
	POStatusReceived          = "received"
	POStatusInvoiced          = "invoiced"
	POStatusPaid              = "paid"
	POStatusCancelled         = "cancelled"
)

// Allowed PO state transitions
var validPOTransitions = map[string][]string{
	POStatusDraft:             {POStatusSubmitted, POStatusCancelled},
	POStatusSubmitted:         {POStatusManagerApproved, POStatusRejected, POStatusCancelled},
	POStatusManagerApproved:   {POStatusOwnerApproved, POStatusSent, POStatusRejected, POStatusCancelled},
	POStatusOwnerApproved:     {POStatusSent, POStatusRejected, POStatusCancelled},
	POStatusSent:              {POStatusPartiallyReceived, POStatusReceived, POStatusCancelled},
	POStatusPartiallyReceived: {POStatusReceived, POStatusCancelled},
	POStatusReceived:          {POStatusInvoiced},
	POStatusInvoiced:          {POStatusPaid},
	POStatusPaid:              {}, // Terminal
	POStatusRejected:          {POStatusDraft}, // Can revise
	POStatusCancelled:         {}, // Terminal
}

// ValidatePOTransition validates if a PO status transition is legal.
func ValidatePOTransition(currentStatus, targetStatus string) error {
	if currentStatus == targetStatus {
		return nil
	}

	allowed, exists := validPOTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("%w: unknown current status '%s'", ErrInvalidPOTransition, currentStatus)
	}

	for _, next := range allowed {
		if next == targetStatus {
			return nil
		}
	}

	return fmt.Errorf("%w: cannot transition from '%s' to '%s'", ErrInvalidPOTransition, currentStatus, targetStatus)
}

// RequiresOwnerApproval determines whether a PO amount warrants Owner approval.
func RequiresOwnerApproval(poTotalAmount, threshold float64) bool {
	return poTotalAmount >= threshold
}

// ThreeWayMatchStatus represents reconciliation outcome.
type ThreeWayMatchStatus string

const (
	MatchSuccess       ThreeWayMatchStatus = "matched"
	MatchQuantityError ThreeWayMatchStatus = "quantity_mismatch"
	MatchPriceError    ThreeWayMatchStatus = "price_mismatch"
)

// ThreeWayMatchResult reports audit matching details between PO, GRN, and Vendor Invoice.
type ThreeWayMatchResult struct {
	Status           ThreeWayMatchStatus `json:"status"`
	POQuantity       float64             `json:"po_quantity"`
	ReceivedQuantity float64             `json:"received_quantity"`
	BilledQuantity   float64             `json:"billed_quantity"`
	POUnitPrice      float64             `json:"po_unit_price"`
	BilledUnitPrice  float64             `json:"billed_unit_price"`
	PriceVariancePct float64             `json:"price_variance_pct"`
	AllowedVariance  float64             `json:"allowed_variance_pct"`
	ErrorMessage     string              `json:"error_message,omitempty"`
}

// ValidateThreeWayMatch verifies Purchase Order vs Goods Receipt vs Vendor Invoice.
func ValidateThreeWayMatch(
	poQty, poUnitPrice float64,
	grnReceivedQty float64,
	invBilledQty, invUnitPrice float64,
	tolerancePercent float64,
) (ThreeWayMatchResult, error) {
	result := ThreeWayMatchResult{
		Status:           MatchSuccess,
		POQuantity:       poQty,
		ReceivedQuantity: grnReceivedQty,
		BilledQuantity:   invBilledQty,
		POUnitPrice:      poUnitPrice,
		BilledUnitPrice:  invUnitPrice,
		AllowedVariance:  tolerancePercent,
	}

	// 1. Quantity validation: Cannot invoice more than received
	if invBilledQty > grnReceivedQty {
		result.Status = MatchQuantityError
		result.ErrorMessage = fmt.Sprintf("Billed quantity (%.2f) exceeds received quantity (%.2f)", invBilledQty, grnReceivedQty)
		return result, ErrOverBilledQuantity
	}

	// 2. Quantity validation: Cannot invoice more than ordered
	if invBilledQty > poQty {
		result.Status = MatchQuantityError
		result.ErrorMessage = fmt.Sprintf("Billed quantity (%.2f) exceeds PO quantity (%.2f)", invBilledQty, poQty)
		return result, ErrPOQuantityExceeded
	}

	// 3. Price validation: Check variance percentage
	if poUnitPrice > 0 {
		variancePct := ((invUnitPrice - poUnitPrice) / poUnitPrice) * 100.0
		variancePct = math.Round(variancePct*100) / 100
		result.PriceVariancePct = variancePct

		if variancePct > tolerancePercent {
			result.Status = MatchPriceError
			result.ErrorMessage = fmt.Sprintf("Invoice price (%.2f) exceeds PO price (%.2f) by %.2f%% (limit: %.2f%%)",
				invUnitPrice, poUnitPrice, variancePct, tolerancePercent)
			return result, ErrPriceVarianceExceeded
		}
	}

	return result, nil
}
