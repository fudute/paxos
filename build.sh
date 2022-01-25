#!/bin/bash

rm -rf output
mkdir output

cp -R config output/config

go build -o output/server ./app/server/main.go
