package v1alpha1

type ScheduleAction struct {
	Name     string `json:"name"`
	Cron     string `json:"cron"`
	Replicas int    `json:"replicas"`
}

type CronJob struct {
	Name string `json:"name"`
}
