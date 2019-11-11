main:
	go build -o app cmd/unifier/main.go

run:
	go run cmd/unifier/main.go

tests:
	go test -v ./...