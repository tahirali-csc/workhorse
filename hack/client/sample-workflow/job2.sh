#!/bin/sh
# set -x
machine=$(hostname)
# echo "This is script #2:: $machine"
# ls -la /etc
echo "Sleeping...(Job2)"
apk add git
git version
# sleep 20s
# echo "Yeaaaaaaa!"
# date
# echo "Finished script##2"

echo "This is RI"