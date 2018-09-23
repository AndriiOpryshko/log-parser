#!/bin/bash
docker-compose up -d
(cd ../ && go get &&  go run *.go)
