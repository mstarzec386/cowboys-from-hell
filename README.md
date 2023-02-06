# cowboys-from-hell

# Preparations

Install [minikube](https://minikube.sigs.k8s.io/docs/start/)
Install redis-cli (deb)
```bash
apt install -y redis-tools
```
or for yum
```bash
yum install redis
```

Start minikube
```bash
minikube start
```

# Deployment
Envs and alias for kubectl
```bash
eval $(minikube docker-env)
alias kubectl="minikube kubectl --"
```

## Building
Build docker images
```bash
cd cowboys
docker build -t cowboy:latest -f ../docker/cowboy/Dockerfile .
docker build -t game-master:latest -f ../docker/game-master/Dockerfile .
```

## Redis deployment
K8s deployment redis
```bash
cd ../k8s
kubectl apply -f redis-service.yaml
sleep 2
kubectl apply -f redis.yaml
sleep 2
```

In a different terminal run port-forward to inject data to the redis
```bash
alias kubectl="minikube kubectl --"
kubectl port-forward service/redis 6379:6379
```

Inject data (use first terminal)
```bash
redis-cli -x SET init < cowboys.json
```

## Watch logs
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


## Apps deployment
K8s deployment apps
```bash
kubectl apply -f game-master-service.yaml
sleep 2
kubectl apply -f game-master.yaml
sleep 2
kubectl apply -f cowboys.yaml
```


# Cleanup
```bash
kubectl delete -f game-master.yaml
kubectl delete -f game-master-service.yaml
kubectl delete -f cowboys.yaml
kubectl delete -f redis.yaml
kubectl delete -f redis-service.yaml

minikube stop

#optional
minikube delete
```

