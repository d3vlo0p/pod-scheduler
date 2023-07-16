# Pod Scheduler
automates the number of replicas of pods with cron expression

Schedule CR example
```yaml
apiVersion: pod.loop.dev/v1alpha1
kind: Schedule
metadata:
  labels:
    app.kubernetes.io/name: schedule
    app.kubernetes.io/instance: schedule-sample
    app.kubernetes.io/part-of: pod-scheduler
    app.kubernetes.io/created-by: pod-scheduler
  name: deployment-sample
spec:
  matchLabels:
    schedule: scalein-at-0-scaleout-at-5
  matchType: Deployment
  schedules:
    - cron: 0,10,20,30,40,50 * * * *
      name: scale-in
      replicas: 1
    - cron: 5,15,25,35,45,55 * * * *
      name: scale-out
      replicas: 2

apiVersion: pod.loop.dev/v1alpha1
kind: Schedule
metadata:
  labels:
    app.kubernetes.io/name: schedule
    app.kubernetes.io/instance: schedule-sample
    app.kubernetes.io/part-of: pod-scheduler
    app.kubernetes.io/created-by: pod-scheduler
  name: hpa-sample
spec:
  matchLabels:
    schedule: lower-at-0-higher-at-5
  matchType: HorizontalPodAutoscaler
  schedules:
    - cron: 0,10,20,30,40,50 * * * *
      name: lower-min
      minReplicas: 1
      maxReplicas: 5
    - cron: 5,15,25,35,45,55 * * * *
      name: higher-min
      minReplicas: 3
      maxReplicas: 5
```


## Installation
```bash
# installing CRDs
kubectl apply -f https://raw.githubusercontent.com/d3vlo0p/pod-scheduler/main/config/crd/bases/pod.loop.dev_schedules.yaml

# installing operator with helm
helm install my-scheduler oci://ghcr.io/d3vlo0p/pod-scheduler-operator
```