package event

import (
	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/orderItem"
)

type OrderPlacedEvent struct {
	OrderItems orderItem.OrderItems
	Customer   customer.Customer
}
