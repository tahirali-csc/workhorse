#!/bin/sh
machine=$(hostname)
echo "This is script #2:: $machine"
ls -la /etc
echo "Sleeping..."
sleep 5s
date
echo "Finished script##2"
