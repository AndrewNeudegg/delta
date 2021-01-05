package definitions

import (
	"context"

	resourceDefinitions "github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Pipeline defines the construction of I->P->O.
type Pipeline interface {
	resourceDefinitions.Resource

	Do(context.Context)
}
