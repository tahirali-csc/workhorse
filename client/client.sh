#/bin/sh

go run client/*.go \
    --master-node-address=localhost \
    --workflow-file=/home/tahir/workspace/rnd-projects/workhorse/client/sample-workflow/workflow.yaml
