1. Install MetaLLB

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.2/config/manifests/metallb-native.yaml

2. Get IP Address Range for MetaLLB 
    


k8s-client  192.168.2.10 
k8s-master  192.168.2.7 
k8s-storage 192.168.2.9     
k8s-worker  192.168.2.8      

Address-Range: 192.168.2.40-192.168.2.50

3. 