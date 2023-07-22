#!/bin/bash
list="{{StringsJoin .Schedule.Spec.Namespaces ","}}"
for namespace in ${list//,/ }; do
    echo "Namespace: $namespace, applying scale to statefulset"
    kubectl scale statefoulset -n $namespace {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} --replicas={{.Action.Replicas}}
done