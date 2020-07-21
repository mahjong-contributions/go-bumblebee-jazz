# Observability

## Expose endpoints through nginx-ingress controller

1. Install nginx-ingress controller

```bash

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/aws/deploy.yaml

POD_NAMESPACE=ingress-nginx
POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app.kubernetes.io/name=ingress-nginx --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}')

kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- /nginx-ingress-controller --

2. Expose related services through nginx-ingress controller.

```bash
cd manifests/monitoring
# Apply yaml to expose endpoints
htpasswd -c auth admin abcd1234
kubectl create secret generic basic-auth --from-file=auth -n istio-system
kubectl get secret basic-auth -n istio-system -o yaml

kubectl apply -f observability-console.yaml -n istio-system

```

Notice:

All endpoints are exposed with header `Host`, so you can install `modheader` as header modifier for conveniency. 

[!!! Screen Recording !!!](https://aws-solutions-assembler.s3-ap-southeast-1.amazonaws.com/sample-box/monitoring.mp4)

## Grafana
## Prometheus 
## Jaeger Tracing
## Kiali