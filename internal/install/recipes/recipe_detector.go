package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
	recipeEvaluated  map[string]bool // same recipe(ref) should only be evaluated one time
}

func newRecipeDetector(processEvaluator DetectionStatusProvider, scriptEvaluator DetectionStatusProvider) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: processEvaluator,
		scriptEvaluator:  scriptEvaluator,
		recipeEvaluated:  make(map[string]bool),
	}
}

func NewRecipeDetector() *RecipeDetector {
	return newRecipeDetector(NewProcessEvaluator(), NewScriptEvaluator())
}

func (dt *RecipeDetector) detectBundleRecipe(ctx context.Context, bundleRecipe *BundleRecipe) {

	// if already evaluated
	if dt.recipeEvaluated[bundleRecipe.Recipe.Name] {
		return
	}

	dt.recipeEvaluated[bundleRecipe.Recipe.Name] = true

	for i := 0; i < len(bundleRecipe.Dependencies); i++ {
		dependencyBundleRecipe := bundleRecipe.Dependencies[i]
		dt.detectBundleRecipe(ctx, dependencyBundleRecipe)
	}

	status := dt.detectRecipe(ctx, bundleRecipe.Recipe)
	bundleRecipe.AddDetectionStatus(status)
}

func (dt *RecipeDetector) detectRecipe(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {

	status := dt.processEvaluator.DetectionStatus(ctx, recipe)

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(ctx, recipe)
	}

	return status
}