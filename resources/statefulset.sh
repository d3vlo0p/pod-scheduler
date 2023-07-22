#!/bin/bash
kubectl scale statefoulset {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} --replicas={{.Action.Replicas}}