# 1. Minikube Commands

Minikube is used to run a **local Kubernetes cluster** for development and learning.

## Start a Cluster

```bash
minikube start
```

Starts a single-node Kubernetes cluster locally.

```bash
minikube start --nodes 2
```

Starts a **multi-node cluster** with two nodes.

Useful for simulating production-like scheduling.

---

## Check Cluster Status

```bash
minikube status
```

Shows the state of:

- Kubernetes
- Host VM
- kubelet
- apiserver

---

## View Logs

```bash
minikube logs
```

Shows complete Minikube logs.

```bash
minikube logs --tail 20
```

Shows last 20 lines of logs.

```bash
minikube logs -n 20
```

Alternative way to display the last few log lines.

---

## Node Management

Add a new node:

```bash
minikube node add
```

Add a node with a name:

```bash
minikube node add minikube-m02
```

Delete a node:

```bash
minikube node delete minikube-m02
```

Delete another node:

```bash
minikube node delete minikube-m03
```

---

## Delete Cluster

```bash
minikube delete
```

Completely removes the Minikube cluster.

Useful when you want to recreate a clean environment.

---

# 2. Kubernetes Commands (kubectl)

`kubectl` is the **Kubernetes command-line client** used to interact with the cluster.

---

# Cluster Information

Show cluster endpoints:

```bash
kubectl cluster-info
```

Dump full cluster configuration:

```bash
kubectl cluster-info dump
```

---

# Node Information

List nodes:

```bash
kubectl get nodes
kubectl get node
kubectl get no
```

Show detailed node information:

```bash
kubectl get node -o wide
```

Displays:

- Node IP
- OS
- Kernel
- Container runtime

---

# Namespace Management

Create namespace:

```bash
kubectl create ns test
```

List namespaces:

```bash
kubectl get ns
```

Delete namespace:

```bash
kubectl delete ns test
```

Deleting a namespace removes **all resources inside it**.

---

# Pod Creation

Create an nginx pod:

```bash
kubectl run nginx-pod --image nginx:latest --restart=Never --port=80 -n test
```

Create pod with wrong image (demonstration):

```bash
kubectl run nginx-pod1 --image nginx:latest1 --restart=Never --port=80 -n test
```

Run busybox container:

```bash
kubectl run box1 --image busybox -n test
```

Run a one-time command container:

```bash
kubectl run box2 -n test --image=busybox --restart=Never -- echo "Hello Kubernetes"
```

---

# Viewing Pods

List pods:

```bash
kubectl get po -n test
kubectl get pods -n test
```

Watch pod status:

```bash
kubectl get po -n test -w
```

Show extended details:

```bash
kubectl get pods -n test -o wide
```

Continuous monitoring:

```bash
watch kubectl get pods -n test
```

---

# Pod Debugging

Describe pod:

```bash
kubectl describe po nginx-pod -n test
```

Describe another pod:

```bash
kubectl describe po nginx-pod1 -n test
```

Describe busybox pod:

```bash
kubectl describe po box1 -n test
```

---

# Pod Logs

View logs:

```bash
kubectl logs box1 -n test
```

Logs with container name:

```bash
kubectl logs box1 -c box1 -n test
```

Logs from second container:

```bash
kubectl logs box2 -n test
```

---

# Delete Pods

Delete pod:

```bash
kubectl delete po nginx-pod1 -n test
```

Delete busybox pod:

```bash
kubectl delete po nginx-box -n test
```

---

# Working with YAML

Apply resource configuration:

```bash
kubectl apply -f create-nginx-pod.yaml
```

Apply service configuration:

```bash
kubectl apply -f create-nginx-service.yaml
```

Apply all YAML files in directory:

```bash
kubectl apply -f ./
```

Create resource directly:

```bash
kubectl create -f create-nginx-pod.yaml
```

---

# View Resources

List everything in namespace:

```bash
kubectl get all -n test
```

Show extended information:

```bash
kubectl get all -n test -o wide
```

List everything across namespaces:

```bash
kubectl get all -A
```

---

# MetalLB Installation

Install MetalLB:

```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.3/config/manifests/metallb-native.yaml
```

Verify installation:

```bash
kubectl get pods -n metallb-system
```

Apply IP pool:

```bash
kubectl apply -f metallb-pool.yaml
```

Apply Layer2 configuration:

```bash
kubectl apply -f metallb-l2.yaml
```

Test LoadBalancer IP:

```bash
curl http://192.168.1.240
curl http://192.168.49.240
```

---

# Useful Inspection Commands

List supported API resources:

```bash
kubectl api-resources
```

View kubeconfig contexts:

```bash
kubectl config get-contexts
```

View configuration:

```bash
kubectl config view
```

Use custom kubeconfig:

```bash
kubectl --kubeconfig=/home/labuser/.kube/config get pods -A
```

---

# Quick Cheat Sheet

## Minikube

```bash
minikube start --nodes 2
minikube status
minikube logs
minikube node add
minikube node delete
minikube delete
```

## Kubernetes

```bash
kubectl get nodes
kubectl get pods
kubectl describe pod
kubectl logs pod
kubectl apply -f file.yaml
kubectl get all -A
kubectl delete pod
```

---

## Deployment/Scale/Rollout

```bash

kubectl apply -f create-nginx-deployment.yaml
kubectl scale deployment nginx-deployment --replicas=10 -n test
kubectl set image deployment/nginx-deployment nginx=nginx:1.26 -n test
kubectl rollout status deployment/nginx-deployment -n test
kubectl rollout undo deployment/nginx-deployment -n test
```

---

# Summary

Typical workflow:

Start cluster:

```bash
minikube start
```

Deploy application:

```bash
kubectl create ns test
kubectl apply -f create-nginx-pod.yaml
kubectl apply -f create-nginx-service.yaml
```

Inspect resources:

```bash
kubectl get all -n test
kubectl describe pod nginx-pod -n test
```

Test service:

```bash
curl http://<node-ip>:<nodeport>
```
