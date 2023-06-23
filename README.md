# Pod Scheduler
automates the number of replicas of pods

Schedule CRD
+ selector labels
+ scale: []
  + cron
  + replicas


this operator manages cron jobs that executes scaling task.

```bash
kubectl scale deployments -l <label-key>=<label-value> --replicas=0
kubectl scale statefulsets -l <label-key>=<label-value> --replicas=0
```