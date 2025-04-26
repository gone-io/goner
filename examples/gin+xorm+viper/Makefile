build:
	go install github.com/gone-io/gonectr@latest
	go install go.uber.org/mock/mockgen@latest
	gonectr build -ldflags="-w -s" -tags musl -o bin/server ./cmd/server

run:
	go install github.com/gone-io/gonectr@latest
	go install go.uber.org/mock/mockgen@latest
	gonectr run ./cmd/server


build-docker:
	docker compose build