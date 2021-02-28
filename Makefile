SHELL := /bin/bash

run:
	go run app/sales-api/main.go

tidy:
	go mod tidy
	# go mid vendor (only to version the mods)

expvarmon:
	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

hey:
	hey -m GET -c 100 -n 1000000 "http://localhost:3000/readiness"