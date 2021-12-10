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
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

const ceExtensionsKey = "extensions"

type key struct {
	extensions bool
	path       string
	indices    []int
}

func newKey(jsonPath string) (*key, error) {
	result := &key{
		indices: make([]int, 0, 2),
	}

	// parse key path
	parts := strings.Split(jsonPath, ".")
	switch len(parts) {
	case 1:
		result.path = parts[0]
	case 2:
		if parts[0] != ceExtensionsKey {
			return result, fmt.Errorf("only %q is supported as a key nested path, got %q", ceExtensionsKey, parts[0])
		}
		result.path = parts[1]
		result.extensions = true
	default:
		return result, fmt.Errorf("key path %q is not supported", jsonPath)
	}

	// parse key indices
	index := strings.Index(result.path, "[")
	if index > 0 {
		indicesSuffix := result.path[index+1 : len(result.path)-1]
		indices := strings.Split(indicesSuffix, ":")
		i1, err := strconv.Atoi(indices[0])
		if err != nil {
			return result, fmt.Errorf("failed to parse key index %q: %w", indices[0], err)
		}
		i2, err := strconv.Atoi(indices[1])
		if err != nil {
			return result, fmt.Errorf("failed to parse key index %q: %w", indices[1], err)
		}
		if i1 < 0 || i2 < 0 || i1 >= i2 {
			return result, fmt.Errorf("indices %d:%d are not valid", i1, i2)
		}
		result.indices = []int{i1, i2}
		result.path = result.path[:index]
	}

	if !result.extensions && result.path != "id" {
		return result, fmt.Errorf("correlation keys other than \"id\" are not unique, hence not supported: %q", result.path)
	}

	return result, nil
}

func (k *key) get(ceCtx event.EventContext) (string, error) {
	var value string

	if k.extensions {
		k, err := ceCtx.GetExtension(k.path)
		if err != nil {
			return "", nil
		}
		value = k.(string)
	} else {
		value = ceCtx.GetID()
	}

	if len(k.indices) == 2 {
		if k.indices[1] >= len(value) {
			return "", fmt.Errorf("key size is out of event value range: %d > %d", k.indices[1], len(value))
		}
		return value[k.indices[0]:k.indices[1]], nil
	}
	return value, nil
}

func (k *key) set(ceCtx *event.EventContext) (string, error) {
	id := uuid.NewString()
	if len(k.indices) == 2 {
		id = id[k.indices[0]:k.indices[1]]
	}
	if k.extensions {
		c := *ceCtx
		if err := c.SetExtension(k.path, id); err != nil {
			return "", err
		}
		ceCtx = &c
	}
	return id, nil
}
