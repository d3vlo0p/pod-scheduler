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
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: pod-scheduler
  name: schedule-sample
spec:
  matchLabels:
    app: test
  matchType: Deployment
  schedules:
    - cron: 0,10,20,30,40,50 * * * *
      name: prova-scale-in
      replicas: 2
    - cron: 5,15,25,35,45,55 * * * *
      name: prova-scale-out
      replicas: 5
```