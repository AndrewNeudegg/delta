package definitions

import (
	"context"
)

// Identification for this pipeline.
type Identification interface {
	ID() string // ID defines what this thing is.
}

// Pipeline defines the construction of I->P->O.
type Pipeline interface {
	Identification

	Do(context.Context) error
}
