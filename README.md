
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
- [Automated release pipline](docs/cicd.md#automated-release-pipline#automated-release-pipline)
- Blue/Green deployment with ArgoCD Rollouts 
- [Canary deployment with ArgoCD Rollouts](docs/cicd.md#canary-deployment-with-argocd-rollouts) 

### Traffic management
- [Tripping the circuit breaker](docs/traffic-management.md#)
- [Fault Injection](docs/traffic-management.md#)

### Observerbility: metrics, tracing, logs
- [Monitoring with Prometheus](docs/observability.md)
- [Dashboard with Grafana](docs/observability.md)
- [Distributed tracing with Jaeger](docs/observability.md)
- [Log analysis with FluentBit & Elasticsearch](docs/observability.md)

### Load testing
- [Load testing with Locust](docs/load-testing.md)
- Load testing with Fortio

### Secret management
- Public secret with Sealed-Secret


