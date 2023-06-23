package order_complete

import (
	"encoding/json"
	"main/pkg/utils"
)

type OrderCompleteDto struct {
	CourierId int    `json:"courier_id"`
	OrderId   int    `json:"order_id"`
	OrderTime string `json:"order_time"`
}

func (oc *OrderCompleteDto) Valid() bool {
	return oc.OrderId >= 0 && oc.CourierId >= 0 && utils.ValidTime(oc.OrderTime)
}

func (oc *OrderCompleteDto) MarshalBinary() ([]byte, error) {
	return json.Marshal(oc)
}
