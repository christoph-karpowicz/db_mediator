main:
	go build -o server cmd/server/main.go
	go build -o client cmd/client/main.go

run:
	go run cmd/server/main.go
	go run cmd/client/main.go

test:
	go test -v ./...