
# PromQL Queries for Node Exporter (Kubernetes Nodes)

## 1. Check if Node Exporter Targets Are Up
```promql
up{job="node-exporter"}
```

---

## 2. CPU Usage Percentage Per Node
```promql
100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
```

---

## 3. CPU Usage Per Core
```promql
100 - (rate(node_cpu_seconds_total{mode="idle"}[5m]) * 100)
```

---

## 4. Total CPU Cores Per Node
```promql
count(node_cpu_seconds_total{mode="system"}) by (instance)
```

---

## 5. Memory Usage Percentage
```promql
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100
```

---

## 6. Memory Available
```promql
node_memory_MemAvailable_bytes
```

---

## 7. Total Memory
```promql
node_memory_MemTotal_bytes
```

---

## 8. Disk Usage Percentage
```promql
100 * (1 - (node_filesystem_avail_bytes{fstype!="tmpfs"} / node_filesystem_size_bytes{fstype!="tmpfs"}))
```

---

## 9. Disk Read Throughput
```promql
rate(node_disk_read_bytes_total[5m])
```

---

## 10. Disk Write Throughput
```promql
rate(node_disk_written_bytes_total[5m])
```

---

## 11. Network Receive Bandwidth
```promql
rate(node_network_receive_bytes_total[5m])
```

---

## 12. Network Transmit Bandwidth
```promql
rate(node_network_transmit_bytes_total[5m])
```

---

## 13. Node Load Average
```promql
node_load1
```

---

## 14. System Uptime
```promql
node_time_seconds - node_boot_time_seconds
```

---

## 15. Filesystem Free Space
```promql
node_filesystem_avail_bytes
```

---

## 16. Number of Running Processes
```promql
node_procs_running
```

---

## 17. Open File Descriptors
```promql
node_filefd_allocated
```

---

## 18. Network Errors
```promql
rate(node_network_receive_errs_total[5m])
```

---

## 19. Context Switch Rate
```promql
rate(node_context_switches_total[5m])
```

---

## 20. Interrupt Rate
```promql
rate(node_intr_total[5m])
```
