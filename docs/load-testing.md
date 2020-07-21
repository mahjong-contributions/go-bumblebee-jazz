# Load testing

## Locust

1. Access to Locust web console.

```bash

kubectl -n go-bumblebee port-forward svc/locust 8989:8989

# open http://localhost:8989 (MacOS)

```