docker build -t jpalaparthi/gofile-store:1.0.0 .

docker push jpalaparthi/gofile-store:1.0.0

kubectl apply -f k8s-hostpath-app.yaml


kubectl get pods -n gofile-store -o wide


ls -l /var/lib/gofile-store