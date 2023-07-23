# Pod Scheduler
automates the number of replicas of pods with cron expression

Schedule CR example
```yaml
apiVersion: pod.loop.dev/v1alpha1
kind: Schedule
metadata:
  name: deployment-sample
spec:
  matchLabels:
    schedule: scalein-at-1-scaleout-at-5
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
  name: hpa-sample
spec:
  matchLabels:
    schedule: lower-at-1-higher-at-5
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

apiVersion: pod.loop.dev/v1alpha1
kind: ClusterSchedule
metadata:
  name: hpa-cluster-sample
spec:
  matchLabels:
    schedule: lower-at-1-higher-at-5
  matchType: HorizontalPodAutoscaler
  namespaces:
    - default
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
kubectl apply -f https://raw.githubusercontent.com/d3vlo0p/pod-scheduler/main/config/crd/bases/pod.loop.dev_clusterschedules.yaml

# installing operator with helm
helm install my-scheduler oci://ghcr.io/d3vlo0p/pod-scheduler-operator
```