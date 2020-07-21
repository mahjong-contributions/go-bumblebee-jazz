# Traffic management



## Tripping the circuit breaker

1. Observing the circuit breaker with Fortio

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

[!!! Screen Recording !!!](https://aws-solutions-assembler.s3-ap-southeast-1.amazonaws.com/sample-box/circuit-breaking.mp4)

2. Observing the circuit breaker with Locust

```bash

# Apply Destination Rule to service Air within namespace: go-bumblebee
kubectl -n go-bumblebee apply -f air-destination-rule.yaml
kubectl get destinationrule -n go-bumblebee
kubectl get destinationrule -n go-bumblebee -o yaml

# start with Totol number:1 & Hatch rate:1

# start with Totol number:20 & Hatch rate:2
# start with Totol number:20 & Hatch rate:3
```

## Fault Injection

1. HTTP deplay
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

2. HTTP abort
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