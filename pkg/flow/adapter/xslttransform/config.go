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
	"errors"

	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	pkgadapter.EnvConfig

	XSLT string `envconfig:"XSLTTRANSFORM_XSLT"`
	// If set to true, enables consuming structured CloudEvents that include
	// fields for the XML and XSLT field.
	AllowXSLTOverride bool `envconfig:"XSLTTRANSFORM_ALLOW_XSLT_OVERRIDE" required:"true"`

	// BridgeIdentifier is the name of the bridge workflow this target is part of
	BridgeIdentifier string `envconfig:"EVENTS_BRIDGE_IDENTIFIER"`
}

func (e *envAccessor) validate() error {
	if !e.AllowXSLTOverride && e.XSLT == "" {
		return errors.New("if XSLT cannot be overriden by CloudEvent payloads, configured XSLT cannot be empty")
	}
	return nil
}
