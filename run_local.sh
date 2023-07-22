export WATCH_NAMESPACE=""
export JOB_CONTAINER_IMAGE=bitnami/kubectl:latest
export CLUSTER_ROLE_NAME=pod-scheduler-replica-modifier-role
export OPERATOR_NAMESPACE=pod-scheduler
export OPERATOR_SERVICE_ACCOUNT_NAME=controller-manager
kubectl apply -k ./config/rbac/
make install run