#!/bin/bash
set -e
set -v

export PEERS="127.0.0.1:8888 127.0.0.1:8889 127.0.0.1:8890"
go run . -port=8888 -log_dir=./log/peer1 && 
go run . -port=8889 -log_dir=./log/peer2 && 
go run . -port=8890 -log_dir=./log/peer3