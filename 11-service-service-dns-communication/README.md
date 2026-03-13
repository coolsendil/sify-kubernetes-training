kubectl apply -f 0-nginx.yaml

kubectl get svc -n ns-b

kubectl apply -f 1.client.yaml

kubectl get pods -n ns-a

kubectl exec -it service-a -n ns-a -- sh

curl service-b.ns-b.svc.cluster.local