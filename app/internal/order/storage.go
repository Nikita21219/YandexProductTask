package order

import (
	"context"
	"main/internal/order_complete"
)

type Repository interface {
	CreateOne(ctx context.Context, c *Order) error
	FindAll(ctx context.Context) (c []Order, err error)
	CreateAll(ctx context.Context, orders []*Order) error
	FindByLimitAndOffset(ctx context.Context, l, o int) (order []Order, err error)
	FindOne(ctx context.Context, id int) (Order, error)
	Update(ctx context.Context, o Order, oc *order_complete.OrderCompleteDto) error
	Delete(ctx context.Context, id int) error
	FindAllInTimeInterval(ctx context.Context, startDate, endDate string, courierId int) ([]Order, error)
}
