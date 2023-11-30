#!/bin/bash

cd /home/ec2-user/src/Back-End
git fetch origin
git pull origin
echo "git pull..."

pid=$(ps -ef | grep nftime | grep -v grep |  awk '{print $2}')
echo "found nftime process with pid" $pid
sudo kill -9 $pid
echo "kill command exited with status" $?

sleep 1

echo $(pwd)
swag init
go build
sudo nohup ./nftime >> nftime.log -env=dev &
echo "nohup done with status" $?

exit