/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logzmetricstarget

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/kmeta"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	"github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	pkgreconciler "github.com/triggermesh/triggermesh/pkg/targets/reconciler"
	"github.com/triggermesh/triggermesh/pkg/targets/reconciler/resources"
)

const (
	adapterName = "logzmetricstarget"

	envCortexEndpoint           = "OPENTELEMETRY_CORTEX_ENDPOINT"
	envCortexBearerToken        = "OPENTELEMETRY_CORTEX_BEARER_TOKEN"
	envOpenTelemetryInstruments = "OPENTELEMETRY_INSTRUMENTS"

	envEventsPayloadPolicy = "EVENTS_PAYLOAD_POLICY"
)

// adapterConfig contains properties used to configure the target's adapter.
// Public fields are automatically populated by envconfig.
type adapterConfig struct {
	// Configuration accessor for logging/metrics/tracing
	configs source.ConfigAccessor
	// Container image
	Image string `envconfig:"OPENTELEMETRYTARGET_IMAGE" default:"gcr.io/triggermesh/opentelemetrytarget-adapter"`
}

// makeTargetAdapterKService generates (but does not insert into K8s) the Target Adapter KService.
func makeTargetAdapterKService(o *v1alpha1.LogzMetricsTarget, cfg *adapterConfig) (*servingv1.Service, error) {
	envApp, err := makeAppEnv(o)
	if err != nil {
		return nil, err
	}

	ksvcLabels := pkgreconciler.MakeAdapterLabels(adapterName, o.Name)
	podLabels := pkgreconciler.MakeAdapterLabels(adapterName, o.Name)
	name := kmeta.ChildName(adapterName+"-", o.Name)
	envSvc := pkgreconciler.MakeServiceEnv(o.Name, o.Namespace)
	envObs := pkgreconciler.MakeObsEnv(cfg.configs)
	envs := append(envSvc, append(envApp, envObs...)...)

	return resources.MakeKService(o.Namespace, name, cfg.Image,
		resources.KsvcLabels(ksvcLabels),
		resources.KsvcLabelVisibilityClusterLocal,
		resources.KsvcPodLabels(podLabels),
		resources.KsvcOwner(o),
		resources.KsvcPodEnvVars(envs)), nil
}

func makeAppEnv(o *v1alpha1.LogzMetricsTarget) ([]corev1.EnvVar, error) {
	instruments, err := json.Marshal(o.Spec.Instruments)
	if err != nil {
		return nil, err
	}

	env := []corev1.EnvVar{
		{
			Name: envCortexBearerToken,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: o.Spec.Connection.Token.SecretKeyRef,
			},
		}, {
			Name:  envCortexEndpoint,
			Value: o.Spec.Connection.ListenerURL,
		}, {
			Name:  envOpenTelemetryInstruments,
			Value: string(instruments),
		}, {
			Name:  pkgreconciler.EnvBridgeID,
			Value: pkgreconciler.GetStatefulBridgeID(o),
		},
	}

	if o.Spec.EventOptions != nil && o.Spec.EventOptions.PayloadPolicy != nil {
		env = append(env, corev1.EnvVar{
			Name:  envEventsPayloadPolicy,
			Value: string(*o.Spec.EventOptions.PayloadPolicy),
		})
	}

	return env, nil
}
