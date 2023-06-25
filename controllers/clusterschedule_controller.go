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

package controllers

import (
	"context"
	"fmt"
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	podv1alpha1 "github.com/d3vlo0p/pod-scheduler/api/v1alpha1"
)

// ClusterScheduleReconciler reconciles a ClusterSchedule object
type ClusterScheduleReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	JobImage       string
	ServiceAccount string
	Namespace      string
}

//+kubebuilder:rbac:groups=pod.loop.dev,resources=clusterschedules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pod.loop.dev,resources=clusterschedules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pod.loop.dev,resources=clusterschedules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterSchedule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ClusterScheduleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info(fmt.Sprintf("reconciling object %#q", req.NamespacedName))

	instance := &podv1alpha1.ClusterSchedule{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("ClusterSchedule resource not found. object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Info("Failed to get ClusterSchedule resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}

	// find cronJobs managed by this resource
	oldResouces := map[string]*batchv1.CronJob{}
	for _, cj := range instance.Status.CronJobs {
		cronJob := &batchv1.CronJob{}
		err = r.Get(ctx, client.ObjectKey{Name: cj.Name, Namespace: r.Namespace}, cronJob)
		if err != nil {
			if !errors.IsNotFound(err) {
				logger.Info("Failed to get cronjob. Re-running reconcile.")
				return ctrl.Result{}, err
			}
			logger.Info(fmt.Sprintf("cronjob %s/%s not found", cj.Name, req.Namespace))
		} else {
			// managed resource do exist
			oldResouces[cj.Name] = cronJob
		}
	}

	instance.Status.CronJobs = []podv1alpha1.CronJob{}
	// verify if the cronjobs are matching the schedule spec.
	// if not create a new one or modify the existing one
	for _, scheduleAction := range instance.Spec.Schedules {
		var newCj *batchv1.CronJob
		jobName := GetCronJobName(instance.Name, scheduleAction.Name)
		cj, ok := oldResouces[jobName]
		if !ok {
			// create a new cronjob and add name to status
			newCj = r.cronJobForSchedule(instance, scheduleAction)
			err = r.Create(ctx, newCj)
			if err != nil {
				logger.Info("Failed to create Schedule CronJob. Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			// cronjob exist, remove the cronjob from the map to check later if there are job left to remove
			delete(oldResouces, jobName)
			if cj.Spec.Schedule != scheduleAction.Cron ||
				cj.Annotations["pod-scheduler.loop.dev/replicas"] != strconv.Itoa(scheduleAction.Replicas) ||
				cj.Annotations["pod-scheduler.loop.dev/labelSelectors"] != ConvertMapToString(instance.Spec.MatchLabels) {
				// configuration changed, replace cronjob
				err := r.Delete(ctx, cj, &client.DeleteOptions{})
				if err != nil {
					logger.Info("Failed to delete cronjob")
					return ctrl.Result{}, err
				}
				newCj = r.cronJobForSchedule(instance, scheduleAction)
				err = r.Create(ctx, newCj)
				if err != nil {
					logger.Info("Failed to create Schedule CronJob. Re-running reconcile.")
					return ctrl.Result{}, err
				}
			}
		}
		instance.Status.CronJobs = append(instance.Status.CronJobs, podv1alpha1.CronJob{Name: jobName})
	}

	// check if some action has been removed from the schedule spec, but cronjob is still active then remove it
	for _, cj := range oldResouces {
		err := r.Delete(ctx, cj, &client.DeleteOptions{})
		if err != nil {
			logger.Info("Failed to delete cronjob")
			return ctrl.Result{}, err
		}
	}

	// Set status
	instance.Status.LastRunTime = metav1.Now()
	r.Status().Update(ctx, instance)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podv1alpha1.ClusterSchedule{}).
		Complete(r)
}

func (r *ClusterScheduleReconciler) cronJobForSchedule(schedule *podv1alpha1.ClusterSchedule, action podv1alpha1.ScheduleAction) *batchv1.CronJob {
	jobName := GetCronJobName(schedule.Name, action.Name)
	job := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: r.Namespace,
			Annotations: map[string]string{
				"pod-scheduler.loop.dev/replicas":       strconv.Itoa(action.Replicas),
				"pod-scheduler.loop.dev/labelSelectors": ConvertMapToString(schedule.Spec.MatchLabels),
			},
		},
		Spec: batchv1.CronJobSpec{
			ConcurrencyPolicy: "Replace",
			Schedule:          action.Cron,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: GenerateLabelsForApp(jobName),
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: GenerateLabelsForApp(jobName),
						},
						Spec: corev1.PodSpec{
							RestartPolicy:      corev1.RestartPolicyNever,
							ServiceAccountName: r.ServiceAccount,
							Containers: []corev1.Container{
								{
									Name:  "pod-scheduler-deployment",
									Image: r.JobImage,
									Args:  GenerateArgs("deployment", schedule.Spec.MatchLabels, action.Replicas, true),
								},
								{
									Name:  "pod-scheduler-statefulset",
									Image: r.JobImage,
									Args:  GenerateArgs("statefulset", schedule.Spec.MatchLabels, action.Replicas, true),
								},
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(schedule, job, r.Scheme)
	return job
}
