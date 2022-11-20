server:
	go run cmd/main.go

test:
	go test -v -count=1 -cover ./...

.PHONY: server test
