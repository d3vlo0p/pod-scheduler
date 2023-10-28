//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 Marco Bonacchi.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSchedule) DeepCopyInto(out *ClusterSchedule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSchedule.
func (in *ClusterSchedule) DeepCopy() *ClusterSchedule {
	if in == nil {
		return nil
	}
	out := new(ClusterSchedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterSchedule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterScheduleList) DeepCopyInto(out *ClusterScheduleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterSchedule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterScheduleList.
func (in *ClusterScheduleList) DeepCopy() *ClusterScheduleList {
	if in == nil {
		return nil
	}
	out := new(ClusterScheduleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterScheduleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterScheduleSpec) DeepCopyInto(out *ClusterScheduleSpec) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Schedules != nil {
		in, out := &in.Schedules, &out.Schedules
		*out = make([]ScheduleAction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterScheduleSpec.
func (in *ClusterScheduleSpec) DeepCopy() *ClusterScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterScheduleStatus) DeepCopyInto(out *ClusterScheduleStatus) {
	*out = *in
	in.LastRunTime.DeepCopyInto(&out.LastRunTime)
	if in.CronJobs != nil {
		in, out := &in.CronJobs, &out.CronJobs
		*out = make([]CronJob, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterScheduleStatus.
func (in *ClusterScheduleStatus) DeepCopy() *ClusterScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterScheduleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronJob) DeepCopyInto(out *CronJob) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronJob.
func (in *CronJob) DeepCopy() *CronJob {
	if in == nil {
		return nil
	}
	out := new(CronJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Schedule) DeepCopyInto(out *Schedule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Schedule.
func (in *Schedule) DeepCopy() *Schedule {
	if in == nil {
		return nil
	}
	out := new(Schedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Schedule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleAction) DeepCopyInto(out *ScheduleAction) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleAction.
func (in *ScheduleAction) DeepCopy() *ScheduleAction {
	if in == nil {
		return nil
	}
	out := new(ScheduleAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleList) DeepCopyInto(out *ScheduleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Schedule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleList.
func (in *ScheduleList) DeepCopy() *ScheduleList {
	if in == nil {
		return nil
	}
	out := new(ScheduleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScheduleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleSpec) DeepCopyInto(out *ScheduleSpec) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Schedules != nil {
		in, out := &in.Schedules, &out.Schedules
		*out = make([]ScheduleAction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleSpec.
func (in *ScheduleSpec) DeepCopy() *ScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(ScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleStatus) DeepCopyInto(out *ScheduleStatus) {
	*out = *in
	in.LastRunTime.DeepCopyInto(&out.LastRunTime)
	if in.CronJobs != nil {
		in, out := &in.CronJobs, &out.CronJobs
		*out = make([]CronJob, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleStatus.
func (in *ScheduleStatus) DeepCopy() *ScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(ScheduleStatus)
	in.DeepCopyInto(out)
	return out
}
