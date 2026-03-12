kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.35/deploy/local-path-storage.yaml

kubectl get storageclass

By default it uses /opt/local-path-provisioner. You can change that to something like:
	•	k8s-master → /data/local-path
	•	k8s-worker → /data/local-path

kubectl -n local-path-storage get configmap local-path-config -o yaml