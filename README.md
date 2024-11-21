# Docker (used --build when change has edit code)  
docker compose up --build

# Remove Docker
docker compose down 

# apply the changes:
kubectl apply -f k8s/player-service.yaml
kubectl apply -f k8s/table-service.yaml
kubectl apply -f k8s/mongo.yaml

# For check k8s is working 
kubectl get services
kubectl get pods
kubectl get deployments

# If need to remove k8s 
kubectl delete deployment player-service table-service mongodb
kubectl delete svc player-service table-service mongodb

# For player-service forward to localhost 
kubectl port-forward player-service 8888:8888

# For table-service forward to localhost
kubectl port-forward table-service 8889:8889

# For test  forward to localhost (response :=> etc. ,"name of service")
curl http://localhost:8889 
curl http://localhost:8888

# MongoDB
show dbs
use match_results
show collections
db.results.find().pretty()



