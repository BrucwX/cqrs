package ordercheck

import "context"

func (s *Service) OrderBelongCustomer(ctx context.Context, customerID int64, orderID int64) (bool, error) {
	customer, err := s.customerCommand.MustGetCustomer(ctx, customerID)
	if err != nil {
		return false, err
	}

	order, err := s.orderCommand.MustGetOrder(ctx, orderID)
	if err != nil {
		return false, err
	}

	return customer.ID() == order.CustomerID(), nil
}
