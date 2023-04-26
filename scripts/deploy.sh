#!/bin/bash

source ./config/.env

HOST=94.198.216.89

CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -o $(pwd)/bin/auth $(pwd)/cmd/auth/main.go

read -p "enter server username:" name;
scp -o LogLevel=Error -r $(pwd)/* "${name:-root}@${HOST}:~/${SVC}"
if [ $? -ne 0 ]; then
    exit echo $?
fi

sshpass ssh -o stricthostkeychecking=no -T root@${HOST} "
  cd ${SVC} &&
  echo $SSHPASS | sudo -S docker compose \
    -f ~/${SVC}/deployments/docker-compose.production.yml \
    -p deployments up \
    -d \
    --build && \
    exit
"