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

package synchronizer

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"
)

var _ pkgadapter.Adapter = (*adapter)(nil)

type adapter struct {
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger

	requestKey      ceContextParser
	responseKey     ceContextParser
	responseTimeout int

	sessions *storage
	sinkURL  string
}

// NewTarget adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	requestParser, err := newKeyGetter(env.RequestKey)
	if err != nil {
		panic(err)
	}
	responseParser, err := newKeyGetter(env.ResponseCorrelationKey)
	if err != nil {
		panic(err)
	}

	return &adapter{
		ceClient:        ceClient,
		logger:          logger,
		requestKey:      requestParser,
		responseKey:     responseParser,
		responseTimeout: env.ResponseWaitTimeout,
		sessions:        newStorage(),
		sinkURL:         env.Sink,
	}
}

// Returns if stopCh is closed or Send() returns an error.
func (a *adapter) Start(ctx context.Context) error {
	a.logger.Info("Starting Synchronizer Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *adapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	correlationID, err := a.responseKey(event.Context)
	if err != nil {
		// That's fine, it may be a request
		a.logger.Warnf("Event correlation ID: %v", err)
	}

	eventID, err := a.requestKey(event.Context)
	if err != nil {
		// Something doesn't work, fail
		a.logger.Errorf("Event request ID: %v", err)
		return nil, err
	}

	switch {
	case correlationID != "":
		return a.handleResponse(ctx, correlationID, event)
	default:
		return a.handleRequest(ctx, eventID, event)
	}
}

func (a *adapter) handleRequest(ctx context.Context, id string, event cloudevents.Event) (*cloudevents.Event, error) {
	if err := a.ceClient.Send(cloudevents.ContextWithTarget(ctx, a.sinkURL), event); err != nil {
		return nil, err
	}

	resultChan := a.sessions.add(id)
	t := time.NewTimer(time.Duration(a.responseTimeout) * time.Second)
	defer t.Stop()

	select {
	case result := <-resultChan:
		return result, nil
	case <-t.C:
		return nil, fmt.Errorf("timeout")
	}
}

func (a *adapter) handleResponse(ctx context.Context, correlationID string, event cloudevents.Event) (*cloudevents.Event, error) {
	responseChan, exists := a.sessions.get(correlationID)
	if !exists {
		return nil, fmt.Errorf("session expired")
	}

	responseChan <- &event
	a.sessions.delete(correlationID)

	return nil, nil
}
