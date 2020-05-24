#/bin/sh

go run server/*.go --worker-node-address=localhost:8080,192.168.56.103:8080 --schedule=random
#go run server/*.go --worker-node-address=localhost:8080
