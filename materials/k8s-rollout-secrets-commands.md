# Kubernetes Commands Reference

## Rollout Commands

### Show rollout history
```bash
kubectl rollout history deployment/go-app -n db-demo-ns
kubectl rollout history deploy go-app -n db-demo-ns
```

### Show a specific revision
```bash
kubectl rollout history deployment/go-app -n db-demo-ns --revision=2
kubectl rollout history deployment/nginx-deployment --revision=2
```

### Roll back to previous revision
```bash
kubectl rollout undo deployment/go-app -n db-demo-ns
```

### Roll back to a specific revision
```bash
kubectl rollout undo deployment/go-app -n db-demo-ns --to-revision=1
kubectl rollout undo deployment/go-app -n db-demo-ns --to-revision=2
```

### Check rollout status
```bash
kubectl rollout status deployment/go-app -n db-demo-ns
```

### Scale deployment
```bash
kubectl scale deployment go-app --replicas=5 -n db-demo-ns
```

### Edit deployment
```bash
kubectl edit deployment go-app -n db-demo-ns
```

---

## Secret Commands

### List secrets
```bash
kubectl get secrets -A
kubectl get secret -n db-demo-ns
```

### Edit a secret
```bash
kubectl edit secret mytls -n db-demo-ns
kubectl edit secret bootstrap-token-rp5np5 -n kube-system
```

### Explain secret type
```bash
kubectl explain secret.type
```

### Create a TLS secret
```bash
kubectl create secret tls mytls --cert=tls.crt --key=tls.key -n db-demo-ns
```

### Create a token for a ServiceAccount
```bash
kubectl create token app-token
```

### View a secret in YAML
```bash
kubectl get secret mytls -n db-demo-ns -o yaml
```

### Decode a secret value
```bash
kubectl get secret <secret-name> -n <namespace> -o jsonpath='{.data.<key>}' | base64 -d
```

---

## ConfigMap Commands

### List ConfigMaps
```bash
kubectl get configmap -n db-demo-ns
```

### Edit ConfigMap
```bash
kubectl edit configmap postgres-config -n db-demo-ns
```

### Describe ConfigMap
```bash
kubectl describe configmap postgres-config -n db-demo-ns
```

---

## Namespace Commands

### Create namespace
```bash
kubectl create ns db-demo-ns
```

### Delete namespace
```bash
kubectl delete ns db-demo-ns
```

---

## Apply / Delete Manifests

### Apply all YAML files in a directory
```bash
kubectl apply -f .
```

### Delete all YAML files in a directory
```bash
kubectl delete -f .
```

### Apply a specific file
```bash
kubectl apply -f 04-app-deployment.yaml
```

---

## General Inspection Commands

### Get all resources in a namespace
```bash
kubectl get all -n db-demo-ns
kubectl get all -n db-demo-ns -o wide
```

### Describe a pod
```bash
kubectl describe pod/go-app-cd5fb6c5c-gvcff -n db-demo-ns
kubectl describe pod/go-app-cd5fb6c5c-2vgq5 -n db-demo-ns
kubectl describe pod/go-app-7674c74646-2tq72 -n db-demo-ns
```

### Check logs
```bash
kubectl logs pod/go-app-cd5fb6c5c-gvcff -n db-demo-ns
```

### Get nodes
```bash
kubectl get node -o wide
```

---

## Minikube Image Commands

### Build image inside Minikube
```bash
minikube image build -t jpalaparthi/sify-go-demo:v1.0 .
minikube image build -t jpalaparthi/sify-go-demo:v1.1 .
minikube image build -t jpalaparthi/sify-go-demo:v1.2 .
```

### List images in Minikube
```bash
minikube image list
```

---

## Docker Commands

### Build Docker image
```bash
docker build -t jpalaparthi/sify-go-demo:v1.0 .
```

### List Docker images
```bash
docker images
docker image ls
```

### Docker login
```bash
docker login
docker login -u jpalaparthi
```

---

## TLS Certificate Generation

### Generate self-signed certificate
```bash
openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout tls.key \
  -out tls.crt \
  -subj "/CN=myapp.local/O=myapp.local"
```

---

## App Test Commands

### Test application endpoints
```bash
curl 192.168.49.2:30080/
curl 192.168.49.2:30080/ping
curl 192.168.49.2:30080/greet
curl 192.168.49.2:30080/health
curl 192.168.49.2:30080/healthz
```

---

## Recommended Rollout Workflow

```bash
# 1. Check current rollout history
kubectl rollout history deployment/go-app -n db-demo-ns

# 2. Edit or update deployment
kubectl edit deployment go-app -n db-demo-ns

# 3. Check rollout status
kubectl rollout status deployment/go-app -n db-demo-ns

# 4. Verify new pods
kubectl get all -n db-demo-ns -o wide

# 5. Test app
curl 192.168.49.2:30080/greet

# 6. If needed, rollback
kubectl rollout undo deployment/go-app -n db-demo-ns --to-revision=1
```

---

## Notes

- `deployment` and `deploy` are the clearest forms to use.
- For passwords and tokens, prefer `Secret` over `ConfigMap`.
- For production, add a change-cause annotation so rollout history is meaningful.
