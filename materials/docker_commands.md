
# Docker Commands Cheat Sheet

## 1. Docker Version & Info

```bash
docker version
docker info
```

- `docker version` → Shows Docker client and server version.
- `docker info` → Displays system-wide Docker information.

---

# 2. Docker Images

## List Images

```bash
docker images
docker image ls
```

## Pull Image

```bash
docker pull nginx
docker pull ubuntu:22.04
```

## Remove Image

```bash
docker rmi IMAGE_ID
docker image rm IMAGE_ID
```

## Build Image

```bash
docker build -t myimage .
docker build -t myimage:1.0 .
```

---

# 3. Docker Containers

## Run Container

```bash
docker run nginx
docker run -d nginx
docker run -it ubuntu bash
```

Options:

- `-d` → Detached mode
- `-it` → Interactive terminal
- `--name` → Assign container name

Example:

```bash
docker run -d --name web nginx
```

---

## List Containers

```bash
docker ps
docker ps -a
```

- `docker ps` → Running containers
- `docker ps -a` → All containers

---

## Stop Container

```bash
docker stop CONTAINER_ID
```

---

## Start Container

```bash
docker start CONTAINER_ID
```

---

## Restart Container

```bash
docker restart CONTAINER_ID
```

---

## Remove Container

```bash
docker rm CONTAINER_ID
```

Force remove:

```bash
docker rm -f CONTAINER_ID
```

---

# 4. Container Logs & Monitoring

## View Logs

```bash
docker logs CONTAINER_ID
docker logs -f CONTAINER_ID
```

---

## Container Statistics

```bash
docker stats
```

Shows:

- CPU usage
- Memory usage
- Network IO

---

# 5. Execute Commands in Container

```bash
docker exec -it CONTAINER_ID bash
```

Example:

```bash
docker exec -it nginx_container bash
```

---

# 6. Port Mapping

Run container with port mapping:

```bash
docker run -p 8080:80 nginx
```

Meaning:

Host Port → Container Port  
8080 → 80

Access in browser:

http://localhost:8080

---

# 7. Volume Mounting

Mount local directory to container:

```bash
docker run -v /host/path:/container/path nginx
```

Example:

```bash
docker run -v $(pwd):/usr/share/nginx/html nginx
```

---

# 8. Docker Networks

## List Networks

```bash
docker network ls
```

## Create Network

```bash
docker network create mynetwork
```

## Run Container in Network

```bash
docker run -d --network=mynetwork nginx
```

---

# 9. Docker Inspect

Get container details:

```bash
docker inspect CONTAINER_ID
```

Example:

```bash
docker inspect nginx_container
```

---

# 10. Docker System Cleanup

Remove unused objects:

```bash
docker system prune
```

Remove everything unused:

```bash
docker system prune -a
```

---

# 11. Docker Compose

Start services:

```bash
docker compose up
```

Detached mode:

```bash
docker compose up -d
```

Stop services:

```bash
docker compose down
```

---

# 12. Docker Save & Load Images

Save image:

```bash
docker save nginx > nginx.tar
```

Load image:

```bash
docker load < nginx.tar
```

---

# 13. Docker Copy Files

Copy file from host to container:

```bash
docker cp file.txt CONTAINER_ID:/tmp
```

Copy file from container:

```bash
docker cp CONTAINER_ID:/tmp/file.txt .
```

---

# 14. Docker Daemon

Check daemon status (Linux):

```bash
systemctl status docker
```

Start Docker:

```bash
systemctl start docker
```

Restart Docker:

```bash
systemctl restart docker
```

---

# Example: Run Nginx Web Server

```bash
docker run -d -p 18080:80 --name nginx-server nginx
```

Access:

http://localhost:18080
