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
run-admin:
	@go run app/admin/main.go

# Running tests within the local computer
test:
	go test -v ./... -count=1
	staticcheck ./...

# Running from within k8s/dev

kind-up:
	kind create cluster --image kindest/node:v1.19.1 --name jean-cluster --config zarf/k8s/dev/kind-config.yaml

kind-down:
	kind delete cluster --name jean-cluster

kind-load:
	kind load docker-image sales-api-amd64:1.0 --name jean-cluster
	#kind load docker-image metrics-amd64:1.0 --name jean-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

kind-sales-api: sales-api # reload the pod
	kind load docker-image sales-api-amd64:1.0 --name jean-cluster
	kubectl delete pods -lapp=sales-api

kind-logs:
	kubectl logs -lapp=sales-api --all-containers=true -f

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-status-full:
	kubectl describe pod -lapp=sales-api

kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep sales-api | cut -c1-26) --container app -- bash

kind-database:
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed

kind-delete:
	kustomize build zarf/k8s/dev | kubectl delete -f -