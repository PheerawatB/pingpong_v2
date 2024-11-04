# apply the changes:
kubectl apply -f k8s/player-service-service.yaml
kubectl apply -f k8s/table-service-service.yaml

# For check 
kubectl get services
kubectl get pods
kubectl get deployments

# For player-service
kubectl port-forward service/player-service 8888:8888

# For table-service
kubectl port-forward service/table-service 8889:8889

# For test 
curl http://localhost:8889
curl http://localhost:8888
