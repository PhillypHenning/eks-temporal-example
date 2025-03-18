# Instructions
## Prepare Blue deployment
### 1. Create blue deployment
```bash
kubectl create namespace blue
helm install temporal-training-blue -n blue temporalio/temporal --version 0.52.0 -f Server/values.yaml
```

### 2. Create a standard service that forwards from default namespace to blue

Open the BlueGreen/frontend-service.yaml file and update the `externalName` to match the pod+namespace expected.

`externalName: temporal-training-blue-frontend.blue.svc.cluster.local`

```bash
kubectl apply -f BlueGreen/frontend-service.yaml
```

### 3. Verify the worker is able to connect to the blue namespace deployment via the default pass-through service

```bash
kubectl port-forward -n blue pod/<web pod> 8080 8080
```

## Prepare Green deployment

### 1. Create green deployment

```bash
kubectl create namespace green
helm install temporal-training-green -n green temporalio/temporal --version 0.52.0 -f Server/values.yaml
```

### 2. Update the default service that the worker is pointing at

Open the BlueGreen/frontend-service.yaml file and update the `externalName` to match the pod+namespace expected.

`temporal-training-blue-frontend.blue.svc.cluster.local` -> `temporal-training-green-frontend.green.svc.cluster.local`

```bash
kubectl apply -f BlueGreen/frontend-service.yaml
```
### 3. Verify that the green deployment can use the worker

```bash
kubectl port-forward pod/temporal-training-green-web-58d8b4b6d-66tgl -n green 8080 8080
```



