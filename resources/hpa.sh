#!/bin/bash
hpa_names=$(kubectl get hpa {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} -o name | cut -d '/' -f 2)
for hpa_name in $hpa_names; do
  kubectl patch hpa $hpa_name -p '{"spec":{"minReplicas": {{.Action.MinReplicas}}, "maxReplicas": {{.Action.MaxReplicas}} }}'
done