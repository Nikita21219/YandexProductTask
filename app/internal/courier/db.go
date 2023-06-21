package courier

import (
	"context"
	"fmt"
	"main/pkg"
	"strings"
)

type repository struct {
	client pkg.DBClient
}

func (r *repository) CreateOne(ctx context.Context, c *CourierDto) error {
	q := `INSERT INTO courier (courier_type, regions, working_hours) VALUES ($1, $2, $3) RETURNING id`
	row := r.client.QueryRow(ctx, q, c.CourierType, c.Regions, c.WorkingHours)
	if err := row.Scan(&c.Id); err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateAll(ctx context.Context, couriers []*CourierDto) error {
	q := `INSERT INTO courier (courier_type, regions, working_hours) VALUES %s`

	values := make([]string, 0, 3)
	params := make([]interface{}, 0, len(couriers))

	for _, c := range couriers {
		paramsLength := len(params)
		rowValues := fmt.Sprintf("($%d, $%d, $%d)", paramsLength+1, paramsLength+2, paramsLength+3)
		values = append(values, rowValues)
		params = append(params, string(c.CourierType), c.Regions, c.WorkingHours)
	}

	q = fmt.Sprintf(q, strings.Join(values, ","))
	_, err := r.client.Exec(ctx, q, params...)
	return err
}

func (r *repository) FindAll(ctx context.Context) ([]CourierDto, error) {
	q := `SELECT id, courier_type, regions, working_hours FROM courier`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	couriers := make([]CourierDto, 0, 50)

	for rows.Next() {
		var c CourierDto
		err = rows.Scan(&c.Id, &c.CourierType, &c.Regions, &c.WorkingHours)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (CourierDto, error) {
	q := `SELECT id, courier_type, regions, working_hours FROM courier WHERE id = $1`
	var c CourierDto
	if err := r.client.QueryRow(ctx, q, id).Scan(&c.Id, &c.CourierType, &c.Regions, &c.WorkingHours); err != nil {
		return CourierDto{}, err
	}
	return c, nil
}

func (r *repository) FindByLimitAndOffset(ctx context.Context, l, o int) (c []CourierDto, err error) {
	q := `SELECT id, courier_type, regions, working_hours FROM courier ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.client.Query(ctx, q, l, o)
	if err != nil {
		return nil, err
	}

	couriers := make([]CourierDto, 0, l)

	for rows.Next() {
		var c CourierDto
		err = rows.Scan(&c.Id, &c.CourierType, &c.Regions, &c.WorkingHours)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *repository) Update(ctx context.Context, courier CourierDto) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewRepo(client pkg.DBClient) Repository {
	return &repository{
		client: client,
	}
}
