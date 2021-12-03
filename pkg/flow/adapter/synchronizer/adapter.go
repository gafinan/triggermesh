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
		a.logger.Warnf("Event correlation ID: %v", err)
	} else {
		return a.handleResponse(ctx, correlationID, event)
	}

	eventID, err := a.requestKey(event.Context)
	if err != nil {
		a.logger.Errorf("Cannot parse request ID: %v", err)
		return nil, err
	} else {
		return a.handleRequest(ctx, eventID, event)
	}
}

func (a *adapter) handleRequest(ctx context.Context, id string, event cloudevents.Event) (*cloudevents.Event, error) {
	a.logger.Infof("Handling request %q", id)

	respChan := a.sessions.add(id)

	go func() {
		if res := a.ceClient.Send(cloudevents.ContextWithTarget(ctx, a.sinkURL), event); cloudevents.IsUndelivered(res) {
			a.logger.Errorf("Unable to forward the request: %v", res)
			// a.sessions.delete(id)
		}
	}()

	a.logger.Infof("Request forwarded to %q", a.sinkURL)

	t := time.NewTimer(time.Duration(a.responseTimeout) * time.Second)
	defer t.Stop()

	a.logger.Infof("Waiting response for %q", id)

	select {
	case result := <-respChan:
		if result == nil {
			a.logger.Infof("Nil response for %q", id)
			return nil, nil
		}
		a.logger.Infof("Received response for %q: %+v", id, result)
		return result, cloudevents.ResultACK
	case <-t.C:
		a.logger.Errorf("Request %q timed out", id)
		return nil, cloudevents.ResultNACK
	}
}

func (a *adapter) handleResponse(ctx context.Context, correlationID string, event cloudevents.Event) (*cloudevents.Event, error) {
	a.logger.Infof("Handling response %q", correlationID)

	responseChan, exists := a.sessions.open(correlationID)
	if !exists {
		a.logger.Errorf("Session for %q does not exist", correlationID)
		return nil, cloudevents.ResultNACK
	}
	defer a.sessions.close(correlationID)

	a.logger.Infof("Forwarding response %q", correlationID)

	select {
	case responseChan <- &event:
		a.sessions.delete(correlationID)
		a.logger.Infof("Response %q complete", correlationID)
		return nil, cloudevents.ResultACK
	default:
		a.logger.Errorf("Unable to forward the response %q", correlationID)
		return nil, cloudevents.ResultNACK
	}
}
