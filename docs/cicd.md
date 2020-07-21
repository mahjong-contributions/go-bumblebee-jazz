# CI/CD

## Automated release pipline 

1. High-level diagram for automated release pipeline.

<img src="cicd-diagram.png"
     alt="CI/CD"
     style="margin-left: 10px;" />

2. CI build up with AWS Codebuild.

3. CD handler bu ArgoCD. Go through the following steps to access ArgoCD console.

```bash

ARGOCD_ENDPOINT=`kubectl get svc argocd-server -n argocd -o json  |jq -r '.status.loadBalancer.ingress[].hostname'`
ARGOCD_WEB_URL=https://$ARGOCD_ENDPOINT/argocd

# Only for MacOS
open $ARGOCD_WEB_URL

# Otherwise you can copy URL from `echo $ARGOCD_WEB_URL` to your browser

```
Initial user and password are admin/argocd-server-cf87b5c86-5xlv2 to login into ArgoCD.


4. Modify demo application code in order to trigger automated release pipeline.

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

[!!! Screen Recording !!!](https://aws-solutions-assembler.s3-ap-southeast-1.amazonaws.com/sample-box/cicd-demo.mp4)


## Blue/Green deployment with ArgoCD Rollouts 

## Canary deployment with ArgoCD Rollouts

1. Install Argocd Rollouts plugin to watch status: https://argoproj.github.io/argo-rollouts/installation/#kubectl-plugin-installation.


2. Setup trial for canary deployment 
```bash

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
```

3. Check out result from demo application before deployment.
```bash

# Check out jsons output with specific version
INGRESS_ENDPOINT=`kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'`
curl -s -H "Host: was.canary" http://$INGRESS_ENDPOINT/air/v1/city/beijing |jq
```

4. Edit air-rollout.yaml to trial canary deployment and promote 
```bash
kubectl apply -f air-rollout.yaml -n go-bumblebee-canary

curl -s -H "Host: was.canary" http://$INGRESS_ENDPOINT/air/v1/city/beijing |jq
kubectl argo rollouts promote air -n go-bumblebee-canary

```


5. Experience abort deployment if you want to. 
```bash

# Abort deployment for was
kubectl argo rollouts abort was -n go-bumblebee-canary

```

[!!! Screen Recording !!!](https://aws-solutions-assembler.s3-ap-southeast-1.amazonaws.com/sample-box/canary-deployment.mp4)