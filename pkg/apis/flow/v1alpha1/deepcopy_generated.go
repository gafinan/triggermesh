//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1 "knative.dev/pkg/apis/duck/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Path) DeepCopyInto(out *Path) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Path.
func (in *Path) DeepCopy() *Path {
	if in == nil {
		return nil
	}
	out := new(Path)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Synchronizer) DeepCopyInto(out *Synchronizer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Synchronizer.
func (in *Synchronizer) DeepCopy() *Synchronizer {
	if in == nil {
		return nil
	}
	out := new(Synchronizer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Synchronizer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SynchronizerList) DeepCopyInto(out *SynchronizerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Synchronizer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SynchronizerList.
func (in *SynchronizerList) DeepCopy() *SynchronizerList {
	if in == nil {
		return nil
	}
	out := new(SynchronizerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SynchronizerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SynchronizerSpec) DeepCopyInto(out *SynchronizerSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SynchronizerSpec.
func (in *SynchronizerSpec) DeepCopy() *SynchronizerSpec {
	if in == nil {
		return nil
	}
	out := new(SynchronizerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SynchronizerStatus) DeepCopyInto(out *SynchronizerStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	in.AddressStatus.DeepCopyInto(&out.AddressStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SynchronizerStatus.
func (in *SynchronizerStatus) DeepCopy() *SynchronizerStatus {
	if in == nil {
		return nil
	}
	out := new(SynchronizerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Transform) DeepCopyInto(out *Transform) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]Path, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Transform.
func (in *Transform) DeepCopy() *Transform {
	if in == nil {
		return nil
	}
	out := new(Transform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Transformation) DeepCopyInto(out *Transformation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Transformation.
func (in *Transformation) DeepCopy() *Transformation {
	if in == nil {
		return nil
	}
	out := new(Transformation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Transformation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TransformationList) DeepCopyInto(out *TransformationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Transformation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TransformationList.
func (in *TransformationList) DeepCopy() *TransformationList {
	if in == nil {
		return nil
	}
	out := new(TransformationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TransformationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TransformationSpec) DeepCopyInto(out *TransformationSpec) {
	*out = *in
	in.Sink.DeepCopyInto(&out.Sink)
	if in.Context != nil {
		in, out := &in.Context, &out.Context
		*out = make([]Transform, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make([]Transform, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TransformationSpec.
func (in *TransformationSpec) DeepCopy() *TransformationSpec {
	if in == nil {
		return nil
	}
	out := new(TransformationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TransformationStatus) DeepCopyInto(out *TransformationStatus) {
	*out = *in
	in.SourceStatus.DeepCopyInto(&out.SourceStatus)
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(v1.Addressable)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TransformationStatus.
func (in *TransformationStatus) DeepCopy() *TransformationStatus {
	if in == nil {
		return nil
	}
	out := new(TransformationStatus)
	in.DeepCopyInto(out)
	return out
}
