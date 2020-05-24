#/bin/sh

go run client/*.go \
    --master-node-address=localhost \
    --workflow-file=./client/sample-workflow/workflow.yaml