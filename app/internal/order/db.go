package order

import (
	"context"
	"fmt"
	"main/internal/order_complete"
	"main/pkg"
	"strings"
	"time"
)

type repository struct {
	client pkg.DBClient
}

func (r *repository) CreateOne(ctx context.Context, c *Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) CreateAll(ctx context.Context, orders []*Order) error {
	q := `INSERT INTO "order" (weight, region, delivery_time, price) VALUES %s`

	values := make([]string, 0, 4)
	params := make([]interface{}, 0, len(orders))

	for _, o := range orders {
		paramsLength := len(params)
		rowValues := fmt.Sprintf(
			"($%d, $%d, $%d, $%d)",
			paramsLength+1,
			paramsLength+2,
			paramsLength+3,
			paramsLength+4,
		)
		values = append(values, rowValues)
		params = append(params, o.Weight, o.Region, o.DeliveryTime, o.Price)
	}

	q = fmt.Sprintf(q, strings.Join(values, ","))
	_, err := r.client.Exec(ctx, q, params...)
	return err
}

func (r *repository) FindAll(ctx context.Context) ([]Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindOne(ctx context.Context, id int) (Order, error) {
	q := `SELECT id, courier_id, weight, region, delivery_time, price FROM "order" WHERE id = $1`
	var o Order
	if err := r.client.QueryRow(ctx, q, id).Scan(&o.Id, &o.CourierId, &o.Weight, &o.Region, &o.DeliveryTime, &o.Price); err != nil {
		return Order{}, err
	}
	return o, nil
}

func (r *repository) FindByLimitAndOffset(ctx context.Context, l, o int) ([]Order, error) {
	q := `SELECT id, weight, region, delivery_time, price FROM "order" ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.client.Query(ctx, q, l, o)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0, l)

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.Id, &order.Weight, &order.Region, &order.DeliveryTime, &order.Price)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) Update(ctx context.Context, o Order, oc *order_complete.OrderCompleteDto) error {
	q := `UPDATE "order" SET courier_id = $1, complete_time = $2 WHERE id = $3`
	completeTime, err := time.Parse("2006-01-02", oc.OrderCompleteTime)
	if err != nil {
		return err
	}
	_, err = r.client.Exec(ctx, q, oc.CourierId, completeTime, o.Id)
	return err
}

func (r *repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindAllInTimeInterval(ctx context.Context, startDate, endDate string, courierId int) ([]Order, error) {
	q := `
	SELECT id, courier_id, weight, region, delivery_time, price
	FROM "order" 
	WHERE delivery_time >= $1 AND delivery_time < $2 AND courier_id = $3
	`

	rows, err := r.client.Query(ctx, q, startDate, endDate, courierId)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)

	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Id, &o.CourierId, &o.Weight, &o.Region, &o.DeliveryTime, &o.Price)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func NewRepo(client pkg.DBClient) Repository {
	return &repository{
		client: client,
	}
}
