package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validParams = Params{
		Credentials: "gke-production",
		App:         "myapp",
		Namespace:   "mynamespace",
		Autoscale: AutoscaleParams{
			MinReplicas:   3,
			MaxReplicas:   100,
			CPUPercentage: 80,
		},
		Container: ContainerParams{
			ImageRepository: "estafette",
			ImageName:       "my-app",
			ImageTag:        "1.0.0",
			Port:            5000,

			CPU: CPUParams{
				Request: "100m",
				Limit:   "150m",
			},
			Memory: MemoryParams{
				Request: "768Mi",
				Limit:   "1024Mi",
			},
			LivenessProbe: ProbeParams{
				Path:                "/liveness",
				InitialDelaySeconds: 30,
				TimeoutSeconds:      1,
			},
			ReadinessProbe: ProbeParams{
				Path:                "/readiness",
				InitialDelaySeconds: 0,
				TimeoutSeconds:      1,
			},
			Metrics: MetricsParams{
				Scrape: "true",
				Path:   "/metrics",
				Port:   5000,
			},
		},
		Visibility: "private",
		Hosts:      []string{"gke.estafette.io"},
	}
)

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsAppToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			App: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "myapp", params.App)
	})

	t.Run("KeepsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			App: "yourapp",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "yourapp", params.App)
	})

	t.Run("DefaultsImageNameToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "",
			},
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "myapp", params.Container.ImageName)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "my-app",
			},
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "my-app", params.Container.ImageName)
	})

	t.Run("DefaultsImageTagToBuildVersionIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "",
			},
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "", map[string]string{})

		assert.Equal(t, "1.0.0", params.Container.ImageTag)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "2.1.3",
			},
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "", map[string]string{})

		assert.Equal(t, "2.1.3", params.Container.ImageTag)
	})

	t.Run("DefaultsCredentialsToReleaseNamePrefixedByGKEIfEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName, map[string]string{})

		assert.Equal(t, "gke-production", params.Credentials)
	})

	t.Run("KeepsCredentialsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "staging",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName, map[string]string{})

		assert.Equal(t, "staging", params.Credentials)
	})

	t.Run("DefaultsLabelsToEstafetteLabelsIfEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

	t.Run("KeepsLabelsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"app":  "yourapp",
				"team": "yourteam",
			},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("AddsAppLabelToLabelsIfNotSet", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"team": "yourteam",
			},
		}
		appLabel := "myapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults(appLabel, "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("OverwritesAppLabelToAppIfSetFromEstafetteLabels", func(t *testing.T) {

		params := Params{}
		appLabel := "yourapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults(appLabel, "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

	t.Run("DefaultsVisibilityToPrivateIfEmpty", func(t *testing.T) {

		params := Params{
			Visibility: "",
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "private", params.Visibility)
	})

	t.Run("KeepsVisibilityIfNotEmpty", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "public", params.Visibility)
	})

	t.Run("DefaultsCpuRequestTo100MIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "100m", params.Container.CPU.Request)
	})

	t.Run("DefaultsCpuRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "300m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "300m", params.Container.CPU.Request)
	})

	t.Run("KeepsCpuRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "250m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Container.CPU.Request)
	})

	t.Run("DefaultsCpuLimitTo125MIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "125m", params.Container.CPU.Limit)
	})

	t.Run("DefaultsCpuLimitToRequestIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "300m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "300m", params.Container.CPU.Limit)
	})

	t.Run("KeepsCpuLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "250m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Container.CPU.Limit)
	})

	t.Run("DefaultsMemoryRequestTo128MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "128Mi", params.Container.Memory.Request)
	})

	t.Run("DefaultsMemoryRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "256Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "256Mi", params.Container.Memory.Request)
	})

	t.Run("KeepsMemoryRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "512Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "512Mi", params.Container.Memory.Request)
	})

	t.Run("DefaultsMemoryLimitTo128MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "128Mi", params.Container.Memory.Limit)
	})

	t.Run("DefaultsMemoryLimitToRequestIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "768Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "768Mi", params.Container.Memory.Limit)
	})

	t.Run("KeepsMemoryLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "1024Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "1024Mi", params.Container.Memory.Limit)
	})

	t.Run("DefaultsContainerPortTo5000IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 5000, params.Container.Port)
	})

	t.Run("KeepsContainerPortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 3000,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 3000, params.Container.Port)
	})

	t.Run("DefaultsAutoscaleMinReplicasTo3IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 3, params.Autoscale.MinReplicas)
	})

	t.Run("KeepsAutoscaleMinReplicasIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 2,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 2, params.Autoscale.MinReplicas)
	})

	t.Run("DefaultsAutoscaleMaxReplicasTo100IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 100, params.Autoscale.MaxReplicas)
	})

	t.Run("KeepsAutoscaleMaxReplicasIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 50,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 50, params.Autoscale.MaxReplicas)
	})

	t.Run("DefaultsAutoscaleCPUPercentageTo80IfZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				CPUPercentage: 0,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 80, params.Autoscale.CPUPercentage)
	})

	t.Run("KeepsAutoscaleCPUPercentageIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				CPUPercentage: 30,
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 30, params.Autoscale.CPUPercentage)
	})

	t.Run("DefaultsLivenessInitialDelaySecondsTo30IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					InitialDelaySeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 30, params.Container.LivenessProbe.InitialDelaySeconds)
	})

	t.Run("KeepsLivenessInitialDelaySecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					InitialDelaySeconds: 120,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 120, params.Container.LivenessProbe.InitialDelaySeconds)
	})

	t.Run("DefaultsLivenessTimeoutSecondsTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					TimeoutSeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.LivenessProbe.TimeoutSeconds)
	})

	t.Run("KeepsLivenessTimeoutSecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					TimeoutSeconds: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.LivenessProbe.TimeoutSeconds)
	})

	t.Run("DefaultsLivenessPathToLivenessIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/liveness", params.Container.LivenessProbe.Path)
	})

	t.Run("KeepsLivenessPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Path: "/healthz",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/healthz", params.Container.LivenessProbe.Path)
	})

	t.Run("DefaultsReadinessInitialDelaySecondsTo0IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					InitialDelaySeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 0, params.Container.ReadinessProbe.InitialDelaySeconds)
	})

	t.Run("KeepsReadinessInitialDelaySecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					InitialDelaySeconds: 120,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 120, params.Container.ReadinessProbe.InitialDelaySeconds)
	})

	t.Run("DefaultsReadinessTimeoutSecondsTo1IfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					TimeoutSeconds: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 1, params.Container.ReadinessProbe.TimeoutSeconds)
	})

	t.Run("KeepsReadinessTimeoutSecondsIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					TimeoutSeconds: 5,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 5, params.Container.ReadinessProbe.TimeoutSeconds)
	})

	t.Run("DefaultsReadinessPathToReadinessIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/readiness", params.Container.ReadinessProbe.Path)
	})

	t.Run("KeepsReadinessPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "/healthz",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/healthz", params.Container.ReadinessProbe.Path)
	})

	t.Run("DefaultsMetricsPathToMetricsIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/metrics", params.Container.Metrics.Path)
	})

	t.Run("KeepsMetricsPathIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "/mymetrics",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "/mymetrics", params.Container.Metrics.Path)
	})

	t.Run("DefaultsMetricsPortToContainerPortIfZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 5000,
				Metrics: MetricsParams{
					Port: 0,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 5000, params.Container.Metrics.Port)
	})

	t.Run("KeepsMetricsPortIfLargerThanZero", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 5000,
				Metrics: MetricsParams{
					Port: 5001,
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, 5001, params.Container.Metrics.Port)
	})

	t.Run("DefaultsMetricsScrapeToTrueIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Scrape: "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "true", params.Container.Metrics.Scrape)
	})

	t.Run("KeepsMetricsScrapeIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Scrape: "false",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "false", params.Container.Metrics.Scrape)
	})

	t.Run("DefaultsSidecarTypeToOpenrestyIfEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: "",
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "openresty", params.Sidecar.Type)
	})

	t.Run("KeepsVisibilityIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: "istio",
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "istio", params.Sidecar.Type)
	})

	t.Run("DefaultsSidecarCpuRequestTo10MIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "10m", params.Sidecar.CPU.Request)
	})

	t.Run("DefaultsSidecarCpuRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "300m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "300m", params.Sidecar.CPU.Request)
	})

	t.Run("KeepsSidecarCpuRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "250m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Sidecar.CPU.Request)
	})

	t.Run("DefaultsSidecarCpuLimitTo50MIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "50m", params.Sidecar.CPU.Limit)
	})

	t.Run("DefaultsSidecarCpuLimitToRequestIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "300m",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "300m", params.Sidecar.CPU.Limit)
	})

	t.Run("KeepsSidecarCpuLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "",
					Limit:   "250m",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "250m", params.Sidecar.CPU.Limit)
	})

	t.Run("DefaultsSidecarMemoryRequestTo10MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "10Mi", params.Sidecar.Memory.Request)
	})

	t.Run("DefaultsSidecarMemoryRequestToLimitIfRequestIsEmptyButLimitIsNot", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "256Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "256Mi", params.Sidecar.Memory.Request)
	})

	t.Run("KeepsSidecarMemoryRequestIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "512Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "512Mi", params.Sidecar.Memory.Request)
	})

	t.Run("DefaultsSidecarMemoryLimitTo50MiIfBothRequestAndLimitAreEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "50Mi", params.Sidecar.Memory.Limit)
	})

	t.Run("DefaultsSidecarMemoryLimitToRequestIfLimitIsEmptyButRequestIsNot", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "768Mi",
					Limit:   "",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "768Mi", params.Sidecar.Memory.Limit)
	})

	t.Run("KeepsSidecarMemoryLimitIfNotEmpty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "",
					Limit:   "1024Mi",
				},
			},
		}

		// act
		params.SetDefaults("", "", "", map[string]string{})

		assert.Equal(t, "1024Mi", params.Sidecar.Memory.Limit)
	})

}

func TestSetDefaultsFromCredentials(t *testing.T) {

	t.Run("DefaultsNamespaceToCredentialDefaultNamespaceIfEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				DefaultNamespace: "mynamespace",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "mynamespace", params.Namespace)
	})

	t.Run("KeepsNamespaceIfNotEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "yournamespace",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				DefaultNamespace: "mynamespace",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "yournamespace", params.Namespace)
	})

	t.Run("DefaultsImageRepositoryToCredentialProjectIfEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageRepository: "",
			},
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				Project: "myproject",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "myproject", params.Container.ImageRepository)
	})

	t.Run("KeepsImageRepositoryIfNotEmpty", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageRepository: "extensions",
			},
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				Project: "myproject",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "extensions", params.Container.ImageRepository)
	})
}

func TestValidateRequiredProperties(t *testing.T) {

	t.Run("ReturnsFalseIfAppIsNotSet", func(t *testing.T) {

		params := validParams
		params.App = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAppIsSet", func(t *testing.T) {

		params := validParams
		params.App = "myapp"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfNamespaceIsNotSet", func(t *testing.T) {

		params := validParams
		params.Namespace = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfNamespaceIsSet", func(t *testing.T) {

		params := validParams
		params.Namespace = "mynamespace"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageRepositoryIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageRepository = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageRepositoryIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageRepository = "myrepository"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageNameIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageName = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageNameIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageName = "myimage"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfImageTagIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageTag = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageTagIsSet", func(t *testing.T) {

		params := validParams
		params.Container.ImageTag = "1.0.0"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfCredentialsIsNotSet", func(t *testing.T) {

		params := validParams
		params.Credentials = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCredentialsIsSet", func(t *testing.T) {

		params := validParams
		params.Credentials = "gke-production"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfVisibilityIsNotSet", func(t *testing.T) {

		params := validParams
		params.Visibility = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfVisibilityIsSetToUnsupportedValue", func(t *testing.T) {

		params := validParams
		params.Visibility = "everywhere"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToPublic", func(t *testing.T) {

		params := validParams
		params.Visibility = "public"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfVisibilityIsSetToPrivate", func(t *testing.T) {

		params := validParams
		params.Visibility = "private"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfCpuRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Request = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCpuRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Request = "100m"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfCpuLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Limit = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCpuLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Container.CPU.Limit = "100m"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMemoryRequestIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Request = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMemoryRequestIsSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Request = "100m"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMemoryLimitIsNotSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Limit = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMemoryLimitIsSet", func(t *testing.T) {

		params := validParams
		params.Container.Memory.Limit = "100m"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfContainerPortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.Port = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfContainerPortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.Port = 5000

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfHostsAreNotSet", func(t *testing.T) {

		params := validParams
		params.Hosts = []string{}

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfOneOrMoreHostsAreSet", func(t *testing.T) {

		params := validParams
		params.Hosts = []string{"gke.estafette.io"}

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAutoscaleMinReplicasIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.MinReplicas = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleMinReplicasIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.MinReplicas = 5

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAutoscaleMaxReplicasIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.MaxReplicas = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleMaxReplicasIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.MaxReplicas = 15

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAutoscaleCPUPercentageIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Autoscale.CPUPercentage = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAutoscaleCPUPercentageIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Autoscale.CPUPercentage = 35

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessPathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Path = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessPathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.Path = "/liveness"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessInitialDelaySecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.InitialDelaySeconds = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessInitialDelaySecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.InitialDelaySeconds = 30

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfLivenessTimeoutSecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.TimeoutSeconds = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfLivenessTimeoutSecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.LivenessProbe.TimeoutSeconds = 2

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfReadinessProbePathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Path = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfReadinessProbePathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.Path = "/readiness"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfReadinessProbeTimeoutSecondsIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.TimeoutSeconds = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfReadinessProbeTimeoutSecondsIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.ReadinessProbe.TimeoutSeconds = 2

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsPathIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Path = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsPathIsNotEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Path = "/metrics"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsPathIsEmptyButScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = "false"
		params.Container.Metrics.Path = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsPortIsZeroOrLess", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Port = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsPortIsLargerThanZero", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Port = 5000

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsPortIsZeroOrLessButScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = "false"
		params.Container.Metrics.Port = 0

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsScrapeIsEmpty", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfMetricsScrapeIsTrue", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = "true"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsTrueIfMetricsScrapeIsFalse", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = "false"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfMetricsScrapeIsNonBoolean", func(t *testing.T) {

		params := validParams
		params.Container.Metrics.Scrape = "yessir"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})
}
