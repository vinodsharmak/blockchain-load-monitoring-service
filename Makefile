include .env
export $(shell sed 's/=.*//' .env)
current_dir = $(shell pwd)

build:
	@go build -o blockchain-load-monitoring-service

run:
	@./blockchain-load-monitoring-service

buildandrun:
	@go build -o blockchain-load-monitoring-service && ./blockchain-load-monitoring-service