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
)

const ceExtensionsKey = "extensions"

type ceContextParser func(event.EventContext) (string, error)

func newKeyGetter(jsonPath string) (ceContextParser, error) {
	var extensions bool
	var key string
	var err error

	// parse key path
	parts := strings.Split(jsonPath, ".")
	switch len(parts) {
	case 1:
		key = parts[0]
	case 2:
		if parts[0] != ceExtensionsKey {
			return nil, fmt.Errorf("expected path is %q, got %q", ceExtensionsKey, parts[0])
		}
		key = parts[1]
		extensions = true
	default:
		return nil, fmt.Errorf("key path %q is not supported", jsonPath)
	}

	// parse key indices
	i1 := 0
	i2 := len(key) - 1
	index := strings.Index(key, "[")
	if index > 0 {
		section := key[index+1 : len(key)-1]
		indices := strings.Split(section, ":")
		if i1, err = strconv.Atoi(indices[0]); err != nil {
			return nil, fmt.Errorf("failed to parse key index %q: %w", indices[0], err)
		}
		if i2, err = strconv.Atoi(indices[1]); err != nil {
			return nil, fmt.Errorf("failed to parse key index %q: %w", indices[1], err)
		}
		if i1 < 0 || i2 < 0 || i1 >= i2 {
			return nil, fmt.Errorf("indices %d:%d are not valid", i1, i2)
		}
		key = key[:index]
	}

	return func(ceCtx event.EventContext) (string, error) {
		var result string

		if extensions {
			k, err := ceCtx.GetExtension(key)
			if err != nil {
				return "", err
			}
			result = k.(string)
		} else {
			if key != "id" {
				return "", fmt.Errorf("expected \"id\" as a unique key, got %q", key)
			}
			result = ceCtx.GetID()
		}

		if index > 0 {
			if i2 >= len(result) {
				return "", fmt.Errorf("selected key section is out of event key range")
			}
			return result[i1:i2], nil
		}

		return result, nil
	}, nil
}
