@PHONY:run
run:
	go run main.go crl --debug

@PHONY:test
test:
	go test ./...

@PHONY: lint
lint:
	golangci-lint run --out-format=github-actions

@PHONY: build
build:
	go build -o cg

@PHONY: gif
gif: build
	vhs cassette.tape

@PHONY: sqlc
sqlc:
	sqlc generate -f internal/adapter/db/sqlc.yaml