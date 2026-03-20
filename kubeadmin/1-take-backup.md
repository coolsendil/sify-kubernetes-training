- Take compelte role back up

kubectl get clusterrolebinding -o yaml > crb-backup.yaml

- take complete role binding back up

kubectl get rolebinding -A -o yaml > rb-backup.yaml

-- Another Approach 

Stop Api server port access during backup ..


# Freeze
iptables -A INPUT -p tcp --dport 6443 -j DROP

sleep 10

# Take backup
etcdctl ... snapshot save backup.db

# Resume
iptables -D INPUT -p tcp --dport 6443 -j DROP



### Cert transfer

multipass exec k8s-master -- sudo cp /etc/kubernetes/pki/etcd/healthcheck-client.key /home/ubuntu/
multipass exec k8s-master -- sudo chown ubuntu:ubuntu /home/ubuntu/healthcheck-client.key
multipass transfer k8s-master:/home/ubuntu/healthcheck-client.key .