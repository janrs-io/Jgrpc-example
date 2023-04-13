.PHONY: tidy
tidy:
	# go get -u ./...
	go mod tidy

.PHONY: wire
wire:
	wire ./...

.PHONY: buf
buf:
	buf generate

.PHONY: run
run:
	go run cmd/main.go