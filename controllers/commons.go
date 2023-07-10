package controllers

import (
	"fmt"
	"strings"
)

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

func GenerateArgs(resourceType string, labels map[string]string, replicas int, global bool) []string {
	args := []string{"scale", resourceType, fmt.Sprintf("--replicas=%d", replicas)}
	for k, v := range labels {
		args = append(args, "-l", fmt.Sprintf("%s=%s", k, v))
	}
	if global {
		args = append(args, "-A")
	}
	return args
}

func GenerateLabelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name, "app": "pod-scheduler", "cr_type": "schedule"}
}
