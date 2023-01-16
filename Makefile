include .env
export $(shell sed 's/=.*//' .env)
current_dir = $(shell pwd)

build:
	@go build -o ./bin/blockchain-load-monitoring-service

run:
	@./bin/blockchain-load-monitoring-service

buildandrun:
	@go build -o ./bin/blockchain-load-monitoring-service && ./bin//blockchain-load-monitoring-service