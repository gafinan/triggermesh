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

package v1alpha1

import (
	"sort"
	"strings"
)

// Accepted event types
const (
	// EventTypeWildcard represents events of any type.
	EventTypeWildcard = "*"
)

// Returned event types
const (
	// EventTypeResponse is the event type of a generic response returned by a target.
	EventTypeResponse = "io.triggermesh.targets.response"
)

// String implements fmt.Stringer.
func (kv EnvKeyValue) String() string {
	keys := make([]string, 0, len(kv))

	for k := range kv {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var b strings.Builder

	for i, k := range keys {
		b.WriteString(k)
		b.WriteByte(':')
		b.WriteString(kv[k])

		if i+1 < len(keys) {
			b.WriteByte(',')
		}
	}

	return b.String()
}
