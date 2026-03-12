
# Multipass Commands Cheat Sheet

Multipass is a lightweight VM manager by Canonical that uses Ubuntu cloud images.  
It is commonly used for development environments, Kubernetes labs, and container testing.

---

# 1. Check Multipass Version

```bash
multipass version
```

---

# 2. List Virtual Machines

```bash
multipass list
```

Shows:
- Instance name
- State
- IP address
- CPU / RAM / Disk

---

# 3. Launch a New VM

```bash
multipass launch
```

Launch with a name:

```bash
multipass launch --name myvm
```

Launch with resources:

```bash
multipass launch --name devvm --cpus 2 --mem 2G --disk 10G
```

---

# 4. Enter a VM Shell

```bash
multipass shell myvm
```

---

# 5. Execute Command Inside VM

```bash
multipass exec myvm -- ls /
```

Example:

```bash
multipass exec myvm -- uname -a
```

---

# 6. Stop a VM

```bash
multipass stop myvm
```

Stop all instances:

```bash
multipass stop --all
```

---

# 7. Start a VM

```bash
multipass start myvm
```

Start all:

```bash
multipass start --all
```

---

# 8. Restart a VM

```bash
multipass restart myvm
```

---

# 9. Delete a VM

```bash
multipass delete myvm
```

Delete all:

```bash
multipass delete --all
```

---

# 10. Purge Deleted VMs

Multipass keeps deleted VMs in trash until purged.

```bash
multipass purge
```

---

# 11. Show VM Info

```bash
multipass info myvm
```

Example output:

- IP address
- Disk usage
- CPU count
- Memory

---

# 12. Mount Host Directory into VM

```bash
multipass mount /host/path myvm:/vm/path
```

Example:

```bash
multipass mount ~/projects myvm:/home/ubuntu/projects
```

---

# 13. Unmount Directory

```bash
multipass umount myvm
```

---

# 14. Transfer Files

Copy file to VM:

```bash
multipass transfer file.txt myvm:/home/ubuntu/
```

Copy file from VM:

```bash
multipass transfer myvm:/home/ubuntu/file.txt .
```

---

# 15. Get IP Address

```bash
multipass info myvm
```

or

```bash
multipass list
```

---

# 16. Show Available Ubuntu Images

```bash
multipass find
```

Example:

- 20.04
- 22.04
- 24.04

---

# 17. Launch Specific Ubuntu Version

```bash
multipass launch 22.04 --name ubuntu22
```

---

# 18. Run Kubernetes Lab Example

Create master node:

```bash
multipass launch --name k8s-master --cpus 2 --mem 4G --disk 20G
```

Create worker node:

```bash
multipass launch --name k8s-worker --cpus 2 --mem 2G --disk 20G
```

---

# 19. Networking

Check VM IP:

```bash
multipass list
```

Multipass networking types:

- NAT (default)
- Bridged networking

Example bridged launch:

```bash
multipass launch --network en0 --name bridged-vm
```

---

# 20. Multipass Daemon Troubleshooting

Check service:

Linux:

```bash
systemctl status multipass
```

Mac:

```bash
sudo launchctl list | grep multipass
```

Restart daemon (Mac):

```bash
sudo launchctl kickstart -k system/com.canonical.multipassd
```

---

# Example Workflow

```bash
multipass launch --name devvm
multipass shell devvm
sudo apt update
sudo apt install docker.io -y
```

Now your VM is ready for development.

---

# Clean Up Everything

```bash
multipass delete --all
multipass purge
```

---

Multipass is extremely useful for:

- Kubernetes labs
- Docker testing
- Cloud-init experiments
- Ubuntu development environments
