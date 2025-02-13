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

import (
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/system"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	"github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	"github.com/triggermesh/triggermesh/pkg/sources/reconciler/common/resource"
)

const metricsPrometheusPort uint16 = 9092

// ComponentName returns the component name for the given source object.
func ComponentName(src kmeta.OwnerRefable) string {
	return strings.ToLower(src.GetGroupVersionKind().Kind)
}

// MTAdapterObjectName returns a unique name to apply to all objects related to
// the given source's multi-tenant adapter (RBAC, Deployment/KnService, ...).
func MTAdapterObjectName(src kmeta.OwnerRefable) string {
	return ComponentName(src) + "-" + componentAdapter
}

// ServiceAccountName returns the name to set on the ServiceAccount associated
// with the given source.
func ServiceAccountName(src v1alpha1.EventSource) string {
	if v1alpha1.WantsOwnServiceAccount(src) {
		srcName := src.GetName()

		// Edge case: we need to make sure some characters are inserted
		// between the component name and the source name to avoid
		// clashing with the shared "{kind}-adapter" ServiceAccount in
		// case the source is named "adapter". We picked 'i' for
		// "instance" to keep it short yet distinguishable.
		return kmeta.ChildName(ComponentName(src)+"-i-", srcName)
	}

	return MTAdapterObjectName(src)
}

// NewAdapterDeployment is a wrapper around resource.NewDeployment which
// pre-populates attributes common to all adapters backed by a Deployment.
func NewAdapterDeployment(src v1alpha1.EventSource, sinkURI *apis.URL, opts ...resource.ObjectOption) *appsv1.Deployment {
	srcNs := src.GetNamespace()
	srcName := src.GetName()

	var sinkURIStr string
	if sinkURI != nil {
		sinkURIStr = sinkURI.String()
	}

	return resource.NewDeployment(srcNs, kmeta.ChildName(ComponentName(src)+"-", srcName),
		append(commonAdapterDeploymentOptions(src), append([]resource.ObjectOption{
			resource.Controller(src),

			resource.Label(appInstanceLabel, srcName),
			resource.Selector(appInstanceLabel, srcName),

			resource.EnvVar(envSink, sinkURIStr),
		}, opts...)...)...,
	)
}

// NewMTAdapterDeployment is a wrapper around resource.NewDeployment which
// pre-populates attributes common to all multi-tenant adapters backed by a
// Deployment.
func NewMTAdapterDeployment(src v1alpha1.EventSource, opts ...resource.ObjectOption) *appsv1.Deployment {
	srcNs := src.GetNamespace()

	return resource.NewDeployment(srcNs, MTAdapterObjectName(src),
		append(commonAdapterDeploymentOptions(src), append([]resource.ObjectOption{
			resource.EnvVar(EnvNamespace, srcNs),
			resource.EnvVar(system.NamespaceEnvKey, srcNs), // required to enable HA
		}, opts...)...)...,
	)
}

// commonAdapterDeploymentOptions returns a set of ObjectOptions common to all
// adapters backed by a Deployment.
func commonAdapterDeploymentOptions(src v1alpha1.EventSource) []resource.ObjectOption {
	app := ComponentName(src)

	return []resource.ObjectOption{
		resource.TerminationErrorToLogs,

		resource.Label(appNameLabel, app),
		resource.Label(appComponentLabel, componentAdapter),
		resource.Label(appPartOfLabel, partOf),
		resource.Label(appManagedByLabel, managedBy),

		resource.Selector(appNameLabel, app),
		resource.Selector(appComponentLabel, componentAdapter),
		resource.PodLabel(appPartOfLabel, partOf),
		resource.PodLabel(appManagedByLabel, managedBy),

		resource.ServiceAccount(ServiceAccountName(src)),

		resource.EnvVar(envComponent, app),
	}
}

// NewAdapterKnService is a wrapper around resource.NewKnService which
// pre-populates attributes common to all adapters backed by a Knative Service.
func NewAdapterKnService(src v1alpha1.EventSource, sinkURI *apis.URL, opts ...resource.ObjectOption) *servingv1.Service {
	srcNs := src.GetNamespace()
	srcName := src.GetName()

	var sinkURIStr string
	if sinkURI != nil {
		sinkURIStr = sinkURI.String()
	}

	return resource.NewKnService(srcNs, kmeta.ChildName(ComponentName(src)+"-", srcName),
		append(commonAdapterKnServiceOptions(src), append([]resource.ObjectOption{
			resource.Controller(src),

			resource.Label(appInstanceLabel, srcName),
			resource.PodLabel(appInstanceLabel, srcName),

			resource.EnvVar(envSink, sinkURIStr),
		}, opts...)...)...,
	)
}

// NewMTAdapterKnService is a wrapper around resource.NewKnService which
// pre-populates attributes common to all multi-tenant adapters backed by a
// Knative Service.
func NewMTAdapterKnService(src v1alpha1.EventSource, opts ...resource.ObjectOption) *servingv1.Service {
	srcNs := src.GetNamespace()

	return resource.NewKnService(srcNs, MTAdapterObjectName(src),
		append(commonAdapterKnServiceOptions(src), append([]resource.ObjectOption{
			resource.EnvVar(EnvNamespace, srcNs),
			resource.EnvVar(system.NamespaceEnvKey, srcNs), // required to enable HA
		}, opts...)...)...,
	)
}

// commonAdapterKnServiceOptions returns a set of ObjectOptions common to all
// adapters backed by a Knative Service.
func commonAdapterKnServiceOptions(src v1alpha1.EventSource) []resource.ObjectOption {
	app := ComponentName(src)

	return []resource.ObjectOption{
		resource.Label(appNameLabel, app),
		resource.Label(appComponentLabel, componentAdapter),
		resource.Label(appPartOfLabel, partOf),
		resource.Label(appManagedByLabel, managedBy),

		resource.PodLabel(appNameLabel, app),
		resource.PodLabel(appComponentLabel, componentAdapter),
		resource.PodLabel(appPartOfLabel, partOf),
		resource.PodLabel(appManagedByLabel, managedBy),

		resource.ServiceAccount(MTAdapterObjectName(src)),

		resource.EnvVar(envComponent, app),
		resource.EnvVar(envMetricsPrometheusPort, strconv.FormatUint(uint64(metricsPrometheusPort), 10)),
	}
}

// newServiceAccount returns a ServiceAccount object with its OwnerReferences
// metadata attribute populated from the given owners.
func newServiceAccount(src v1alpha1.EventSource, owners []kmeta.OwnerRefable) *corev1.ServiceAccount {
	return resource.NewServiceAccount(src.GetNamespace(), ServiceAccountName(src),
		resource.Owners(owners...),
		resource.Labels(CommonObjectLabels(src)),
	)
}

// newRoleBinding returns a RoleBinding object that binds a ServiceAccount
// (namespace-scoped) to a ClusterRole (cluster-scoped).
func newRoleBinding(src v1alpha1.EventSource, owner *corev1.ServiceAccount) *rbacv1.RoleBinding {
	crGVK := rbacv1.SchemeGroupVersion.WithKind("ClusterRole")
	saGVK := corev1.SchemeGroupVersion.WithKind("ServiceAccount")

	ns := src.GetNamespace()
	n := MTAdapterObjectName(src)

	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      n,
			Labels:    CommonObjectLabels(src),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: crGVK.Group,
			Kind:     crGVK.Kind,
			Name:     n,
		},
		Subjects: []rbacv1.Subject{{
			APIGroup:  saGVK.Group,
			Kind:      saGVK.Kind,
			Namespace: ns,
			Name:      n,
		}},
	}

	OwnByServiceAccount(rb, owner)

	return rb
}

// OwnByServiceAccount sets the owner of obj to the given ServiceAccount.
func OwnByServiceAccount(obj metav1.Object, owner *corev1.ServiceAccount) {
	saGVK := corev1.SchemeGroupVersion.WithKind("ServiceAccount")

	obj.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(owner, saGVK),
	})
}

// CommonObjectLabels returns a set of labels which are always applied to
// objects reconciled for the given source type.
func CommonObjectLabels(src kmeta.OwnerRefable) labels.Set {
	return labels.Set{
		appNameLabel:      ComponentName(src),
		appComponentLabel: componentAdapter,
		appPartOfLabel:    partOf,
		appManagedByLabel: managedBy,
	}
}

// MaybeAppendValueFromEnvVar conditionally appends an EnvVar to env based on
// the contents of valueFrom.
// ValueFromSecret takes precedence over Value in case the API didn't reject
// the object despite the CRD's schema validation
func MaybeAppendValueFromEnvVar(envs []corev1.EnvVar, key string, valueFrom v1alpha1.ValueFromField) []corev1.EnvVar {
	if vfs := valueFrom.ValueFromSecret; vfs != nil {
		return append(envs, corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: vfs,
			},
		})
	}

	if v := valueFrom.Value; v != "" {
		return append(envs, corev1.EnvVar{
			Name:  key,
			Value: v,
		})
	}

	return envs
}

// MakeAWSAuthEnvVars returns environment variables for the given AWS
// authentication method.
func MakeAWSAuthEnvVars(auth v1alpha1.AWSAuth) []corev1.EnvVar {
	var authEnvVars []corev1.EnvVar

	if creds := auth.Credentials; creds != nil {
		authEnvVars = MaybeAppendValueFromEnvVar(authEnvVars, EnvAccessKeyID, creds.AccessKeyID)
		authEnvVars = MaybeAppendValueFromEnvVar(authEnvVars, EnvSecretAccessKey, creds.SecretAccessKey)
	}

	return authEnvVars
}

// MakeAWSEndpointEnvVars returns environment variables for the given AWS
// endpoint parameters.
func MakeAWSEndpointEnvVars(endpoint *v1alpha1.AWSEndpoint) []corev1.EnvVar {
	if endpoint == nil {
		return nil
	}

	var endpointEnvVars []corev1.EnvVar

	if url := endpoint.URL; url != nil {
		endpointEnvVars = append(endpointEnvVars, corev1.EnvVar{
			Name:  EnvEndpointURL,
			Value: url.String(),
		})
	}

	return endpointEnvVars
}
