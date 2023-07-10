package v1alpha1

type ScheduleAction struct {
	Name     string `json:"name"`
	Cron     string `json:"cron"`
	Replicas int    `json:"replicas"`
}

type CronJob struct {
	Name      string `json:"name"`
	Job       string `json:"job"`
	ConfigMap string `json:"configMap"`
}

type ResourceType string

const (
	Deployment  ResourceType = "Deployment"
	StatefulSet ResourceType = "StatefulSet"
	// HorizontalPodAutoscaler ResourceType = "HorizontalPodAutoscaler"
)

func (r ResourceType) String() string {
	return string(r)
}
