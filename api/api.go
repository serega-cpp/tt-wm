package api

import (
	"fmt"
	"math"
)

type City struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
}

func (c *City) Validate() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("City: Name is empty")
	}
	if math.Abs(c.Latitude) > 90.0 {
		return fmt.Errorf("Latitude: Wrong value (allowed -90..90)")
	}
	if math.Abs(c.Longitude) > 180.0 {
		return fmt.Errorf("Longitude: Wrong value (allowed -180..180)")
	}
	return nil
}

type Temperature struct {
	ID        int
	CityID    int
	MaxC      float32
	MinC      float32
	Timestamp int64
}

func (t *Temperature) Validate() error {
	if t.CityID <= 0 {
		return fmt.Errorf("Temperature: Wrong city_id (must be > 0)")
	}
	if t.MinC < -273.0 || 100.0 < t.MinC {
		return fmt.Errorf("Temperature: Wrong min temperature: %f", t.MinC)
	}
	if t.MaxC < -273.0 || 100.0 < t.MaxC {
		return fmt.Errorf("Temperature: Wrong max temperature: %f", t.MaxC)
	}
	return nil
}
