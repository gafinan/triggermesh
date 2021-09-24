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

package common

// Common environment variables propagated to adapters.
const (
	EnvName      = "NAME"
	EnvNamespace = "NAMESPACE"

	envSink                  = "K_SINK"
	envComponent             = "K_COMPONENT"
	envMetricsPrometheusPort = "METRICS_PROMETHEUS_PORT"

	// Common AWS attributes
	EnvARN             = "ARN"
	EnvAccessKeyID     = "AWS_ACCESS_KEY_ID"
	EnvSecretAccessKey = "AWS_SECRET_ACCESS_KEY" //nolint:gosec

	// Common Azure attributes
	EnvTenantID     = "AZURE_TENANT_ID"
	EnvClientID     = "AZURE_CLIENT_ID"
	EnvClientSecret = "AZURE_CLIENT_SECRET"
	EnvConnStr      = "AZURE_CONN_STR"

	// Azure EventHub auth attributes
	EnvHubName      = "AZURE_HUB_NAME"
	EnvHubNamespace = "AZURE_HUB_NAMESPACE"
	EnvHubKeyName   = "EVENTHUB_KEY_NAME"
	EnvHubKeyValue  = "EVENTHUB_KEY_VALUE"

	// Google Cloud
	EnvGCloudProject = "GCLOUD_PROJECT"
	EnvGCloudSAKey   = "GCLOUD_SERVICEACCOUNT_KEY"
)