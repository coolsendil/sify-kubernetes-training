docker build -t jpalaparthi/gofile-store:1.0.0 .

docker push jpalaparthi/gofile-store:1.0.0

kubectl apply -f k8s-hostpath-app.yaml


kubectl get pods -n gofile-store -o wide


ls -l /var/lib/gofile-store


curl -X POST http://:30080/write-text \
  -H "Content-Type: application/json" \
  -d '{
    "file_name": "hello.txt",
    "content": "Hello from Go REST service running on Kubernetes hostPath volume"
  }'