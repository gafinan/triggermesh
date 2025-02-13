//go:build !noclibs

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

package ibmmqtarget

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"
	"github.com/triggermesh/triggermesh/pkg/targets/adapter/ibmmqtarget/mq"
)

var _ pkgadapter.Adapter = (*ibmmqtargetAdapter)(nil)

type ibmmqtargetAdapter struct {
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
	mqEnvs   *TargetEnvAccessor
	reply    *mq.ReplyTo
	queue    *mq.Object
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*TargetEnvAccessor)
	logger := logging.FromContext(ctx)

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType(v1alpha1.IBMMQTargetGenericResponseEventType),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}
	return &ibmmqtargetAdapter{
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
		mqEnvs:   env,
		reply: &mq.ReplyTo{
			Manager: env.ReplyToManager,
			Queue:   env.ReplyToQueue,
		},
	}
}

// Returns if stopCh is closed or Send() returns an error.
func (a *ibmmqtargetAdapter) Start(ctx context.Context) error {
	a.logger.Info("Starting IBMMQTarget Adapter")

	conn, err := mq.NewConnection(a.mqEnvs.EnvConnectionConfig.ConnectionConfig())
	if err != nil {
		return fmt.Errorf("failed to create IBM MQ connection: %w", err)
	}
	defer conn.Disc()

	queue, err := mq.OpenQueue(a.mqEnvs.EnvConnectionConfig.QueueName, a.reply, conn)
	if err != nil {
		return fmt.Errorf("failed to open IBM MQ queue: %w", err)
	}
	defer queue.Close()

	a.queue = &queue

	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *ibmmqtargetAdapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	var msg []byte

	if a.mqEnvs.DiscardCEContext {
		msg = event.Data()
	} else {
		jsonEvent, err := json.Marshal(event)
		if err != nil {
			return a.replier.Error(&event, targetce.ErrorCodeRequestParsing, err, nil)
		}
		msg = jsonEvent
	}

	correlationID := event.ID()
	extensions := event.Extensions()
	if extensions != nil {
		if cid, ok := extensions[mq.CECorrelIDAttr]; ok {
			correlationID = cid.(string)
		}
	}

	if err := a.queue.Put(msg, correlationID); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, nil)
	}

	return a.replier.Ok(&event, "ok")
}
