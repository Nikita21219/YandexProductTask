package courier

import "context"

type Repository interface {
	CreateOne(ctx context.Context, c *CourierDto) error
	FindAll(ctx context.Context) (c []Courier, err error)
	CreateAll(ctx context.Context, couriers []*CourierDto) error
	FindByLimitAndOffset(ctx context.Context, l, o int) (c []Courier, err error)
	FindOne(ctx context.Context, id int) (Courier, error)
	Update(ctx context.Context, courier Courier) error
	Delete(ctx context.Context, id int) error
}
