
# Bumblebee ![CircleCI](https://circleci.com/gh/mahjong-contributions/go-bumblebee-jazz.svg?style=svg) ![GoReport](https://goreportcard.com/badge/github.com/mahjong-contributions/go-bumblebee-jazz)

<img src="./docs/anything.jpg" alt="bumblebee" width="80">

It's a sample application, called Bumblebee,  to query air quality around world, which's designed as containerized microservices on Kubernetes. Developers can leverage that to experience various features on Kubernets, such DevOps tooling, observability components, service mesh, etc. 




## What's inside

<img src="./docs/what.png" alt="what">

- Containerized Microservices
- Multiple Language: Go / Typescript / Python
- Multiple protocols: WebSocket / GraphQL / gRPC / HTTP
- Integration with third parties with secret management
- Storage: Local files / MySQL / Caching (Redis)



## How to experiment 

Currently Bumblebee was running on environment, which's provisioned by [Mahjong](https://github.com/awslabs/aws-solutions-assembler) on AWS. Walk through following guides to trial out various features:

### Provision on AWS with Mahjong

### CI/CD
- [Automated release pipline](docs/cicd.md#)
- Blue/Green deployment with ArgoCD Rollouts 
- Canary deployment with ArgoCD Rollouts 

### Traffic management
- Tripping the circuit breaker
- Fault Injection

### Observerbility: metrics, tracing, logs
- Monitoring with Prometheus
- Dashboard with Grafana
- Distributed tracing with Jaeger
- Log analysis with FluentBit & Elasticsearch

### Load testing
- Load testing with Locust
- Load testing with Fortio

### Secret management
- Public secret with Sealed-Secret


# microservice-all-in-one


## Setup KUBECONFIG


## Setup Ingress Controller

```bash
# Install ingress controller 
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/aws/deploy.yaml

```

## Key features
### CI/CD

1. ArgoCD web console

```bash
ARGOCD_ENDPOINT=`kubectl get svc argocd-server -n argocd -o json  |jq -r '.status.loadBalancer.ingress[].hostname'`
ARGOCD_WEB_URL=https://$ARGOCD_ENDPOINT/argocd

# Only for MacOS
open $ARGOCD_WEB_URL

# Otherwise you can copy URL from `echo $ARGOCD_WEB_URL` to your browser

# admin/argocd-server-cf87b5c86-5xlv2 for user and password to login into web console
```

2. Modify code and trigger automated release pipeline
```bash
# Visit demo application
INGRESS_ENDPOINT=`kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'`

DEMO_URL_BJ=http://$INGRESS_ENDPOINT/air/v1/city/beijing
DEMO_URL_HK=http://$INGRESS_ENDPOINT/air/v1/city/hongkong

open $DEMO_URL_BJ
open $DEMO_URL_HK

# Modify code as you want, for example: AQI Standar API

git clone git@github.com:mahjong-contributions/go-bumblebee-jazz.git
cd go-bumblebee-jazz/src/air/apis/v1/
ls
vi handlerv1.go
#Go to AQIStandard(), modify content and save
git status
git add -A
git commit -m "chore: modify AQIStandard API for demo"
git push

# Once push to GitHub and automated pipeline will be triggered

```

3. Canary deployment demo

```bash
# Install Argocd Rollouts plugin: https://argoproj.github.io/argo-rollouts/installation/#kubectl-plugin-installation

# Setup canary deployment
cd manifests/atgocd/canary
kubectl apply -f ns.yaml
kubectl apply -f air-secret.yaml -n go-bumblebee-canary
kubectl apply -f redis.yaml -n go-bumblebee-canary
kubectl apply -f air-rollout.yaml -n go-bumblebee-canary
kubectl argo rollouts get rollout air -n go-bumblebee-canary --watch


kubectl apply -f was-gateway.yaml  -n go-bumblebee-canary
kubectl apply -f was-virtual-service.yaml  -n go-bumblebee-canary
kubectl apply -f was-rollout.yaml -n go-bumblebee-canary
kubectl argo rollouts get rollout was -n go-bumblebee-canary --watch

# Check out jsons output with specific version
INGRESS_ENDPOINT=`kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'`
curl -s -H "Host: was.canary" http://$INGRESS_ENDPOINT/air/v1/city/beijing |jq

# Promote deployment for was
kubectl argo rollouts promote was -n go-bumblebee-canary

# Abort deployment for was
kubectl argo rollouts abort was -n go-bumblebee-canary

# Edit air-rollout.yaml to trial canary deployment and promote 
kubectl argo rollouts abort was -n go-bumblebee-canary
curl -s -H "Host: was.canary" http://$INGRESS_ENDPOINT/air/v1/city/beijing |jq

```


### Traffic management
1. Request routing

2. Circuit Breaking

- Tripping the circuit breaker with command

```bash
# Apply load testing tool to add pressure 
kubectl apply -f fortio-deploy.yaml -n go-bumblebee
FORTIO_POD=$(kubectl get pods -lapp=fortio -o 'jsonpath={.items[0].metadata.name}' -n go-bumblebee)

# Check out result to confirm things' ready
kubectl -n go-bumblebee exec -it "$FORTIO_POD"  -c fortio -- /usr/bin/fortio load -curl http://air:9011/air/v1/city/auckland

kubectl -n go-bumblebee exec -it "$FORTIO_POD"  -c fortio -- /usr/bin/fortio load -c 3 -qps 0 -n 20 -loglevel Warning http://air:9011/air/v1/city/auckland


# Apply Destination Rule to service Air within namespace: go-bumblebee
kubectl -n go-bumblebee apply -f air-destination-rule.yaml
kubectl get destinationrule -n go-bumblebee
kubectl get destinationrule -n go-bumblebee -o yaml

# Adding more presure on it with currency=2
kubectl -n go-bumblebee exec -it "$FORTIO_POD"  -c fortio -- /usr/bin/fortio load -c 2 -qps 0 -n 20 -loglevel Warning http://air:9011/air/v1/city/auckland

# Adding more presure on it with currency=3
kubectl -n go-bumblebee exec -it "$FORTIO_POD"  -c fortio -- /usr/bin/fortio load -c 3 -qps 0 -n 20 -loglevel Warning http://air:9011/air/v1/city/auckland

# Check out key logs to confirm overflow -> circuit breaking
kubectl -n go-bumblebee exec "$FORTIO_POD" -c istio-proxy -- pilot-agent request GET stats | grep air | grep pending

# Clean up 
kubectl -n go-bumblebee delete -f air-destination-rule.yaml
kubectl -n go-bumblebee delete -f fortio-deploy.yaml

```

- Observing the circuit breaker with Locust

```bash

# Apply Destination Rule to service Air within namespace: go-bumblebee
kubectl -n go-bumblebee apply -f air-destination-rule.yaml
kubectl get destinationrule -n go-bumblebee
kubectl get destinationrule -n go-bumblebee -o yaml

# start with Totol number:1 & Hatch rate:1

# start with Totol number:20 & Hatch rate:2
# start with Totol number:20 & Hatch rate:3
```

3. Fault Injection
3.1 HTTP deplay
```bash
# Setup testing services: Was
kubectl -n go-bumblebee apply -f was-deployment-v1.yaml
kubectl -n go-bumblebee apply -f was-deployment-v1.yaml

kubectl -n go-bumblebee apply -f was-httpx-service.yaml
kubectl -n go-bumblebee apply -f was-destination-rule.yaml
kubectl -n go-bumblebee apply -f gateway.yaml
kubectl -n go-bumblebee apply -f virtual-service-was-test-delay.yaml

# Check out virtual service / destination rule
kubectl -n go-bumblebee get virtualservices/was-fault-injection -o yaml
kubectl -n go-bumblebee get destinationrules/was-httpx -o yaml

# Triral

INGRESS_ENDPOINT=`kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'`
curl -v --trace-time -H "Host: was.fault.injection" -H "end-user: aws.builder" http://$INGRESS_ENDPOINT/air/v1/city/beijing


curl -v --trace-time -H "Host: was.fault.injection" http://$INGRESS_ENDPOINT/air/v1/city/beijing

```

3.2 HTTP abort
```bash
# Setup virtual service
kubectl -n go-bumblebee apply -f virtual-service-was-test-abort.yaml

# Check out virtual service
kubectl -n go-bumblebee get virtualservices/was-abort-injection -o yaml

# Triral

INGRESS_ENDPOINT=`kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'`
curl -v --trace-time -H "Host: was.abort.injection" -H "end-user: aws.builder" http://$INGRESS_ENDPOINT/air/v1/city/beijing


curl -v --trace-time -H "Host: was.abort.injection" http://$INGRESS_ENDPOINT/air/v1/city/beijing

````

### Observability 

1. Expose endpoints through nginx-ingress controller

```bash
POD_NAMESPACE=ingress-nginx
POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app.kubernetes.io/name=ingress-nginx --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}')

kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- /nginx-ingress-controller --

# Apply yaml to expose endpoints
htpasswd -c auth admin abcd1234
kubectl create secret generic basic-auth --from-file=auth -n istio-system
kubectl get secret basic-auth -n istio-system -o yaml

kubectl apply -f observability-console.yaml -n istio-system
```

Notice:

All endpoints are exposed with header `Host`, so you can install `modheader` as header modifier for conveniency. 

2. Grafana
3. Prometheus 
4. Jaeger Tracing
5. Kiali

### Load Testing

1. Locust
2. Fortio


