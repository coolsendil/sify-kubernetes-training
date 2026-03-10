
### Download and apply metallb resources 

```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.3/config/manifests/metallb-native.yaml
```

1. check for the ip address pool to be given , note: number of ips are more or equual to the number of services to run 
2. give the valid ip address range in the metallb-pool.yaml file
3. run the metallb-pool.yaml and then the metallb-l2.yaml file for the ARP protocol

```bash
kubectl apply -f metallb-pool.yaml
kubectl apply -f metallb-l2.yaml
```

4. when run the following command, you get the EXTERNAL-IP of that service

```bash
kubectl get svc -n test
```