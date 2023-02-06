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

Game-master logs in third terminal (optional)
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

# Example logs from cowboys

```
Registered as Name: Sam, Health: 10, Damage: 1
Registered as Name: Philip, Health: 15, Damage: 1
Registered as Name: Peter, Health: 5, Damage: 3
Registered as Name: John, Health: 10, Damage: 1
Registered as Name: Bill, Health: 8, Damage: 2
...
Philip: Oh someones hits (2) me (0) ಥ_ಥ and killed (✖╭╮✖)
Bill: Oh someones hits (1) me (4) ಥ_ಥ
John: Hit victim Philip (3)
Bill: Hit victim Philip (11)
Bill: Hit victim Philip (7)
Bill: Hit victim Sam (2)
John: Oh someones hits (1) me (2) ಥ_ಥ
John: Hit victim Bill (4)
John: Oh someones hits (1) me (1) ಥ_ಥ
Bill: Hit victim Philip (2)
Bill: Oh someones hits (1) me (3) ಥ_ಥ
John: Oh someones hits (2) me (0) ಥ_ಥ and killed (✖╭╮✖)
Bill: Hit victim John (1)
Bill: I WIN!!!! \(ᵔᵕᵔ)/   \(ᵔᵕᵔ)/   \(ᵔᵕᵔ)/
```