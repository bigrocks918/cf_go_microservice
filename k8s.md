# K8S


minikube start --nodes=2

minikube addons enable ingress


minikube status
docker ps
kubectl get nodes
kubectl get pods -A
kubectl get pods
kubectl apply -f deployment.yml
kubectl get pods
kubectl get svc
kubectl get deployments

# hit the app
expose deployment myapp --type=LoadBalancer --port=8080 --target-port=80
minikube tunnel

# stop service/app
kubectl delete deployments <deployment>
kubectl delete services <service>

# stop minikube
minikube stop