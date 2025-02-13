/*
Copyright 2020 TriggerMesh Inc.

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

package main

import (
	"github.com/golang-jwt/jwt/v4"

	"knative.dev/eventing/pkg/adapter/v2"

	"github.com/triggermesh/triggermesh/pkg/sources/adapter/salesforcesource"
)

func main() {
	// JWT package marshals Audience as array even if there is only one element in it. This does not
	// seem to be supported by Salesforce. By setting the following option to false we tell the imported
	// library to marshal single item Audience array as a string.
	jwt.MarshalSingleStringAsArray = false

	adapter.Main("salesforce", salesforcesource.NewEnvConfig, salesforcesource.NewAdapter)
}
