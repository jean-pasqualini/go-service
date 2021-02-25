SHELL := /bin/bash

run:
	go run app/sales-api/main.go

tidy:
	go mod tidy
	# go mid vendor (only to version the mods)