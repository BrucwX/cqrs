package order

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the order status.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCanceled  OrderStatus = "CANCELED"
	OrderStatusRefunded  OrderStatus = "REFUNDED"
)

// Order represents an order. It is the order aggregate root.
type Order struct {
	ID        string
	StudentID string
	Amount    float64
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewOrder creates a new order.
func NewOrder(studentID string, amount float64) *Order {
	return &Order{
		ID:        uuid.New().String(),
		StudentID: studentID,
		Amount:    amount,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MarkPaid marks the order as paid.
func (o *Order) MarkPaid() {
	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()
}

// Cancel cancels the order.
func (o *Order) Cancel() {
	o.Status = OrderStatusCanceled
	o.UpdatedAt = time.Now()
}

// Refund marks the order as refunded.
func (o *Order) Refund() {
	o.Status = OrderStatusRefunded
	o.UpdatedAt = time.Now()
}

// IsPaid checks if the order is paid.
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}
