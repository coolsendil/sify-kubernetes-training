docker build -t jpalaparthi/gofile-store:1.0.0 .

docker push jpalaparthi/gofile-store:1.0.0

kubectl apply -f k8s-hostpath-app.yaml


kubectl get pods -n gofile-store -o wide


ls -l /var/lib/gofile-store


curl -X POST http://<WORKER-NODE-IO>:30080/write-text \
  -H "Content-Type: application/json" \
  -d '{
    "file_name": "hello.txt",
    "content": "Hello from Go REST service running on Kubernetes hostPath volume"
  }'

curl -X POST http://10.38.154.232:30080/write-text   -H "Content-Type: application/json"   -d '{
    "file_name": "hello3.txt",
    "content": "Hello World!"
  }'


curl -X POST http://10.38.154.232:30080/write-text   -H "Content-Type: application/json"   -d '{
    "file_name": "hello2.txt",
    "content": "Hello World hey!"
  }'

  curl -X POST http://10.38.154.232:30080/write-text   -H "Content-Type: application/json"   -d '{
    "file_name": "hello4.txt",
    "content": "Hello World hey!"
  }'

   curl -X POST http://10.38.154.232:30080/write-text   -H "Content-Type: application/json"   -d '{
    "file_name": "hello5.txt",
    "content": "Hello World hey!"
  }'