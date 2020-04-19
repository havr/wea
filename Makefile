lint:
	golint ./...
	golangci-lint run

test:
	go test ./...

run:
	go run cmd/wea/wea.go

build:
	go build cmd/wea/wea.go
