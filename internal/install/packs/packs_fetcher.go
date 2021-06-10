package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// PacksFetcher is responsible for retrieving packs information.
type PacksFetcher interface {
	FetchPacks(context.Context, []types.OpenInstallationRecipe) ([]types.OpenInstallationObservabilityPack, error)
}
