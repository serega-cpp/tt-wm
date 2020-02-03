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
