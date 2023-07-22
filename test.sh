#!/bin/bash
list="default"
for namespace in $list; do
    echo "Namespace: $namespace, finding HPA"
    hpa_names=$(kubectl get hpa -n $namespace -l schedule=test -o name | cut -d '/' -f 2)
    for hpa_name in $hpa_names; do
      echo "HPA: $hpa_name set minReplicas: 2 maxReplicas: 3"
      kubectl patch hpa $hpa_name -n $namespace -p '{"spec":{"minReplicas": 2, "maxReplicas": 5 }}'
    done
done