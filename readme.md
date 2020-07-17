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

For local testing, server.sh uses local work node process.

```bash
go run server/*.go
```

3. Open a new terminal tab and run your workflow.
```bash
./client/client.sh
```

For local testing, client.sh points to local master node.
```bash
go run client/*.go \
    --master-node-address=localhost:8081 \
    --workflow-file=./hack/client/sample-workflow/workflow.yaml
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

1. Setup three Ubuntu VM on VirtualBox. Use NAT gateway adapter for internet and Hostonly adapter for local network.

2. Do a git clone of the project at worker node # 1. Start the worker node process
```bash
./worker/worker.sh
```
3. Do a git clone of the project at worker node # 2. Start the worker node process
```bash
./worker/worker.sh
```

3. Do a git clone of the project at master node VM. Start the master node process
```bash
./worker/server.sh
```


4. From a different machine, run client workflow
```bash
./client/client.sh
```

Make sure to point to master node ip address. Assuming maste node ip is 192.168.56.102


```bash
go run client/*.go \
     --master-node-address=192.168.56.102:8081 \
     --workflow-file=./client/sample-workflow/workflow.yaml
```     