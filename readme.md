**Prerequisites**

1. Docker
2. Go Installation
3. Ubuntu or MacOS
4. Optionally, VirtualBox for testing on Cluster

**How to run all components on local machine?**

From the application directory:

1. Open a new terminal tab and run worker process
```bash
./worker/worker.sh
```

2. Open a new terminal tab and run master node process
```bash
./server/server.sh
```

3. Open a new terminal tab and run your workflow
```bash
./client/client.sh
```

This runs a sample workflow. The example workflow is present at **client/sample-workflow**

```yaml
jobs:
    -
        name: job1
        script: job1.sh
        image: alpine
    -
        name: job2
        script: job2.sh
        image: alpine
```

**Run on Cluster**

For this scenario, we can setup a cluster using any hypervisor like VirtualBox or VMWare. For this example, we setup following:

1. Setup two Ubuntu Worker Nodes.
2. Setup one Ubuntu Master Node.
3. Run worker node process on Worker Node.
```bash
./worker/worker.sh
```
4. Run master node process on master node
```bash
./server/server.sh --worker-node-address=ip address of worker node 1:8080, ip address of worker node 2:8080
```
5. From a different machine, run client workflow
```bash
./client/client.sh --master-node-address=[master node address]:8081
```