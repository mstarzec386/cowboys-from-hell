# cowboys-from-hell

# Preparations

Install [minikube](https://minikube.sigs.k8s.io/docs/start/)


Start minikube
```bash
minikube start
eval $(minikube docker-env)
alias kubectl="minikube kubectl --"
```

Build docker images
```bash
cd cowboys
docker build -t cowboy:latest -f ../docker/cowboy/Dockerfile .
docker build -t game-master:latest -f ../docker/game-master/Dockerfile .
```

Cowboys logs in second terminal
```bash
alias kubectl="minikube kubectl --"
while [ true ]; do kubectl logs -f -l run=cowboys --max-log-requests=20 --tail=10000; done
```

Game-master logs in third terminal
```bash
alias kubectl="minikube kubectl --"
while [ true ]; do kubectl logs -f deployment/game-master ; done
```

K8s deployment
```bash
cd ../k8s
kubectl apply -f game-master-service.yaml
kubectl apply -f game-master.yaml
kubectl apply -f cowboys.yaml
```


Cleanup
```bash
kubectl delete -f game-master-service.yaml
kubectl delete -f game-master.yaml
kubectl delete -f cowboys.yaml
minikube stop
minikube delete
```

