package courier

import "regexp"

type CourierType string

const (
	Foot CourierType = "FOOT"
	Bike CourierType = "BIKE"
	Auto CourierType = "AUTO"
)

type CourierDto struct {
	Id           int         `json:"id"`
	CourierType  CourierType `json:"courier_type"`
	Regions      []int       `json:"regions"`
	WorkingHours []string    `json:"working_hours"`
}

func (c *CourierDto) Valid() (bool, error) {
	switch c.CourierType {
	case Foot, Bike, Auto:
	default:
		return false, nil
	}

	if c.WorkingHours == nil || c.Regions == nil {
		return false, nil
	}

	for _, hours := range c.WorkingHours {
		match, err := regexp.MatchString(`^\d{2}:\d{2}-\d{2}:\d{2}$`, hours)
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}

	for _, region := range c.Regions {
		if region < 0 {
			return false, nil
		}
	}

	return true, nil
}
