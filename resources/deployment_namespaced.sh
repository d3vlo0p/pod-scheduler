#!/bin/bash
list="{{StringsJoin .Schedule.Spec.Namespaces ","}}"
for namespace in ${list//,/ }; do
    echo "Namespace: $namespace, applying scale to deployments"
    kubectl scale deployment -n $namespace {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} --replicas={{.Action.Replicas}}
done