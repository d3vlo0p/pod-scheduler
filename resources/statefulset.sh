#!/bin/bash
kubectl scale statefulset {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} --replicas={{.Action.Replicas}}