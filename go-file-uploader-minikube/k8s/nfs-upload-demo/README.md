# NFS-backed uploader demo on Kubernetes

This package splits your setup into separate manifests so you can apply and troubleshoot them in the correct order.

## Files

- `00-namespace.yaml` - namespace
- `01-nfs-server.yaml` - NFS server Deployment and Service
- `02-pv.yaml` - static PersistentVolume
- `03-pvc.yaml` - PersistentVolumeClaim bound to the PV
- `04-uploader.yaml` - uploader Deployment and NodePort Service

## Apply order

```bash
kubectl apply -f 00-namespace.yaml
kubectl apply -f 01-nfs-server.yaml
kubectl apply -f 02-pv.yaml
kubectl apply -f 03-pvc.yaml
kubectl apply -f 04-uploader.yaml
```

## Verify

```bash
kubectl get pods -n upload-demo
kubectl get svc -n upload-demo
kubectl get pv
kubectl get pvc -n upload-demo
```

Expected:

- `shared-nfs-pv` should become `Bound`
- `shared-nfs-pvc` should become `Bound`
- `nfs-server` pod should be `Running`
- `uploader-nfs` pod should be `Running`

## Debug commands

```bash
kubectl describe pv shared-nfs-pv
kubectl describe pvc shared-nfs-pvc -n upload-demo
kubectl describe pod -n upload-demo -l app=uploader-nfs
kubectl logs -n upload-demo deploy/nfs-server
```

## Clean reset

If you change the PV NFS source (`server`, `path`, etc.), Kubernetes will reject in-place updates because the PV source is immutable.
In that case delete and recreate in this order:

```bash
kubectl delete -f 04-uploader.yaml --ignore-not-found
kubectl delete -f 03-pvc.yaml --ignore-not-found
kubectl delete -f 02-pv.yaml --ignore-not-found
kubectl apply -f 02-pv.yaml
kubectl apply -f 03-pvc.yaml
kubectl apply -f 04-uploader.yaml
```

## Notes

- `storageClassName` is set to `nfs-manual` in both PV and PVC so static binding works cleanly.
- `volumeName: shared-nfs-pv` forces the PVC to bind to the exact PV.
- The NFS export path is `/exports`, which matches the NFS container configuration.
- `go-file-uploader:v1` must already exist on the node or in a registry reachable by your cluster.
- `NodePort 30083` exposes the uploader service externally from the node.
