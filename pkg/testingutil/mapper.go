package testingutil

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
)

// Mapper is a redefinition for the single method that we need to avoid import cycles
type Mapper interface {
	ToResource(context.Context, any, string) (model.Resource, error)
}

// ConvertToResources converts all the raw cloud SDK resources to `model.Resource`s using the given mapper.
func ConvertToResources[T any](t testing.TB, ctx context.Context, mapper Mapper, raw []T) []model.Resource {
	t.Helper()
	output := make([]model.Resource, 0, len(raw))
	for _, in := range raw {
		resource, err := mapper.ToResource(ctx, in, TestRegion)
		if err != nil {
			t.Fatalf("cannot convert resource: %v", err)
		}

		output = append(output, resource)
	}

	// Make sure we only grab resources meant for integration testing
	output = ResourceFilterTagKeyValue(output, "IntegrationTest", "true")

	return output
}