# OpenEBS

- Install Helm on the node or sytem you run to install provisioner

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

help version
```

- Install openebs

```bash
helm repo add openebs https://openebs.github.io/openebs
helm repo update

helm install openebs openebs/openebs --namespace openebs --create-namespace
```

- Verify

```bash
kubectl get pods -n openebs
```
