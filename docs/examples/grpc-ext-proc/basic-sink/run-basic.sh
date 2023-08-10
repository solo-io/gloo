CLUSTER_NAME="kind"
echo ${CLUSTER_NAME}

kind load docker-image quay.io/solo-io/grpc-ext-proc-example --name ${CLUSTER_NAME}
docker exec -it ${CLUSTER_NAME}-control-plane crictl images | grep grpc-ext-proc


kubectl apply -f resources/deployment.yaml
kubectl apply -f resources/service.yaml

kubectl rollout restart deployments/ext-proc-grpc 
kubectl get deployments/ext-proc-grpc


kubectl patch settings -n gloo-system default --patch-file resources/ext-proc-settings-patch.json --type merge

