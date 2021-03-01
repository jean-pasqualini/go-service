SHELL := /bin/bash

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