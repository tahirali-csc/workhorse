#/bin/sh

# go run client/*.go \
#     --master-node-address=192.168.56.102:8081 \
#     --workflow-file=./client/sample-workflow/workflow.yaml

go run ./cmd/client/*.go \
    --master-node-address=localhost:8081 \
    --workflow-file=./hack/client/sample-workflow/workflow.yaml