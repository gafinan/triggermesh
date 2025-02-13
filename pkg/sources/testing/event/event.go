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

// Package event contains test helpers for Kubernetes API events.
package event

import "fmt"

// Eventf returns the attributes of an API event in the format returned by
// Kubernetes' FakeRecorder.
func Eventf(eventtype, reason, messageFmt string, args ...interface{}) string {
	return fmt.Sprintf(eventtype+" "+reason+" "+messageFmt, args...)
}
