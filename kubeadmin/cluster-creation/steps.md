multipass launch 24.04 --name k8s-master --cpus 2 --memory 4G --disk 20G --cloud-init master-cloud-init.yaml

multipass launch 24.04 --name k8s-worker --cpus 2 --memory 2G --disk 20G --cloud-init worker-cloud-init.yaml

multipass transfer k8s-master:/home/ubuntu/join-worker.sh k8s-worker:/home/ubuntu/join-worker.sh

multipass transfer k8s-master:/home/ubuntu/join-worker.sh . && \
multipass transfer ./join-worker.sh k8s-worker:/home/ubuntu/join-worker.sh

