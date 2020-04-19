package types

import (
	"github.com/havr/wea/pkg/util"
)

// Temperature represents temperature in Kelvins.
// To convert it to Celsius use the corresponding method.
type Temperature float64

// Celsius returns the temperature converted to Celsius scale.
func (t Temperature) Celsius() float64 {
	return util.Round(-273.15+float64(t), 2)
}
