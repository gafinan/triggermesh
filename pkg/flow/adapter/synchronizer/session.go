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
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type storage struct {
	sync.Mutex
	sessions map[string]channel
}

type channel struct {
	sync.Mutex
	c chan *cloudevents.Event
}

func newStorage() *storage {
	return &storage{
		sessions: make(map[string]channel),
	}
}

func (s *storage) add(id string) <-chan *cloudevents.Event {
	s.Lock()
	defer s.Unlock()

	respChan := make(chan *cloudevents.Event)
	s.sessions[id] = channel{
		c: respChan,
	}
	return respChan
}

func (s *storage) delete(id string) {
	s.Lock()
	defer s.Unlock()

	session, exists := s.sessions[id]
	if !exists {
		return
	}
	session.Lock()
	defer session.Unlock()

	delete(s.sessions, id)
	close(session.c)
}

func (s *storage) open(id string) (chan<- *cloudevents.Event, bool) {
	respChan, exists := s.sessions[id]
	if exists {
		respChan.Lock()
	}
	return respChan.c, exists
}

func (s *storage) close(id string) {
	respChan, exists := s.sessions[id]
	if exists {
		respChan.Unlock()
	}
}
