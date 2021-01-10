package noop

import "math/rand"

const (
	// ID for this collection of resources.
	ID = "utilities/chaos/v1"
)

func isLucky(f float32) bool {
	// rand.Float64() == 0.1, 0.5, 0.8
	// if f == 0.1 (10% chance) then rand has to be above 0.9.
	// if f == 0.90 (90% chance) then rand has to be above 0.1.
	return rand.Float32() > f
}
