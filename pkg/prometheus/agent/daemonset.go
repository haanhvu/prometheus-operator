// Copyright 2023 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusagent

import (
	"net/url"
	"path"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus-operator/prometheus-operator/pkg/operator"
	prompkg "github.com/prometheus-operator/prometheus-operator/pkg/prometheus"
	"github.com/prometheus-operator/prometheus-operator/pkg/webconfig"
)

func makeDaemonSetSpec(
	p monitoringv1.PrometheusInterface,
	c *prompkg.Config,
	cg *prompkg.ConfigGenerator,
	tlsSecrets *operator.ShardedSecret,
) (*appsv1.DaemonSetSpec, error) {
	cpf := p.GetCommonPrometheusFields()

	pImagePath, err := operator.BuildImagePath(
		operator.StringPtrValOrDefault(cpf.Image, ""),
		operator.StringValOrDefault("", c.PrometheusDefaultBaseImage),
		operator.StringValOrDefault(cpf.Version, operator.DefaultPrometheusVersion),
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	cpf.EnableFeatures = append(cpf.EnableFeatures, "agent")
	promArgs := prompkg.BuildCommonPrometheusArgs(cpf, cg)
	promArgs = appendAgentArgs(promArgs, cg, cpf.WALCompression)

	var ports []v1.ContainerPort
	if !cpf.ListenLocal {
		ports = []v1.ContainerPort{
			{
				Name:          cpf.PortName,
				ContainerPort: 9090,
				Protocol:      v1.ProtocolTCP,
			},
		}
	}

	volumes, promVolumeMounts, err := prompkg.BuildCommonVolumes(p, tlsSecrets)
	if err != nil {
		return nil, err
	}

	configReloaderVolumeMounts := []v1.VolumeMount{
		{
			Name:      "config",
			MountPath: prompkg.ConfDir,
		},
		{
			Name:      "config-out",
			MountPath: prompkg.ConfOutDir,
		},
	}

	var configReloaderWebConfigFile string

	// Mount web config and web TLS credentials as volumes.
	// We always mount the web config file for versions greater than 2.24.0.
	// With this we avoid redeploying prometheus when reconfiguring between
	// HTTP and HTTPS and vice-versa.
	webConfigGenerator := cg.WithMinimumVersion("2.24.0")
	if webConfigGenerator.IsCompatible() {
		var fields monitoringv1.WebConfigFileFields
		if cpf.Web != nil {
			fields = cpf.Web.WebConfigFileFields
		}

		webConfig, err := webconfig.New(prompkg.WebConfigDir, prompkg.WebConfigSecretName(p), fields)
		if err != nil {
			return nil, err
		}

		confArg, configVol, configMount, err := webConfig.GetMountParameters()
		if err != nil {
			return nil, err
		}
		promArgs = append(promArgs, confArg)
		volumes = append(volumes, configVol...)
		promVolumeMounts = append(promVolumeMounts, configMount...)

		// To avoid breaking users deploying an old version of the config-reloader image.
		// TODO: remove the if condition after v0.72.0.
		if cpf.Web != nil {
			configReloaderWebConfigFile = confArg.Value
			configReloaderVolumeMounts = append(configReloaderVolumeMounts, configMount...)
		}
	} else if cpf.Web != nil {
		webConfigGenerator.Warn("web.config.file")
	}

	// The /-/ready handler returns OK only after the TSDB initialization has
	// completed. The WAL replay can take a significant time for large setups
	// hence we enable the startup probe with a generous failure threshold (15
	// minutes) to ensure that the readiness probe only comes into effect once
	// Prometheus is effectively ready.
	// We don't want to use the /-/healthy handler here because it returns OK as
	// soon as the web server is started (irrespective of the WAL replay).
	readyProbeHandler := prompkg.ProbeHandler("/-/ready", cpf, webConfigGenerator)
	startupPeriodSeconds, startupFailureThreshold := prompkg.GetStatupProbePeriodSecondsAndFailureThreshold(cpf)
	startupProbe := &v1.Probe{
		ProbeHandler:     readyProbeHandler,
		TimeoutSeconds:   prompkg.ProbeTimeoutSeconds,
		PeriodSeconds:    startupPeriodSeconds,
		FailureThreshold: startupFailureThreshold,
	}

	readinessProbe := &v1.Probe{
		ProbeHandler:     readyProbeHandler,
		TimeoutSeconds:   prompkg.ProbeTimeoutSeconds,
		PeriodSeconds:    5,
		FailureThreshold: 3,
	}

	livenessProbe := &v1.Probe{
		ProbeHandler:     prompkg.ProbeHandler("/-/healthy", cpf, webConfigGenerator),
		TimeoutSeconds:   prompkg.ProbeTimeoutSeconds,
		PeriodSeconds:    5,
		FailureThreshold: 6,
	}

	podAnnotations, podLabels := prompkg.BuildPodMetadata(cpf, cg)
	// In cases where an existing selector label is modified, or a new one is added, new sts cannot match existing pods.
	// We should try to avoid removing such immutable fields whenever possible since doing
	// so forces us to enter the 'recreate cycle' and can potentially lead to downtime.
	// The requirement to make a change here should be carefully evaluated.
	podSelectorLabels := makeSelectorLabels(p.GetObjectMeta().GetName())

	for k, v := range podSelectorLabels {
		podLabels[k] = v
	}

	finalSelectorLabels := c.Labels.Merge(podSelectorLabels)
	finalLabels := c.Labels.Merge(podLabels)

	var additionalContainers, operatorInitContainers []v1.Container

	var watchedDirectories []string

	var minReadySeconds int32
	if cpf.MinReadySeconds != nil {
		minReadySeconds = int32(*cpf.MinReadySeconds)
	}
}

func BuildConfigReloaderForDaemonSet(
	p monitoringv1.PrometheusInterface,
	c *prompkg.Config,
	initContainer bool,
	mounts []v1.VolumeMount,
	watchedDirectories []string,
	opts ...operator.ReloaderOption,
) v1.Container {
	cpf := p.GetCommonPrometheusFields()

	reloaderOptions := []operator.ReloaderOption{
		operator.ReloaderConfig(c.ReloaderConfig),
		operator.LogFormat(cpf.LogFormat),
		operator.LogLevel(cpf.LogLevel),
		operator.VolumeMounts(mounts),
		operator.ConfigFile(path.Join(ConfDir, ConfigFilename)),
		operator.ConfigEnvsubstFile(path.Join(ConfOutDir, ConfigEnvsubstFilename)),
		operator.WatchedDirectories(watchedDirectories),
		operator.ImagePullPolicy(cpf.ImagePullPolicy),
	}
	reloaderOptions = append(reloaderOptions, opts...)

	name := "config-reloader"
	if initContainer {
		name = "init-config-reloader"
		reloaderOptions = append(reloaderOptions, operator.ReloaderRunOnce())
		return operator.CreateConfigReloader(name, reloaderOptions...)
	}

	if ptr.Deref(cpf.ReloadStrategy, monitoringv1.HTTPReloadStrategyType) == monitoringv1.ProcessSignalReloadStrategyType {
		reloaderOptions = append(reloaderOptions,
			operator.ReloaderUseSignal(),
		)
		reloaderOptions = append(reloaderOptions,
			operator.RuntimeInfoURL(url.URL{
				Scheme: cpf.PrometheusURIScheme(),
				Host:   c.LocalHost + ":9090",
				Path:   path.Clean(cpf.WebRoutePrefix() + "/api/v1/status/runtimeinfo"),
			}),
		)
	} else {
		reloaderOptions = append(reloaderOptions,
			operator.ListenLocal(cpf.ListenLocal),
			operator.LocalHost(c.LocalHost),
		)
		reloaderOptions = append(reloaderOptions,
			operator.ReloaderURL(url.URL{
				Scheme: cpf.PrometheusURIScheme(),
				Host:   c.LocalHost + ":9090",
				Path:   path.Clean(cpf.WebRoutePrefix() + "/-/reload"),
			}),
		)
	}

	return operator.CreateConfigReloaderForDaemonSet(name, reloaderOptions...)
}
