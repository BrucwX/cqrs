package payment

import (
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the payment method.
type PaymentMethod string

const (
	PaymentMethodWechat  PaymentMethod = "WECHAT"
	PaymentMethodAlipay  PaymentMethod = "ALIPAY"
	PaymentMethodCash    PaymentMethod = "CASH"
	PaymentMethodBank    PaymentMethod = "BANK_TRANSFER"
)

// PaymentStatus represents the payment status.
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusFailed  PaymentStatus = "FAILED"
	PaymentStatusRefund  PaymentStatus = "REFUND"
)

// Payment represents a payment. It is the payment aggregate root.
type Payment struct {
	ID        string
	OrderID   string
	StudentID string
	Amount    float64
	Method    PaymentMethod
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewPayment creates a new payment.
func NewPayment(orderID, studentID string, amount float64, method PaymentMethod) *Payment {
	return &Payment{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		StudentID: studentID,
		Amount:    amount,
		Method:    method,
		Status:    PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MarkPaid marks the payment as paid.
func (p *Payment) MarkPaid() {
	p.Status = PaymentStatusPaid
	p.UpdatedAt = time.Now()
}

// MarkFailed marks the payment as failed.
func (p *Payment) MarkFailed() {
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
}

// Refund marks the payment as refunded.
func (p *Payment) Refund() {
	p.Status = PaymentStatusRefund
	p.UpdatedAt = time.Now()
}

// IsPaid checks if the payment is paid.
func (p *Payment) IsPaid() bool {
	return p.Status == PaymentStatusPaid
}
