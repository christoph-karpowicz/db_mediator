main:
	go build -o server cmd/server/main.go
	go build -o unifier cmd/client/main.go

run:
	go run cmd/server/main.go
	go run cmd/client/main.go

test:
	go test -v ./...

test-client:
	go test -v ./internal/client/...

test-server:
	go test -v ./internal/server/...

test-server-cfg:
	go test -v ./internal/server/cfg/...

test-db:
	go test -v ./internal/server/db/...

test-synch:
	go test -v ./internal/server/synch/...