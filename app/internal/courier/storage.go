package courier

import "context"

type Repository interface {
	CreateOne(ctx context.Context, c *CourierDto) error
	FindAll(ctx context.Context) (c []CourierDto, err error)
	CreateAll(ctx context.Context, couriers []*CourierDto) error
	FindByLimitAndOffset(ctx context.Context, l, o int) (c []CourierDto, err error)
	FindOne(ctx context.Context, id int) (CourierDto, error)
	Update(ctx context.Context, courier CourierDto) error
	Delete(ctx context.Context, id int) error
}
