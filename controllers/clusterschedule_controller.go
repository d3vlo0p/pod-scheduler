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
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

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
	Templates      *template.Template
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
	oldResouces := map[string]ScheduleResources{}
	for _, cmcj := range instance.Status.CronJobs {
		resource := ScheduleResources{}

		cronJob := &batchv1.CronJob{}
		err1 := r.Get(ctx, client.ObjectKey{Name: cmcj.Job, Namespace: r.Namespace}, cronJob)
		if err1 != nil {
			if !errors.IsNotFound(err1) {
				logger.Info("Failed to get cronjob. Re-running reconcile.")
				return ctrl.Result{}, err1
			}
			logger.Info(fmt.Sprintf("cronjob %s/%s not found", cmcj.Job, r.Namespace))
		} else {
			logger.Info(fmt.Sprintf("cronjob %s/%s found", cmcj.Job, r.Namespace))
			resource.cronjob = cronJob
		}

		configMap := &corev1.ConfigMap{}
		err2 := r.Get(ctx, client.ObjectKey{Name: cmcj.ConfigMap, Namespace: r.Namespace}, configMap)
		if err2 != nil {
			if !errors.IsNotFound(err2) {
				logger.Info("Failed to get configMap. Re-running reconcile.")
				return ctrl.Result{}, err1
			}
			logger.Info(fmt.Sprintf("configMap %s/%s not found", cmcj.ConfigMap, r.Namespace))
		} else {
			logger.Info(fmt.Sprintf("configMap %s/%s found", cmcj.ConfigMap, r.Namespace))
			resource.configmap = configMap
		}

		if err1 == nil && err2 == nil {
			// managed resource do exist
			oldResouces[cmcj.Name] = resource
		}
	}

	instance.Status.CronJobs = []podv1alpha1.CronJob{}
	// verify if the cronjobs are matching the schedule spec.
	// if not create a new one or modify the existing one
	if instance.Spec.Enabled == nil || *instance.Spec.Enabled {
		for _, scheduleAction := range instance.Spec.Schedules {
			if scheduleAction.Enabled == nil || *scheduleAction.Enabled {
				var resource ScheduleResources
				jobName := GetScheduleActionName(instance.Name, scheduleAction.Name)
				cmcj, ok := oldResouces[jobName]
				if !ok {
					logger.Info("New Schedule action found, creating new cronjob & config map")
					resource, err = r.createCronJob(ctx, instance, scheduleAction)
					if err != nil {
						logger.Info("Failed to create Schedule CronJob. Re-running reconcile.")
						return ctrl.Result{}, err
					}
				} else {
					// cronjob exist, remove the cronjob from the map to check later if there are job left to remove
					delete(oldResouces, jobName)
					if cmcj.cronjob.Spec.Schedule != scheduleAction.Cron ||
						cmcj.configmap.Annotations["pod-scheduler.loop.dev/replicas"] != ReplicasForLabels(scheduleAction) ||
						cmcj.configmap.Annotations["pod-scheduler.loop.dev/labelSelectors"] != ConvertMapToString(instance.Spec.MatchLabels) ||
						cmcj.configmap.Annotations["pod-scheduler.loop.dev/resource"] != instance.Spec.MatchType.String() ||
						cmcj.configmap.Annotations["pod-scheduler.loop.dev/namespaces"] != strings.Join(instance.Spec.Namespaces, ",") { // configuration changed, replace cronjob & config map
						logger.Info("Diff found, replacing cronjob & config map")
						err := r.Delete(ctx, cmcj.cronjob, &client.DeleteOptions{})
						if err != nil {
							logger.Info("Failed to delete Schedule cronjob")
							return ctrl.Result{}, err
						}
						err = r.Delete(ctx, cmcj.configmap, &client.DeleteOptions{})
						if err != nil {
							logger.Info("Failed to delete Schedule cronjob configmap")
							return ctrl.Result{}, err
						}
						resource, err = r.createCronJob(ctx, instance, scheduleAction)
						if err != nil {
							logger.Info("Failed to create Schedule CronJob. Re-running reconcile.")
							return ctrl.Result{}, err
						}
					} else {
						logger.Info("No diff found, keeping existing cronjob & config map")
						resource = cmcj
					}
				}
				instance.Status.CronJobs = append(instance.Status.CronJobs, podv1alpha1.CronJob{
					Name:      jobName,
					Job:       resource.cronjob.Name,
					ConfigMap: resource.configmap.Name,
				})
			}
		}
	}

	// check if some action has been removed from the schedule spec, but cronjob is still active then remove it
	for _, cmcj := range oldResouces {
		err := r.Delete(ctx, cmcj.cronjob, &client.DeleteOptions{})
		if err != nil {
			logger.Info("Failed to delete cronjob")
			return ctrl.Result{}, err
		}
		err = r.Delete(ctx, cmcj.configmap, &client.DeleteOptions{})
		if err != nil {
			logger.Info("Failed to delete configmap")
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

func (r *ClusterScheduleReconciler) createCronJob(ctx context.Context, schedule *podv1alpha1.ClusterSchedule, action podv1alpha1.ScheduleAction) (ScheduleResources, error) {
	cm, err := r.configMapForSchedule(schedule, action)
	if err != nil {
		return ScheduleResources{}, fmt.Errorf("failed to generate Schedule ConfigMap: %w", err)
	}
	err = r.Create(ctx, cm)
	if err != nil {
		return ScheduleResources{}, fmt.Errorf("failed to create Schedule ConfigMap: %w", err)
	}
	cj, err := r.cronJobForSchedule(schedule, action, cm)
	if err != nil {
		return ScheduleResources{}, fmt.Errorf("failed to generate Schedule CronJob: %w", err)
	}
	err = r.Create(ctx, cj)
	if err != nil {
		return ScheduleResources{}, fmt.Errorf("failed to create Schedule CronJob: %w", err)
	}
	return ScheduleResources{
		cronjob:   cj,
		configmap: cm,
	}, nil
}

func (r *ClusterScheduleReconciler) configMapForSchedule(schedule *podv1alpha1.ClusterSchedule, action podv1alpha1.ScheduleAction) (*corev1.ConfigMap, error) {
	buf := new(bytes.Buffer)
	data := scriptData[*podv1alpha1.ClusterSchedule]{schedule, action}
	if data.Schedule.Spec.MatchType == podv1alpha1.Deployment {
		err := r.Templates.ExecuteTemplate(buf, "deployment_namespaced.sh", data)
		if err != nil {
			return nil, err
		}
	} else if data.Schedule.Spec.MatchType == podv1alpha1.StatefulSet {
		err := r.Templates.ExecuteTemplate(buf, "statefulset_namespaced.sh", data)
		if err != nil {
			return nil, err
		}
	} else if data.Schedule.Spec.MatchType == podv1alpha1.HorizontalPodAutoscaler {
		err := r.Templates.ExecuteTemplate(buf, "hpa_namespaced.sh", data)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("match type not supported")
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetScheduleActionName(data.Schedule.Name, data.Action.Name),
			Namespace: r.Namespace,
			Annotations: map[string]string{
				"pod-scheduler.loop.dev/replicas":       ReplicasForLabels(data.Action),
				"pod-scheduler.loop.dev/labelSelectors": ConvertMapToString(data.Schedule.Spec.MatchLabels),
				"pod-scheduler.loop.dev/resource":       data.Schedule.Spec.MatchType.String(),
				"pod-scheduler.loop.dev/namespaces":     strings.Join(data.Schedule.Spec.Namespaces, ","),
			},
		},
		Data: map[string]string{
			scriptFileName: buf.String(),
		},
	}
	controllerutil.SetControllerReference(data.Schedule, cm, r.Scheme)
	return cm, nil
}

func (r *ClusterScheduleReconciler) cronJobForSchedule(schedule *podv1alpha1.ClusterSchedule, action podv1alpha1.ScheduleAction, configMap *corev1.ConfigMap) (*batchv1.CronJob, error) {
	jobName := GetScheduleActionName(schedule.Name, action.Name)
	volumeMode := new(int32)
	*volumeMode = 0555 // rxrxrx
	job := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: r.Namespace,
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
									Name:    "pod-scheduler",
									Image:   r.JobImage,
									Command: []string{"sh", "-c", "./" + scriptFileName},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "script-vol",
											MountPath: "/" + scriptFileName,
											SubPath:   scriptFileName,
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "script-vol",
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: configMap.Name,
											},
											Items: []corev1.KeyToPath{
												{
													Key:  scriptFileName,
													Path: scriptFileName,
													Mode: volumeMode,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(schedule, job, r.Scheme)
	return job, nil
}
