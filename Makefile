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
gif:
	vhs cassette.tape
