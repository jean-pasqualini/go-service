SHELL := /bin/bash

# Building containers
all: sales-api

sales-api:
	docker build \
		-f zarf/docker/Dockerfile.sales-api \
		-t sales-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# Running local
run:
	go run app/sales-api/main.go

# Mod management
tidy:
	go mod tidy
	# go mid vendor (only to version the mods)

# Expvar monitoring to see the stats in real time
expvarmon:
	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# Loading capacity test
hey:
	hey -m GET -c 100 -n 1000000 "http://localhost:3000/readiness"

# Administration
genkeys:
	@go run app/admin/main.go

# Running tests within the local computer
test:
	go test -v ./... -count=1
	staticcheck ./...