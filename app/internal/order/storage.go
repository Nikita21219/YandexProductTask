package order

import "context"

type Repository interface {
	CreateOne(ctx context.Context, c *OrderDto) error
	FindAll(ctx context.Context) (c []OrderDto, err error)
	CreateAll(ctx context.Context, couriers []*OrderDto) error
	FindByLimitAndOffset(ctx context.Context, l, o int) (c []OrderDto, err error)
	FindOne(ctx context.Context, id int) (OrderDto, error)
	Update(ctx context.Context, courier OrderDto) error
	Delete(ctx context.Context, id int) error
}
