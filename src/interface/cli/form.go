package main

import (
	"github.com/charmbracelet/huh"
	"github.com/wizenheimer/cascade/internal/config"
)

// createScenarioGroup creates a group of inputs for the scenario configuration
func createScenarioGroup(config *config.Config) *huh.Group {
	scenarioGroup := huh.NewGroup(
		huh.NewInput().
			Title("Scenario ID").
			Description("Name of the chaos experiment").
			Value(&config.Scenario.ID),
		huh.NewInput().
			Title("Scenario Description").
			Description("Description of the chaos experiment").
			Value(&config.Scenario.Description),
	)
	return scenarioGroup
}

// createTargetGroup creates a group of inputs for the target configuration
func createTargetGroup(config *config.Config) *huh.Group {
	targetGroup := huh.NewGroup(
		huh.NewInput().
			Title("Target Namespaces").
			Description("Namespace or set of namespaces to target").
			Value(&config.Target.Namespaces),
		huh.NewInput().
			Title("Included Pod Names").
			Description("Pods would become chaos experiment target if they contain this string").
			Value(&config.Target.IncludedPodNames),
		huh.NewInput().
			Title("Included Node Names").
			Description("Pods would not become target if they reside on a node with this string").
			Value(&config.Target.IncludedNodeNames),
		huh.NewInput().
			Title("Excluded Pod Names").
			Description("Pods would be spared if they contain this string").
			Value(&config.Target.ExcludedPodNames),
	)

	return targetGroup
}

// createRuntimeGroup creates a group of inputs for the runtime configuration
func createRuntimeGroup(config *config.Config) *huh.Group {
	runtimeGroup := huh.NewGroup(
		huh.NewInput().
			Title("Runtime Interval").
			Description("Intervals at which the chaos experiments are to be triggered").
			Value(&config.Runtime.Interval),
		huh.NewInput().
			Title("Runtime Grace").
			Description("The grace time before the chaos experiment starts").
			Value(&config.Runtime.Grace),
		huh.NewSelect[string]().
			Title("Runtime Mode").
			Description("The mode of the chaos experiment").
			Options(huh.NewOptions("evict", "delete", "dry-run")...).
			Value(&config.Runtime.Mode),
		huh.NewSelect[string]().
			Title("Runtime Ordering").
			Description("Determine the priority of the pods to be targeted,").
			Options(huh.NewOptions("oldest", "youngest", "cost", "random", "default")...).
			Value(&config.Runtime.Ordering),
		huh.NewInput().
			Title("Runtime Ratio").
			Description("Ratio of candidate pods to be targeted").
			Value(&config.Runtime.Ratio),
	)

	return runtimeGroup
}

// createClusterGroup creates a group of inputs for the cluster configuration
func createClusterGroup(config *config.Config) *huh.Group {
	clusterGroup := huh.NewGroup(
		huh.NewInput().
			Title("Cluster Kubeconfig").
			Description("Path to the kubeconfig file").
			Value(&config.Cluster.Kubeconfig),
		huh.NewInput().
			Title("Cluster Master").
			Description("Path to the master configuration file").
			Value(&config.Cluster.Master),
		huh.NewSelect[string]().
			Title("Cluster Origin").
			Description("The origin of the chaos experiment").
			Options(huh.NewOptions("host", "cluster")...).
			Value(&config.Cluster.Origin),
		huh.NewInput().
			Title("Cluster Healthcheck").
			Description("Health check port for the pods").
			Value(&config.Cluster.Healthcheck),
	)

	return clusterGroup
}

// createFilePickerGroup creates a group of inputs for the file picker configuration
func createFilePickerGroup(outputFolder, fileName *string, createFolder *bool) *huh.Group {
	filePickerGroup := huh.NewGroup(
		huh.NewInput().
			Title("Export Folder").
			Description("Enter the folder path to save the YAML file").
			Value(outputFolder),
		huh.NewInput().
			Title("File Name").
			Description("Enter the name for the YAML file (default: config.yaml)").
			Value(fileName),
		huh.NewConfirm().
			Title("Create Folder").
			Description("Create the folder if it doesn't exist?").
			Value(createFolder),
	)

	return filePickerGroup
}

// createForm creates a form with all the input groups
func createScenarioForm(config *config.Config, outputFolder, fileName *string, createFolder *bool) *huh.Form {

	scenarioGroup := createScenarioGroup(config)
	targetGroup := createTargetGroup(config)
	runtimeGroup := createRuntimeGroup(config)
	clusterGroup := createClusterGroup(config)
	filePickerGroup := createFilePickerGroup(outputFolder, fileName, createFolder)

	form := huh.NewForm(
		scenarioGroup,
		targetGroup,
		runtimeGroup,
		clusterGroup,
		filePickerGroup,
	)

	return form
}

func createScenarioPickerGroup(outputPath *string) *huh.Group {
	return huh.NewGroup(
		huh.NewFilePicker().
			Title("Export Location").
			Description("Choose where to save the YAML file").
			Value(outputPath),
	)
}

func createSessionForm(outputPath *string) *huh.Form {
	scenarioPickerGroup := createScenarioPickerGroup(outputPath)

	form := huh.NewForm(
		scenarioPickerGroup,
	)

	return form
}
