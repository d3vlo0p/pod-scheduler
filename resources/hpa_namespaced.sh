#!/bin/bash
list="{{StringsJoin .Schedule.Spec.Namespaces ","}}"
for namespace in $list; do
    echo "Namespace: $namespace, finding HPA"
    hpa_names=$(kubectl get hpa -n $namespace {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} -o name | cut -d '/' -f 2)
    for hpa_name in $hpa_names; do
      echo "HPA: $hpa_name set minReplicas: {{.Action.MinReplicas}} maxReplicas: {{.Action.MaxReplicas}}"
      kubectl patch hpa $hpa_name -n $namespace -p '{"spec":{"minReplicas": {{.Action.MinReplicas}}, "maxReplicas": {{.Action.MaxReplicas}} }}'
    done
done