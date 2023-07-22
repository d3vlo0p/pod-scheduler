package controllers

import (
	"fmt"
	"strings"

	podv1alpha1 "github.com/d3vlo0p/pod-scheduler/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

const scriptFileName = "script.sh"

type scheduleData interface {
	*podv1alpha1.ClusterSchedule | *podv1alpha1.Schedule
}

type scriptData[T scheduleData] struct {
	Schedule T
	Action   podv1alpha1.ScheduleAction
}

type ScheduleResources struct {
	cronjob   *batchv1.CronJob
	configmap *corev1.ConfigMap
}

func GetScheduleActionName(scheduleName string, actionName string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", scheduleName, actionName))
}

func ConvertMapToString(m map[string]string) string {
	ra := []string{}
	for k, v := range m {
		ra = append(ra, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(ra, ",")
}

func GenerateLabelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name, "app": "pod-scheduler", "cr_type": "schedule"}
}

func ReplicasForLabels(scheduleAction podv1alpha1.ScheduleAction) string {
	return fmt.Sprintf("%d-%d-%d", scheduleAction.Replicas, scheduleAction.MinReplicas, scheduleAction.MaxReplicas)
}
