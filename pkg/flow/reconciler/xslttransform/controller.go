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

package xslttransform

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"k8s.io/client-go/tools/cache"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	servingclient "knative.dev/serving/pkg/client/injection/client"
	serviceinformerv1 "knative.dev/serving/pkg/client/injection/informers/serving/v1/service"

	"github.com/triggermesh/triggermesh/pkg/apis/flow/v1alpha1"
	informerv1alpha1 "github.com/triggermesh/triggermesh/pkg/client/generated/injection/informers/flow/v1alpha1/xslttransform"
	reconcilerv1alpha1 "github.com/triggermesh/triggermesh/pkg/client/generated/injection/reconciler/flow/v1alpha1/xslttransform"
	libreconciler "github.com/triggermesh/triggermesh/pkg/flow/reconciler"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	adapterCfg := &adapterConfig{
		configs: source.WatchConfigurations(ctx, adapterName, cmw, source.WithLogging, source.WithMetrics),
	}
	envconfig.MustProcess(adapterName, adapterCfg)

	componentInformer := informerv1alpha1.Get(ctx)
	serviceInformer := serviceinformerv1.Get(ctx)

	r := &reconciler{
		ksvcr:      libreconciler.NewKServiceReconciler(servingclient.Get(ctx), serviceInformer.Lister()),
		adapterCfg: adapterCfg,
	}

	impl := reconcilerv1alpha1.NewImpl(ctx, r)

	componentInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	serviceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGVK((&v1alpha1.XSLTTransform{}).GetGroupVersionKind()),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
