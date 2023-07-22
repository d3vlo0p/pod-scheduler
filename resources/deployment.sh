#!/bin/bash
kubectl scale deployment {{range $key, $value := .Schedule.Spec.MatchLabels}} -l {{$key}}={{$value}} {{end}} --replicas={{.Action.Replicas}}