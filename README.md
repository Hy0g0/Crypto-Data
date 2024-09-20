# Crypto-Data

## Setup ## 

Install cert-manager using Helm:

```
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --set crds.enabled=true \
  --namespace cert-manager  \
  --create-namespace
```

The Redpanda Helm chart enables TLS by default and uses cert-manager to manage TLS certificates.

```
kubectl create namespace redpanda
kubectl create configmap redpanda-io-config --namespace redpanda --from-file=io-config.yaml

```


```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus-operator prometheus-community/kube-prometheus-stack \
  --namespace default \
  --create-namespace
```
Install the Redpanda Helm chart to deploy a Redpanda cluster and Redpanda Console.



```
helm repo add redpanda https://charts.redpanda.com
helm install redpanda redpanda/redpanda \
  --version 5.9.4 \
  --namespace redpanda \
  --create-namespace \
  --values redpanda-values.yaml
```


## OPS ##

```
helm upgrade --install redpanda redpanda/redpanda \
  --namespace redpanda \
  --create-namespace \
  --values redpanda-values.yaml
```
## Sources ##

[Redpanda](https://docs.redpanda.com/current/deploy/deployment-option/self-hosted/kubernetes/k-production-deployment/)