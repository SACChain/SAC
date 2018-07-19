.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/SACChain/SAC/cmd/swarm
	go build -o bin/geth     github.com/SACChain/SAC/cmd/geth
	go build -o bin/bootnode github.com/SACChain/SAC/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/SACChain/SAC/swarm/...
