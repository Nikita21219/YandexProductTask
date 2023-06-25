package order

import "database/sql"

type Order struct {
	Id           int
	CourierId    sql.NullInt64
	Weight       int
	Region       int
	DeliveryTime string
	Price        int
}
