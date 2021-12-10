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
	"net/http"
	"time"

	"go.uber.org/zap"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"
)

var _ pkgadapter.Adapter = (*adapter)(nil)

type correlationKey interface {
	get(event.EventContext) (string, error)
	set(*event.EventContext) (string, error)
}

type adapter struct {
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger

	requestKey      correlationKey
	responseKey     correlationKey
	responseTimeout int

	sessions *storage
	sinkURL  string
}

// NewTarget adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	requestCorrelationKey, err := newKey(env.RequestKey)
	if err != nil {
		panic(err)
	}
	responseCorrelationKey, err := newKey(env.ResponseCorrelationKey)
	if err != nil {
		panic(err)
	}

	return &adapter{
		ceClient:    ceClient,
		logger:      logger,
		requestKey:  requestCorrelationKey,
		responseKey: responseCorrelationKey,

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
	a.logger.Debugf("Received the event: %s", event.String())

	correlationID, err := a.responseKey.get(event.Context)
	if err != nil {
		a.logger.Debugf("Correlation key is not set: %v", err)
	} else if correlationID != "" {
		return a.handleResponse(ctx, correlationID, event)
	}

	eventID, err := a.requestKey.get(event.Context)
	if err != nil {
		a.logger.Errorf("Cannot get request key: %v", err)
		return nil, err
	}
	if eventID == "" {
		eventID, err = a.requestKey.set(&event.Context)
		if err != nil {
			a.logger.Errorf("Cannot set request key: %v", err)
			return nil, err
		}
	}
	return a.handleRequest(ctx, eventID, event)
}

func (a *adapter) handleRequest(ctx context.Context, id string, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	a.logger.Debugf("Handling request %q", id)

	respChan := a.sessions.add(id)

	sendErr := make(chan error)
	defer close(sendErr)

	go func(errChan chan error) {
		if res := a.ceClient.Send(cloudevents.ContextWithTarget(ctx, a.sinkURL), event); cloudevents.IsUndelivered(res) {
			errChan <- res
		}
	}(sendErr)

	a.logger.Debugf("Request forwarded to %q", a.sinkURL)

	t := time.NewTimer(time.Duration(a.responseTimeout) * time.Second)
	defer t.Stop()

	a.logger.Debugf("Waiting response for %q", id)

	select {
	case err := <-sendErr:
		a.logger.Errorf("Unable to forward the request: %v", err)
		a.sessions.delete(id)
		return nil, cloudevents.NewHTTPResult(http.StatusBadGateway, "unable to forward the request: %v", err)
	case result := <-respChan:
		if result == nil {
			a.logger.Errorf("No response from %q", id)
			return nil, cloudevents.NewHTTPResult(http.StatusInternalServerError, "failed to communicate the response")
		}
		a.logger.Debugf("Received response for %q", id)
		return result, cloudevents.ResultACK
	case <-t.C:
		a.logger.Errorf("Request %q timed out", id)
		return nil, cloudevents.NewHTTPResult(http.StatusGatewayTimeout, "backend did not respond in time")
	}
}

func (a *adapter) handleResponse(ctx context.Context, correlationID string, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	a.logger.Debugf("Handling response %q", correlationID)

	responseChan, exists := a.sessions.open(correlationID)
	if !exists {
		a.logger.Errorf("Session for %q does not exist", correlationID)
		return nil, cloudevents.NewHTTPResult(http.StatusInternalServerError, "client session does not exist")
	}
	defer a.sessions.close(correlationID)

	a.logger.Debugf("Forwarding response %q", correlationID)

	select {
	case responseChan <- &event:
		a.sessions.delete(correlationID)
		a.logger.Debugf("Response %q complete", correlationID)
		return nil, cloudevents.ResultACK
	default:
		a.logger.Errorf("Unable to forward the response %q", correlationID)
		return nil, cloudevents.NewHTTPResult(http.StatusBadGateway, "response channel closed")
	}
}
