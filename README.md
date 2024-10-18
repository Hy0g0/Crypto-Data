# Crypto-Data

# TODO

- [x] Create basic readme to allow all group to use and work on project
- [x] Create python scraper to get cryptos
- [x] Make scraper push in redpanda using kafka API
- [x] Add monitoring to Kube cluster and Redpanda
- [x] Deploy redpanda cluster with wright configuration
- [ ] Add distributed Database to save transformed data (cassandra,scyllaDB,CockroachDB)
- [ ] Update readme with redpanda dashboard creation
- [ ] Investigate Clickhouse project to find if it fits our plans
- [ ] Create final front to display our datas
- [ ] Create Consumer to transform data and push it in DB

## Setup ## 

Install cert-manager using Helm:

```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --set crds.enabled=true \
  --namespace cert-manager  \
  --create-namespace
```

The Redpanda Helm chart enables TLS by default and uses cert-manager to manage TLS certificates.

```bash
kubectl create namespace redpanda
kubectl create configmap redpanda-io-config --namespace redpanda --from-file=io-config.yaml

```


```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus-operator prometheus-community/kube-prometheus-stack \
  --namespace default \
  --create-namespace
  --values prometheus-values.yaml
```
Install the Redpanda Helm chart to deploy a Redpanda cluster and Redpanda Console.



```bash
helm repo add redpanda https://charts.redpanda.com
helm install redpanda redpanda/redpanda \
  --version 5.9.4 \
  --namespace redpanda \
  --create-namespace \
  --values redpanda-values.yaml
```


## OPS ##

```bash
helm upgrade --install redpanda redpanda/redpanda \
  --namespace redpanda \
  --create-namespace \
  --values redpanda-values.yaml
```

## SETUP CLUSTER ##

### monitoring

```bash
kubectl apply -f redpanda-servicemonitor.yaml
```

### create topic 

```bash
kubectl exec -it redpanda-0 -n redpanda -- rpk topic create crypto-prices --partitions 3 -c compression.type=lz4 
```

### create producer

``` bash
docker build -t redpanda-producer:latest ./scraper  
```

``` bash
kubectl apply -f redpanda-producer.yaml -n redpanda 
```

Note : this way of creating producer work for `orbstack` kube cluster for any other you may need to find a way to push image to the cluster repository

## Sources ##

- ### [Redpanda](https://docs.redpanda.com/current/deploy/deployment-option/self-hosted/kubernetes/k-production-deployment/)
- ### [Helm](https://helm.sh/docs/)
- ### [minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download)
- ### [kind](https://kind.sigs.k8s.io/)
- ### [orbstack](https://orbstack.dev/download)