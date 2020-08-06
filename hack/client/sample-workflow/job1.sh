#!/bin/sh
machine=$(hostname)
echo "This is script #1:: $machine"
ls -la
echo "Sleeping long..."
sleep 20s
echo "Awake..."
date
echo "Finished script##1"

