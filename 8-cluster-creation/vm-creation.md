multipass launch 24.04 \
  --name k8s-master \
  --cpus 2 \
  --memory 4G \
  --disk 20G \
  --cloud-init k8s-node.yaml


  multipass launch 24.04 \
  --name k8s-worker \
  --cpus 2 \
  --memory 4G \
  --disk 20G \
  --cloud-init k8s-node.yaml


  multipass launch 24.04 \
  --name k8s-storage \
  --cpus 2 \
  --memory 2G \
  --disk 20G \
  --cloud-init storage.yaml

   multipass launch 24.04 \
  --name k8s-client \
  --cpus 1 \
  --memory 1G \
  --disk 10G \
  --cloud-init client-node.yaml


kubeadm join 192.168.2.7:6443 --token jq01eb.nwx5d9op067vg7fp \
	--discovery-token-ca-cert-hash sha256:0d3c5633c6ea868f375dd27e6511be28fd42739379b2a5d09a58391a16ce9b85 

  on master

  sudo kubeadm init --pod-network-cidr=10.244.0.0/16

  mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

kubectl get nodes

kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml

kubeadm token create --print-join-command

kubeadm join 10.200.1.10:6443 --token xxx \
--discovery-token-ca-cert-hash sha256:xxxx



### Access related 

- in k8s-master 

```bash
sudo kubeadm kubeconfig user --client-name jp-viewer > jp-viewer.conf
```

```bash
- kubectl --kubeconfig ~/.kube/config create clusterrolebinding jp-viewer-binding \
  --clusterrole=view \
  --user=jp-viewer
  ```
- create a new instance call client 

```bash
multipass ssh k8s-client 
mkdir -p ~/.kube
multipass transfer k8s-master:/home/ubuntu/jp-viewer.conf ./jp-viewer.conf
multipass exec k8s-client -- mkdir -p /home/ubuntu/.kube
multipass transfer ./jp-viewer.conf k8s-client:/home/ubuntu/.kube/config
```

```bash
# create admin for a namespace called dev
kubectl create rolebinding jp-admin \
  --clusterrole=admin \
  --user=jp \
  --namespace=dev
  ```