#/bin/sh

go run client/*.go \
    --master-node-address=localhost:8081 \
    --workflow-file=./client/sample-workflow/workflow.yaml