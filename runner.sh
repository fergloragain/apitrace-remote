#!/bin/bash

t=1
while [ $t -ne 0 ]
do
	git pull
	go build -o main cmd/server/*go
	./main > log 2>&1
	t=$?
	sleep 2
done
