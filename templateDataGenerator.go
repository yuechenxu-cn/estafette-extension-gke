package main

import (
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func generateTemplateData(params Params, currentReplicas int, releaseID string) TemplateData {

	data := TemplateData{
		BuildVersion: params.BuildVersion,

		Name:             params.App,
		NameWithTrack:    params.App,
		Namespace:        params.Namespace,
		Schedule:         params.Schedule,
		Labels:           params.Labels,
		AppLabelSelector: params.App,

		Hosts:               params.Hosts,
		HostsJoined:         strings.Join(params.Hosts, ","),
		InternalHosts:       params.InternalHosts,
		InternalHostsJoined: strings.Join(params.InternalHosts, ","),
		IngressPath:         params.Basepath,
		InternalIngressPath: params.Basepath,

		IncludeReplicas: currentReplicas > 0,

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		Secrets:                 params.Secrets.Keys,
		MountApplicationSecrets: len(params.Secrets.Keys) > 0,
		SecretMountPath:         params.Secrets.MountPath,
		MountConfigmap:          len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0,
		ConfigMountPath:         params.Configs.MountPath,

		MountPayloadLogging:      params.EnablePayloadLogging,
		AddSafeToEvictAnnotation: params.EnablePayloadLogging,

		RollingUpdateMaxSurge:       params.RollingUpdate.MaxSurge,
		RollingUpdateMaxUnavailable: params.RollingUpdate.MaxUnavailable,

		PreferPreemptibles:        params.ChaosProof,
		MountServiceAccountSecret: params.UseGoogleCloudCredentials,

		Container: ContainerData{
			Repository: params.Container.ImageRepository,
			Name:       params.Container.ImageName,
			Tag:        params.Container.ImageTag,
			Port:       params.Container.Port,

			CPURequest:    params.Container.CPU.Request,
			CPULimit:      params.Container.CPU.Limit,
			MemoryRequest: params.Container.Memory.Request,
			MemoryLimit:   params.Container.Memory.Limit,

			EnvironmentVariables: params.Container.EnvironmentVariables,

			Liveness: ProbeData{
				Path:                params.Container.LivenessProbe.Path,
				Port:                params.Container.LivenessProbe.Port,
				InitialDelaySeconds: params.Container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.LivenessProbe.TimeoutSeconds,
				IncludeOnContainer:  true,
			},
			Readiness: ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				Port:                params.Container.ReadinessProbe.Port,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
				IncludeOnContainer:  params.Sidecar.Type != "openresty" || params.Container.ReadinessProbe.Port != params.Container.Port || params.Container.ReadinessProbe.Path != params.Sidecar.HealthCheckPath,
			},
			Metrics: MetricsData{
				Path: params.Container.Metrics.Path,
				Port: params.Container.Metrics.Port,
			},
		},
	}

	if params.UseGoogleCloudCredentials {
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "GOOGLE_APPLICATION_CREDENTIALS", "/gcp-service-account/service-account-key.json")
	}

	mainSidecar := buildSidecar(params.Sidecar, params.Request)
	data.Sidecars = append(data.Sidecars, mainSidecar)

	for _, sidecarParams := range params.Sidecars {
		sidecar := buildSidecar(sidecarParams, params.Request)
		data.Sidecars = append(data.Sidecars, sidecar)
	}

	// set request params on the nginx ingress
	data.NginxIngressProxyConnectTimeout = params.Request.Timeout
	data.NginxIngressProxySendTimeout = params.Request.Timeout
	data.NginxIngressProxyReadTimeout = params.Request.Timeout
	data.NginxIngressProxyBodySize = params.Request.MaxBodySize
	data.NginxIngressClientBodyBufferSize = params.Request.ClientBodyBufferSize
	data.NginxIngressProxyBufferSize = params.Request.ProxyBufferSize
	data.NginxIngressProxyBuffersNumber = strconv.Itoa(params.Request.ProxyBuffersNumber)

	if params.Container.Metrics.Scrape != nil {
		data.Container.Metrics.Scrape = *params.Container.Metrics.Scrape
	}
	if params.Container.Lifecycle.PrestopSleep != nil {
		data.Container.UseLifecyclePreStopSleepCommand = *params.Container.Lifecycle.PrestopSleep
	}
	if params.Container.Lifecycle.PrestopSleepSeconds != nil {
		data.Container.PreStopSleepSeconds = *params.Container.Lifecycle.PrestopSleepSeconds
	}

	if currentReplicas > 0 {
		data.Replicas = currentReplicas
	} else {
		data.Replicas = data.MinReplicas
	}

	if releaseID != "" {
		data.IncludeReleaseIDLabel = true
		data.ReleaseIDLabel = releaseID
	}

	switch params.Action {
	case "deploy-simple":
		data.IncludeTrackLabel = false
	case "deploy-canary":
		data.NameWithTrack += "-canary"
		data.IncludeTrackLabel = true
		data.TrackLabel = "canary"
	case "deploy-stable":
		data.NameWithTrack += "-stable"
		data.IncludeTrackLabel = true
		data.TrackLabel = "stable"
	}

	data.ConfigmapFiles = params.Configs.RenderedFileContent

	data.ManifestData = map[string]interface{}{}
	for k, v := range params.Manifests.Data {
		data.ManifestData[k] = v
	}

	switch params.Visibility {
	case "private":
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false

	case "iap":
		data.ServiceType = "NodePort"
		data.UseNginxIngress = false
		data.UseGCEIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false

	case "public-whitelist":
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = len(params.WhitelistedIPS) > 0
		data.NginxIngressWhitelist = strings.Join(params.WhitelistedIPS, ",")

	case "public":
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
		data.LimitTrustedIPRanges = true
		data.OverrideDefaultWhitelist = false
	}

	if !strings.HasSuffix(data.IngressPath, "/") && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "/"
	}
	if data.UseGCEIngress && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "*"
	}
	if !strings.HasSuffix(data.InternalIngressPath, "/") && !strings.HasSuffix(data.InternalIngressPath, "*") {
		data.InternalIngressPath += "/"
	}

	data.TrustedIPRanges = params.TrustedIPRanges

	data.AdditionalVolumeMounts = []VolumeMountData{}
	for _, vm := range params.VolumeMounts {
		yamlBytes, err := yaml.Marshal(vm.Volume)
		if err == nil {
			data.AdditionalVolumeMounts = append(data.AdditionalVolumeMounts, VolumeMountData{
				Name:       vm.Name,
				MountPath:  vm.MountPath,
				VolumeYAML: string(yamlBytes),
			})
		}
	}
	data.MountAdditionalVolumes = len(data.AdditionalVolumeMounts) > 0

	data.AdditionalContainerPorts = []AdditionalPortData{}
	data.AdditionalServicePorts = []AdditionalPortData{}
	for _, ap := range params.Container.AdditionalPorts {
		additionalPortData := AdditionalPortData{
			Name:     ap.Name,
			Port:     ap.Port,
			Protocol: ap.Protocol,
		}
		data.AdditionalContainerPorts = append(data.AdditionalContainerPorts, additionalPortData)

		includeAsServicePort := ap.Visibility == params.Visibility

		if includeAsServicePort {
			data.AdditionalServicePorts = append(data.AdditionalServicePorts, additionalPortData)
		}
	}

	return data
}

func buildSidecar(sidecar SidecarParams, request RequestParams) SidecarData {
	builtSidecar := SidecarData{
		Type:                     sidecar.Type,
		Image:                    sidecar.Image,
		HealthCheckPath:          sidecar.HealthCheckPath,
		DbInstanceConnectionName: sidecar.DbInstanceConnectionName,
		SQLProxyPort:             sidecar.SQLProxyPort,
		CPURequest:               sidecar.CPU.Request,
		CPULimit:                 sidecar.CPU.Limit,
		MemoryRequest:            sidecar.Memory.Request,
		MemoryLimit:              sidecar.Memory.Limit,

		EnvironmentVariables: sidecar.EnvironmentVariables,
	}

	if builtSidecar.Type == "openresty" {
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "SEND_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_HEADER_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_CONNECT_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_SEND_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_READ_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_MAX_BODY_SIZE", request.MaxBodySize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_BUFFER_SIZE", request.ClientBodyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFER_SIZE", request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_SIZE", request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_NUMBER", strconv.Itoa(request.ProxyBuffersNumber))
	}

	return builtSidecar
}

func addEnvironmentVariableIfNotSet(environmentVariables map[string]interface{}, name, value string) map[string]interface{} {

	if environmentVariables == nil {
		environmentVariables = map[string]interface{}{}
	}
	if _, ok := environmentVariables[name]; !ok {
		environmentVariables[name] = value
	}

	return environmentVariables
}
