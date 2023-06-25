package order

import "context"

type Repository interface {
	CreateOne(ctx context.Context, c *OrderDto) error
	FindAll(ctx context.Context) (c []OrderDto, err error)
	CreateAll(ctx context.Context, orders []*OrderDto) error
	FindByLimitAndOffset(ctx context.Context, l, o int) (order []OrderDto, err error)
	FindOne(ctx context.Context, id int) (Order, error)
	Update(ctx context.Context, o Order, courierId int) error
	Delete(ctx context.Context, id int) error
}
